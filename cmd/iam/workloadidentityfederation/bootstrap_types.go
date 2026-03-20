package workloadidentityfederation

import clientiam "github.com/thalassa-cloud/client-go/iam"

// BootstrapOptions configures workload-identity-federation bootstrap.
type BootstrapOptions struct {
	VCS              string // ValueVCSGitHub, ValueVCSGitLab, or ValueVCSKubernetes
	Repository       string // owner/repo (GitHub), group/project (GitLab), or namespace/sa (Kubernetes)
	RefKind          RefKind
	Ref              string // branch name, tag name, or environment name (GitHub)
	GitLabRefType    string // branch, tag, etc. (GitLab id_token sub ref_type segment)
	RoleRef          string // organisation role identity, slug, or name
	TrustedAudiences []string
	AllowedScopes    []clientiam.AccessCredentialsScope

	GitLabIssuer string // default https://gitlab.com

	// KubernetesIssuer is the OIDC issuer URL for service account tokens (--service-account-issuer).
	// Optional if KubernetesClusterRef is set and the issuer matches the cluster API URL (typical Thalassa-managed clusters).
	KubernetesIssuer string
	// KubernetesClusterRef is a Thalassa cluster identity, slug, or name used to resolve the default issuer.
	KubernetesClusterRef string

	ProviderDisplayName string
	ProviderDescription string

	// ResourceName, if set, is the Thalassa service account Name and the base for the federated identity Name ({name}-fi).
	// When empty, names default to wif-<platform>-<wif-key> and wif-<platform>-<wif-key>-fi.
	ResourceName string

	DryRun bool
}

// BootstrapResult summarises what bootstrap did or would do.
type BootstrapResult struct {
	WIFKey                    string
	Issuer                    string
	ProviderSubject           string
	ProviderIdentity          string
	ServiceAccountIdentity    string
	ServiceAccountSlug        string
	RoleIdentity              string
	RoleSlug                  string
	FederatedIdentityIdentity string

	CreatedProvider              bool
	CreatedServiceAccount        bool
	CreatedRoleBinding           bool
	CreatedFederatedIdentity     bool
	UpdatedFederatedIdentity     bool // existing FI reconciled to bootstrap (scopes, audiences, labels, etc.)
	WouldCreateProvider          bool
	WouldCreateServiceAccount    bool
	WouldCreateRoleBinding       bool
	WouldCreateFederatedIdentity bool
	WouldUpdateFederatedIdentity bool // dry-run: existing FI would be reconciled
}
