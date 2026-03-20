package workloadidentityfederation

// Labels and values for resources created by workload-identity bootstrap, so they can be
// listed or removed consistently (e.g. label selector thalassa.cloud/managed-by=workload-identity-bootstrap).
const (
	LabelManagedBy = "thalassa.cloud/managed-by"
	ValueManagedBy = "workload-identity-bootstrap"

	// LabelWIFKey is a short stable id derived from vcs + repository + JWT subject + OIDC issuer URL.
	LabelWIFKey = "thalassa.cloud/wif-key"

	// LabelWIFVCS is github, gitlab, or kubernetes for filtering.
	LabelWIFVCS = "thalassa.cloud/wif-vcs"

	// LabelKubernetesClusterID is set on the platform-managed federated identity provider for a Thalassa Kubernetes cluster (value = cluster identity).
	LabelKubernetesClusterID = "kubernetes_cluster_id"

	ValueVCSGitHub     = "github"
	ValueVCSGitLab     = "gitlab"
	ValueVCSKubernetes = "kubernetes"
)
