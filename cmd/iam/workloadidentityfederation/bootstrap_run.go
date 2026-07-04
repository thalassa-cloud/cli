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

func createBootstrapFederatedIdentity(
	ctx context.Context,
	c *clientiam.Client,
	name, description, providerID, saID, subject string,
	audiences []string,
	scopes []clientiam.AccessCredentialsScope,
	vcs, key string,
) (*clientiam.FederatedIdentity, error) {
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

	issuerSubject, err := resolveBootstrapIssuerAndSubject(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	key := wifResourceKey(opts.VCS, opts.Repository, issuerSubject.subject, issuerSubject.issuer)
	res.WIFKey = key
	res.ProviderSubject = issuerSubject.subject
	res.Issuer = issuerSubject.issuer

	scopes := defaultBootstrapScopes(opts)

	role, err := resolveOrganisationRole(ctx, iamc, opts.RoleRef)
	if err != nil {
		return nil, err
	}
	res.RoleIdentity = role.Identity
	res.RoleSlug = role.Slug

	provider, err := ensureBootstrapProvider(ctx, iamc, opts, issuerSubject.issuer, issuerSubject.k8sClusterBoundProvider, res)
	if err != nil {
		return nil, err
	}

	saName, fiName := bootstrapResourceNames(opts, key)
	sa, err := ensureBootstrapServiceAccount(ctx, iamc, opts, saName, issuerSubject.subject, key, res)
	if err != nil {
		return nil, err
	}

	fiDesc := fmt.Sprintf("Workload identity (%s %s) → %s", opts.VCS, opts.Repository, issuerSubject.subject)
	if _, err := ensureBootstrapFederatedIdentity(ctx, iamc, opts, provider, sa, fiName, fiDesc, issuerSubject.subject, key, scopes, res); err != nil {
		return nil, err
	}

	if err := ensureBootstrapRoleBinding(ctx, iamc, opts, role, sa, key, res); err != nil {
		return nil, err
	}

	return res, nil
}
