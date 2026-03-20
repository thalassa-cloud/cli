package workloadidentityfederation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/thalassa-cloud/client-go/filters"
	clientiam "github.com/thalassa-cloud/client-go/iam"
	"github.com/thalassa-cloud/client-go/thalassa"
)

// wifResourceKey is a stable id for bootstrap-labelled resources. It must include the OIDC issuer so the
// same repo/ref on GitLab.com vs self-managed (or any two issuers) does not reuse the same service account.
func wifResourceKey(vcs, repository, providerSubject, issuer string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(vcs) + "\n" + strings.TrimSpace(repository) + "\n" + providerSubject + "\n" + normalizeIssuer(issuer)))
	return hex.EncodeToString(sum[:8])
}

func shortHexKey(s string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(s)))
	return hex.EncodeToString(sum[:4])
}

// issuerURLHostname returns the host part of an issuer URL for display names (e.g. gitlab.com).
func issuerURLHostname(issuer string) string {
	issuer = strings.TrimSpace(issuer)
	if issuer == "" {
		return ""
	}
	raw := issuer
	if !strings.Contains(issuer, "://") {
		raw = "https://" + issuer
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if h := u.Hostname(); h != "" {
		return h
	}
	// Issuer may be stored without scheme
	s := normalizeIssuer(issuer)
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	if i := strings.IndexByte(s, '/'); i >= 0 {
		s = s[:i]
	}
	return s
}

func bootstrapLabels(vcs, key string) map[string]string {
	return map[string]string{
		LabelManagedBy: ValueManagedBy,
		LabelWIFKey:    key,
		LabelWIFVCS:    vcs,
	}
}

// labelsMatch reports whether have contains every key/value in want (API list filters may be ignored).
func labelsMatch(have map[string]string, want map[string]string) bool {
	if len(want) == 0 {
		return false
	}
	for k, v := range want {
		if have[k] != v {
			return false
		}
	}
	return true
}

// scopesEqual compares allowed scope sets (order-insensitive).
func scopesEqual(a, b []clientiam.AccessCredentialsScope) bool {
	toSorted := func(s []clientiam.AccessCredentialsScope) []string {
		out := make([]string, 0, len(s))
		for _, x := range s {
			out = append(out, string(x))
		}
		slices.Sort(out)
		return out
	}
	return slices.Equal(toSorted(a), toSorted(b))
}

func trustedAudiencesEqual(a, b []string) bool {
	aa := append([]string(nil), a...)
	bb := append([]string(nil), b...)
	slices.Sort(aa)
	slices.Sort(bb)
	return slices.Equal(aa, bb)
}

func federatedIdentityNeedsBootstrapReconcile(fi *clientiam.FederatedIdentity, fiName, fiDesc, subject, vcs, key string, scopes []clientiam.AccessCredentialsScope, trustedAudiences []string) bool {
	wantLabels := bootstrapLabels(vcs, key)
	annSubj := ""
	if fi.Annotations != nil {
		annSubj = fi.Annotations["thalassa.cloud/wif.provider-subject"]
	}
	return !scopesEqual(fi.AllowedScopes, scopes) ||
		!trustedAudiencesEqual(fi.TrustedAudiences, trustedAudiences) ||
		(fi.AudienceMatchMode != "" && fi.AudienceMatchMode != clientiam.AudienceMatchModeAny) ||
		fi.Name != fiName ||
		fi.Description != fiDesc ||
		!labelsMatch(fi.Labels, wantLabels) ||
		annSubj != subject
}

func findProviderByIssuer(ctx context.Context, c *clientiam.Client, issuer string) (*clientiam.FederatedIdentityProvider, error) {
	want := normalizeIssuer(issuer)
	list, err := c.ListFederatedIdentityProviders(ctx, &clientiam.ListFederatedIdentityProvidersRequest{})
	if err != nil {
		return nil, err
	}
	for i := range list {
		p := &list[i]
		if normalizeIssuer(p.ProviderIssuer) == want {
			return p, nil
		}
	}
	return nil, nil
}

func findProviderByKubernetesClusterID(ctx context.Context, c *clientiam.Client, clusterIdentity string) (*clientiam.FederatedIdentityProvider, error) {
	clusterIdentity = strings.TrimSpace(clusterIdentity)
	if clusterIdentity == "" {
		return nil, nil
	}
	want := map[string]string{LabelKubernetesClusterID: clusterIdentity}
	list, err := c.ListFederatedIdentityProviders(ctx, &clientiam.ListFederatedIdentityProvidersRequest{
		Filters: []filters.Filter{
			&filters.LabelFilter{MatchLabels: want},
		},
	})
	if err != nil {
		return nil, err
	}
	var found *clientiam.FederatedIdentityProvider
	n := 0
	for i := range list {
		p := &list[i]
		if labelsMatch(p.Labels, want) {
			n++
			found = p
		}
	}
	if n > 1 {
		return nil, fmt.Errorf("multiple federated identity providers match label %s=%s", LabelKubernetesClusterID, clusterIdentity)
	}
	return found, nil
}

func createBootstrapProvider(ctx context.Context, c *clientiam.Client, name, description, issuer, vcs string) (*clientiam.FederatedIdentityProvider, error) {
	// Provider shares managed-by / wif-vcs with other bootstrap objects; wif-key is issuer-specific so
	// multiple GitLab issuers in one org do not collide on the provider label value.
	pKey := "issuer-" + shortHexKey(issuer)
	p, err := c.CreateFederatedIdentityProvider(ctx, clientiam.CreateFederatedIdentityProviderRequest{
		Name:           name,
		Description:    description,
		ProviderIssuer: issuer,
		Status:         clientiam.FederatedIdentityProviderStatusActive,
		Labels:         bootstrapLabels(vcs, pKey),
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func findServiceAccountByWIFKey(ctx context.Context, c *clientiam.Client, vcs, key string) (*clientiam.ServiceAccount, error) {
	want := bootstrapLabels(vcs, key)
	list, err := c.ListServiceAccounts(ctx, &clientiam.ListServiceAccountsRequest{
		Filters: []filters.Filter{
			&filters.LabelFilter{MatchLabels: want},
		},
	})
	if err != nil {
		return nil, err
	}
	for i := range list {
		sa := &list[i]
		if labelsMatch(sa.Labels, want) {
			return sa, nil
		}
	}
	return nil, nil
}

func createBootstrapServiceAccount(ctx context.Context, c *clientiam.Client, name, description string, vcs, key, repo, subject string) (*clientiam.ServiceAccount, error) {
	desc := description
	req := clientiam.CreateServiceAccountRequest{
		Name:   name,
		Labels: bootstrapLabels(vcs, key),
		Annotations: map[string]string{
			"thalassa.cloud/wif.repository":       repo,
			"thalassa.cloud/wif.provider-subject": subject,
		},
	}
	if desc != "" {
		req.Description = &desc
	}
	return c.CreateServiceAccount(ctx, req)
}

func resolveOrganisationRole(ctx context.Context, c *clientiam.Client, ref string) (*clientiam.OrganisationRole, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, fmt.Errorf("role is required")
	}
	if role, err := c.GetOrganisationRole(ctx, ref); err == nil && role != nil {
		return role, nil
	}
	roles, err := c.ListOrganisationRoles(ctx, &clientiam.ListOrganisationRolesRequest{})
	if err != nil {
		return nil, fmt.Errorf("list organisation roles: %w", err)
	}
	for i := range roles {
		r := &roles[i]
		if strings.EqualFold(r.Identity, ref) || strings.EqualFold(r.Slug, ref) || strings.EqualFold(r.Name, ref) {
			return r, nil
		}
	}
	return nil, fmt.Errorf("organisation role not found: %s", ref)
}

func hasRoleBindingForServiceAccount(ctx context.Context, c *clientiam.Client, roleIdentity, saIdentity string) (bool, error) {
	bindings, err := c.ListRoleBindings(ctx, roleIdentity, &clientiam.ListRoleBindingsRequest{})
	if err != nil {
		return false, err
	}
	for _, b := range bindings {
		if b.ServiceAccount != nil && b.ServiceAccount.Identity == saIdentity {
			return true, nil
		}
	}
	return false, nil
}

func createRoleBindingForSA(ctx context.Context, c *clientiam.Client, role *clientiam.OrganisationRole, sa *clientiam.ServiceAccount, vcs, key string) (*clientiam.OrganisationRoleBinding, error) {
	name := fmt.Sprintf("wif-%s-%s", vcs, key)
	if len(name) > 63 {
		name = name[:63]
	}
	saID := sa.Identity
	binding, err := c.CreateRoleBinding(ctx, role.Identity, clientiam.CreateRoleBinding{
		Name:                   name,
		Description:            fmt.Sprintf("Workload identity federation (%s) for Thalassa service account %s", vcs, sa.Slug),
		Labels:                 bootstrapLabels(vcs, key),
		ServiceAccountIdentity: &saID,
	})
	if err != nil {
		return nil, err
	}
	return binding, nil
}

func findFederatedIdentity(ctx context.Context, c *clientiam.Client, vcs, key string) (*clientiam.FederatedIdentity, error) {
	want := bootstrapLabels(vcs, key)
	list, err := c.ListFederatedIdentities(ctx, &clientiam.ListFederatedIdentitiesRequest{
		Filters: []filters.Filter{
			&filters.LabelFilter{MatchLabels: want},
		},
	})
	if err != nil {
		return nil, err
	}
	for i := range list {
		fi := &list[i]
		if labelsMatch(fi.Labels, want) {
			return fi, nil
		}
	}
	return nil, nil
}

func createBootstrapFederatedIdentity(ctx context.Context, c *clientiam.Client, name, description, providerID, saID, subject string, audiences []string, scopes []clientiam.AccessCredentialsScope, vcs, key string) (*clientiam.FederatedIdentity, error) {
	return c.CreateFederatedIdentity(ctx, clientiam.CreateFederatedIdentityRequest{
		Name:        name,
		Description: description,
		Labels:      bootstrapLabels(vcs, key),
		Annotations: map[string]string{
			"thalassa.cloud/wif.provider-subject": subject,
		},
		ServiceAccountIdentity: saID,
		ProviderIdentity:       providerID,
		ProviderSubject:        subject,
		TrustedAudiences:       audiences,
		AudienceMatchMode:      clientiam.AudienceMatchModeAny,
		AllowedScopes:          scopes,
	})
}

// RunBootstrap executes the bootstrap flow. client must be a full thalassa.Client (IAM and, for kubernetes, Kubernetes API).
func RunBootstrap(ctx context.Context, client thalassa.Client, opts BootstrapOptions) (*BootstrapResult, error) {
	iamc := client.IAM()
	res := &BootstrapResult{}

	var issuer string
	var subject string
	var err error
	// Set when --kubernetes-cluster resolves to a Thalassa cluster: OIDC provider must exist with label kubernetes_cluster_id=<cluster identity>.
	var k8sClusterBoundProvider *clientiam.FederatedIdentityProvider

	switch opts.VCS {
	case ValueVCSGitHub:
		issuer = GitHubActionsIssuer
		subject, err = BuildGitHubSubject(opts.Repository, opts.RefKind, opts.Ref)
	case ValueVCSGitLab:
		issuer = normalizeIssuer(opts.GitLabIssuer)
		if issuer == "" {
			issuer = "https://gitlab.com"
		}
		subject, err = BuildGitLabSubject(opts.Repository, opts.GitLabRefType, opts.Ref)
	case ValueVCSKubernetes:
		explicitIssuer := normalizeIssuer(opts.KubernetesIssuer)
		if strings.TrimSpace(opts.KubernetesClusterRef) != "" {
			cluster, rerr := resolveKubernetesCluster(ctx, client.Kubernetes(), opts.KubernetesClusterRef)
			if rerr != nil {
				return nil, rerr
			}
			p, perr := findProviderByKubernetesClusterID(ctx, iamc, cluster.Identity)
			if perr != nil {
				return nil, fmt.Errorf("list federated identity providers for cluster: %w", perr)
			}
			if p == nil {
				return nil, fmt.Errorf("no federated identity provider for kubernetes cluster %q (identity %s): expected exactly one provider with label %s=%s (cluster OIDC not provisioned yet?)",
					cluster.Name, cluster.Identity, LabelKubernetesClusterID, cluster.Identity)
			}
			k8sClusterBoundProvider = p
			issuer = normalizeIssuer(p.ProviderIssuer)
			if issuer == "" {
				return nil, fmt.Errorf("federated identity provider %s for cluster %q has an empty issuer URL", p.Identity, cluster.Name)
			}
			if explicitIssuer != "" && explicitIssuer != issuer {
				return nil, fmt.Errorf("kubernetes: --issuer %q does not match cluster OIDC provider issuer %q (omit --issuer when using --cluster)", opts.KubernetesIssuer, p.ProviderIssuer)
			}
		} else if explicitIssuer != "" {
			issuer = explicitIssuer
		} else {
			return nil, fmt.Errorf("kubernetes: set --kubernetes-cluster (Thalassa-managed cluster) and/or --kubernetes-issuer (self-managed or custom service-account issuer)")
		}
		subject, err = BuildKubernetesSubject(opts.Repository)
	default:
		return nil, fmt.Errorf("unsupported --vcs %q (use github, gitlab, or kubernetes)", opts.VCS)
	}
	if err != nil {
		return nil, err
	}
	if len(opts.TrustedAudiences) == 0 {
		return nil, fmt.Errorf("at least one trusted audience is required")
	}

	key := wifResourceKey(opts.VCS, opts.Repository, subject, issuer)
	res.WIFKey = key
	res.ProviderSubject = subject
	res.Issuer = issuer

	scopes := opts.AllowedScopes
	if len(scopes) == 0 {
		scopes = []clientiam.AccessCredentialsScope{
			clientiam.AccessCredentialsScopeAPIRead,
			clientiam.AccessCredentialsScopeAPIWrite,
		}
	}

	role, err := resolveOrganisationRole(ctx, iamc, opts.RoleRef)
	if err != nil {
		return nil, err
	}
	res.RoleIdentity = role.Identity
	res.RoleSlug = role.Slug

	var provider *clientiam.FederatedIdentityProvider
	if k8sClusterBoundProvider != nil {
		provider = k8sClusterBoundProvider
	} else {
		provider, err = findProviderByIssuer(ctx, iamc, issuer)
		if err != nil {
			return nil, fmt.Errorf("list federated identity providers: %w", err)
		}
	}
	if provider == nil {
		if opts.DryRun {
			res.WouldCreateProvider = true
			res.ProviderIdentity = "(dry-run: would create provider)"
		} else {
			pName := opts.ProviderDisplayName
			if pName == "" {
				switch opts.VCS {
				case ValueVCSGitHub:
					pName = "GitHub Actions OIDC"
				case ValueVCSGitLab:
					if h := issuerURLHostname(issuer); h != "" {
						pName = h
					} else {
						pName = "GitLab CI OIDC"
					}
				case ValueVCSKubernetes:
					if h := issuerURLHostname(issuer); h != "" {
						pName = h
					} else {
						pName = "Kubernetes service account OIDC"
					}
				default:
					pName = "OIDC"
				}
			}
			pDesc := opts.ProviderDescription
			if pDesc == "" {
				pDesc = fmt.Sprintf("OIDC issuer %s (created by tcloud iam workload-identity-federation bootstrap <github|gitlab|kubernetes>)", issuer)
			}
			p, err := createBootstrapProvider(ctx, iamc, pName, pDesc, issuer, opts.VCS)
			if err != nil {
				return nil, fmt.Errorf("create federated identity provider: %w", err)
			}
			provider = p
			res.CreatedProvider = true
		}
	}
	if provider != nil {
		res.ProviderIdentity = provider.Identity
	}

	saName := fmt.Sprintf("wif-%s-%s", opts.VCS, key)
	fiName := fmt.Sprintf("wif-%s-%s-fi", opts.VCS, key)
	if n := strings.TrimSpace(opts.ResourceName); n != "" {
		saName = n
		fiName = n + "-fi"
	}
	if len(saName) > 200 {
		saName = saName[:200]
	}
	if len(fiName) > 200 {
		fiName = fiName[:200]
	}

	sa, err := findServiceAccountByWIFKey(ctx, iamc, opts.VCS, key)
	if err != nil {
		return nil, fmt.Errorf("list service accounts: %w", err)
	}
	if sa == nil {
		if opts.DryRun {
			res.WouldCreateServiceAccount = true
			res.ServiceAccountIdentity = "(dry-run: would create service account)"
		} else {
			desc := fmt.Sprintf("Workload identity (%s %s); JWT sub: %s", opts.VCS, opts.Repository, subject)
			sa, err = createBootstrapServiceAccount(ctx, iamc, saName, desc, opts.VCS, key, opts.Repository, subject)
			if err != nil {
				return nil, fmt.Errorf("create service account: %w", err)
			}
			res.CreatedServiceAccount = true
		}
	}
	if sa != nil {
		res.ServiceAccountIdentity = sa.Identity
		res.ServiceAccountSlug = sa.Slug
	}

	fiDesc := fmt.Sprintf("Workload identity (%s %s) → %s", opts.VCS, opts.Repository, subject)
	fi, err := findFederatedIdentity(ctx, iamc, opts.VCS, key)
	if err != nil {
		return nil, fmt.Errorf("list federated identities: %w", err)
	}
	if fi == nil {
		if opts.DryRun {
			res.WouldCreateFederatedIdentity = true
			res.FederatedIdentityIdentity = "(dry-run: would create federated identity)"
		} else {
			if provider == nil || sa == nil {
				return nil, fmt.Errorf("internal: missing provider or service account after provisioning")
			}
			fi, err = createBootstrapFederatedIdentity(ctx, iamc, fiName, fiDesc, provider.Identity, sa.Identity, subject, opts.TrustedAudiences, scopes, opts.VCS, key)
			if err != nil {
				return nil, fmt.Errorf("create federated identity: %w", err)
			}
			res.CreatedFederatedIdentity = true
		}
	}
	if fi != nil {
		res.FederatedIdentityIdentity = fi.Identity
		if federatedIdentityNeedsBootstrapReconcile(fi, fiName, fiDesc, subject, opts.VCS, key, scopes, opts.TrustedAudiences) {
			if opts.DryRun {
				res.WouldUpdateFederatedIdentity = true
			} else {
				_, err := iamc.UpdateFederatedIdentity(ctx, fi.Identity, clientiam.UpdateFederatedIdentityRequest{
					Name:              fiName,
					Description:       fiDesc,
					Labels:            bootstrapLabels(opts.VCS, key),
					Annotations:       map[string]string{"thalassa.cloud/wif.provider-subject": subject},
					TrustedAudiences:  append([]string(nil), opts.TrustedAudiences...),
					AudienceMatchMode: clientiam.AudienceMatchModeAny,
					AllowedScopes:     scopes,
				})
				if err != nil {
					return nil, fmt.Errorf("update federated identity: %w", err)
				}
				res.UpdatedFederatedIdentity = true
			}
		}
	}

	if !opts.DryRun && sa != nil {
		ok, err := hasRoleBindingForServiceAccount(ctx, iamc, role.Identity, sa.Identity)
		if err != nil {
			return nil, fmt.Errorf("list role bindings: %w", err)
		}
		if !ok {
			_, err = createRoleBindingForSA(ctx, iamc, role, sa, opts.VCS, key)
			if err != nil {
				return nil, fmt.Errorf("create role binding: %w", err)
			}
			res.CreatedRoleBinding = true
		}
	} else if opts.DryRun {
		if sa != nil {
			ok, err := hasRoleBindingForServiceAccount(ctx, iamc, role.Identity, sa.Identity)
			if err != nil {
				return nil, fmt.Errorf("list role bindings: %w", err)
			}
			if !ok {
				res.WouldCreateRoleBinding = true
			}
		} else if res.WouldCreateServiceAccount {
			res.WouldCreateRoleBinding = true
		}
	}

	return res, nil
}
