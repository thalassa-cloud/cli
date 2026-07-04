package workloadidentityfederation

import (
	"context"
	"fmt"
	"strings"

	clientiam "github.com/thalassa-cloud/client-go/iam"
	"github.com/thalassa-cloud/client-go/thalassa"
)

type bootstrapIssuerSubject struct {
	issuer                  string
	subject                 string
	k8sClusterBoundProvider *clientiam.FederatedIdentityProvider
}

func resolveBootstrapIssuerAndSubject(
	ctx context.Context,
	client thalassa.Client,
	opts BootstrapOptions,
) (bootstrapIssuerSubject, error) {
	var out bootstrapIssuerSubject
	var err error

	switch opts.VCS {
	case ValueVCSGitHub:
		out.issuer = GitHubActionsIssuer
		out.subject, err = BuildGitHubSubject(opts.Repository, opts.RefKind, opts.Ref)
	case ValueVCSGitLab:
		out.issuer = normalizeIssuer(opts.GitLabIssuer)
		if out.issuer == "" {
			out.issuer = "https://gitlab.com"
		}
		out.subject, err = BuildGitLabSubject(opts.Repository, opts.GitLabRefType, opts.Ref)
	case ValueVCSKubernetes:
		out, err = resolveKubernetesBootstrapIssuer(ctx, client, opts)
	default:
		return out, fmt.Errorf("unsupported --vcs %q (use github, gitlab, or kubernetes)", opts.VCS)
	}
	if err != nil {
		return out, err
	}
	if len(opts.TrustedAudiences) == 0 {
		return out, fmt.Errorf("at least one trusted audience is required")
	}
	return out, nil
}

func resolveKubernetesBootstrapIssuer(
	ctx context.Context,
	client thalassa.Client,
	opts BootstrapOptions,
) (bootstrapIssuerSubject, error) {
	var out bootstrapIssuerSubject
	var err error
	explicitIssuer := normalizeIssuer(opts.KubernetesIssuer)

	if strings.TrimSpace(opts.KubernetesClusterRef) != "" {
		cluster, err := resolveKubernetesCluster(ctx, client.Kubernetes(), opts.KubernetesClusterRef)
		if err != nil {
			return out, err
		}
		p, err := findProviderByKubernetesClusterID(ctx, client.IAM(), cluster.Identity)
		if err != nil {
			return out, fmt.Errorf("list federated identity providers for cluster: %w", err)
		}
		if p == nil {
			return out, fmt.Errorf("no federated identity provider for kubernetes cluster %q (identity %s): expected exactly one provider with label %s=%s (cluster OIDC not provisioned yet?)",
				cluster.Name, cluster.Identity, LabelKubernetesClusterID, cluster.Identity)
		}
		out.k8sClusterBoundProvider = p
		out.issuer = normalizeIssuer(p.ProviderIssuer)
		if out.issuer == "" {
			return out, fmt.Errorf("federated identity provider %s for cluster %q has an empty issuer URL", p.Identity, cluster.Name)
		}
		if explicitIssuer != "" && explicitIssuer != out.issuer {
			return out, fmt.Errorf("kubernetes: --issuer %q does not match cluster OIDC provider issuer %q (omit --issuer when using --cluster)", opts.KubernetesIssuer, p.ProviderIssuer)
		}
	} else if explicitIssuer != "" {
		out.issuer = explicitIssuer
	} else {
		return out, fmt.Errorf("kubernetes: set --kubernetes-cluster (Thalassa-managed cluster) and/or --kubernetes-issuer (self-managed or custom service-account issuer)")
	}

	out.subject, err = BuildKubernetesSubject(opts.Repository)
	return out, err
}

func defaultBootstrapProviderName(vcs, issuer string) string {
	switch vcs {
	case ValueVCSGitHub:
		return "GitHub Actions OIDC"
	case ValueVCSGitLab:
		if h := issuerURLHostname(issuer); h != "" {
			return h
		}
		return "GitLab CI OIDC"
	case ValueVCSKubernetes:
		if h := issuerURLHostname(issuer); h != "" {
			return h
		}
		return "Kubernetes service account OIDC"
	default:
		return "OIDC"
	}
}

func ensureBootstrapProvider(
	ctx context.Context,
	iamc *clientiam.Client,
	opts BootstrapOptions,
	issuer string,
	k8sClusterBoundProvider *clientiam.FederatedIdentityProvider,
	res *BootstrapResult,
) (*clientiam.FederatedIdentityProvider, error) {
	if k8sClusterBoundProvider != nil {
		res.ProviderIdentity = k8sClusterBoundProvider.Identity
		return k8sClusterBoundProvider, nil
	}

	provider, err := findProviderByIssuer(ctx, iamc, issuer)
	if err != nil {
		return nil, fmt.Errorf("list federated identity providers: %w", err)
	}
	if provider != nil {
		res.ProviderIdentity = provider.Identity
		return provider, nil
	}

	if opts.DryRun {
		res.WouldCreateProvider = true
		res.ProviderIdentity = "(dry-run: would create provider)"
		return nil, nil
	}

	pName := opts.ProviderDisplayName
	if pName == "" {
		pName = defaultBootstrapProviderName(opts.VCS, issuer)
	}
	pDesc := opts.ProviderDescription
	if pDesc == "" {
		pDesc = fmt.Sprintf("OIDC issuer %s (created by tcloud iam workload-identity-federation bootstrap <github|gitlab|kubernetes>)", issuer)
	}
	p, err := createBootstrapProvider(ctx, iamc, pName, pDesc, issuer, opts.VCS)
	if err != nil {
		return nil, fmt.Errorf("create federated identity provider: %w", err)
	}
	res.CreatedProvider = true
	res.ProviderIdentity = p.Identity
	return p, nil
}

func ensureBootstrapServiceAccount(
	ctx context.Context,
	iamc *clientiam.Client,
	opts BootstrapOptions,
	saName, subject, key string,
	res *BootstrapResult,
) (*clientiam.ServiceAccount, error) {
	sa, err := findServiceAccountByWIFKey(ctx, iamc, opts.VCS, key)
	if err != nil {
		return nil, fmt.Errorf("list service accounts: %w", err)
	}
	if sa != nil {
		res.ServiceAccountIdentity = sa.Identity
		res.ServiceAccountSlug = sa.Slug
		return sa, nil
	}

	if opts.DryRun {
		res.WouldCreateServiceAccount = true
		res.ServiceAccountIdentity = "(dry-run: would create service account)"
		return nil, nil
	}

	desc := fmt.Sprintf("Workload identity (%s %s); JWT sub: %s", opts.VCS, opts.Repository, subject)
	sa, err = createBootstrapServiceAccount(ctx, iamc, saName, desc, opts.VCS, key, opts.Repository, subject)
	if err != nil {
		return nil, fmt.Errorf("create service account: %w", err)
	}
	res.CreatedServiceAccount = true
	res.ServiceAccountIdentity = sa.Identity
	res.ServiceAccountSlug = sa.Slug
	return sa, nil
}

func ensureBootstrapFederatedIdentity(
	ctx context.Context,
	iamc *clientiam.Client,
	opts BootstrapOptions,
	provider *clientiam.FederatedIdentityProvider,
	sa *clientiam.ServiceAccount,
	fiName, fiDesc, subject, key string,
	scopes []clientiam.AccessCredentialsScope,
	res *BootstrapResult,
) (*clientiam.FederatedIdentity, error) {
	fi, err := findFederatedIdentity(ctx, iamc, opts.VCS, key)
	if err != nil {
		return nil, fmt.Errorf("list federated identities: %w", err)
	}
	if fi == nil {
		if opts.DryRun {
			res.WouldCreateFederatedIdentity = true
			res.FederatedIdentityIdentity = "(dry-run: would create federated identity)"
			return nil, nil
		}
		if provider == nil || sa == nil {
			return nil, fmt.Errorf("internal: missing provider or service account after provisioning")
		}
		fi, err = createBootstrapFederatedIdentity(ctx, iamc, fiName, fiDesc, provider.Identity, sa.Identity, subject, opts.TrustedAudiences, scopes, opts.VCS, key)
		if err != nil {
			return nil, fmt.Errorf("create federated identity: %w", err)
		}
		res.CreatedFederatedIdentity = true
		res.FederatedIdentityIdentity = fi.Identity
		return fi, nil
	}

	res.FederatedIdentityIdentity = fi.Identity
	if !federatedIdentityNeedsBootstrapReconcile(fi, fiName, fiDesc, subject, opts.VCS, key, scopes, opts.TrustedAudiences) {
		return fi, nil
	}
	if opts.DryRun {
		res.WouldUpdateFederatedIdentity = true
		return fi, nil
	}

	_, err = iamc.UpdateFederatedIdentity(ctx, fi.Identity, clientiam.UpdateFederatedIdentityRequest{
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
	return fi, nil
}

func ensureBootstrapRoleBinding(
	ctx context.Context,
	iamc *clientiam.Client,
	opts BootstrapOptions,
	role *clientiam.OrganisationRole,
	sa *clientiam.ServiceAccount,
	key string,
	res *BootstrapResult,
) error {
	if opts.DryRun {
		if sa != nil {
			ok, err := hasRoleBindingForServiceAccount(ctx, iamc, role.Identity, sa.Identity)
			if err != nil {
				return fmt.Errorf("list role bindings: %w", err)
			}
			if !ok {
				res.WouldCreateRoleBinding = true
			}
		} else if res.WouldCreateServiceAccount {
			res.WouldCreateRoleBinding = true
		}
		return nil
	}

	if sa == nil {
		return nil
	}

	ok, err := hasRoleBindingForServiceAccount(ctx, iamc, role.Identity, sa.Identity)
	if err != nil {
		return fmt.Errorf("list role bindings: %w", err)
	}
	if ok {
		return nil
	}

	_, err = createRoleBindingForSA(ctx, iamc, role, sa, opts.VCS, key)
	if err != nil {
		return fmt.Errorf("create role binding: %w", err)
	}
	res.CreatedRoleBinding = true
	return nil
}

func bootstrapResourceNames(opts BootstrapOptions, key string) (saName, fiName string) {
	saName = fmt.Sprintf("wif-%s-%s", opts.VCS, key)
	fiName = fmt.Sprintf("wif-%s-%s-fi", opts.VCS, key)
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
	return saName, fiName
}

func defaultBootstrapScopes(opts BootstrapOptions) []clientiam.AccessCredentialsScope {
	if len(opts.AllowedScopes) > 0 {
		return opts.AllowedScopes
	}
	return []clientiam.AccessCredentialsScope{
		clientiam.AccessCredentialsScopeAPIRead,
		clientiam.AccessCredentialsScopeAPIWrite,
	}
}
