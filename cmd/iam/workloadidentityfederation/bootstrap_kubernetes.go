package workloadidentityfederation

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
)

var (
	flagK8sNamespace string
	flagK8sSAName    string
	flagK8sCluster   string
	flagK8sIssuer    string
)

var bootstrapKubernetesCmd = &cobra.Command{
	Use:     "kubernetes",
	Short:   "Bootstrap workload identity for Kubernetes service accounts",
	Aliases: []string{"k8s"},
	Long: `Binds system:serviceaccount:<namespace>:<name> to a Thalassa service account via a federated identity.

Thalassa clusters: pass --cluster to resolve the cluster and use the platform-managed federated identity
provider labelled kubernetes_cluster_id=<cluster identity>. If that provider is missing, bootstrap fails
(you must wait for cluster OIDC integration).

Other clusters: pass --issuer with the same URL as kube-apiserver --service-account-issuer; the CLI
creates the federated identity provider if it does not exist yet.`,
	Example: `  # Thalassa-managed cluster (OIDC provider from label kubernetes_cluster_id)
  tcloud iam workload-identity-federation bootstrap kubernetes --cluster my-cluster-slug \
    --namespace default --service-account my-app --role deployer

  # Self-managed / custom issuer
  tcloud iam workload-identity-federation bootstrap kubernetes --issuer https://k8s.example.com \
    --namespace cicd --service-account terraform --role deployer`,
	Args: cobra.NoArgs,
	RunE: runBootstrapKubernetes,
}

func runBootstrapKubernetes(cmd *cobra.Command, _ []string) error {
	ns := strings.TrimSpace(flagK8sNamespace)
	sa := strings.TrimSpace(flagK8sSAName)
	if ns == "" || sa == "" {
		return fmt.Errorf("--namespace and --service-account are required")
	}
	if strings.TrimSpace(flagK8sIssuer) == "" && strings.TrimSpace(flagK8sCluster) == "" {
		return fmt.Errorf("set --cluster (Thalassa) and/or --issuer (self-managed)")
	}
	repo := ns + "/" + sa
	return executeBootstrap(cmd, BootstrapOptions{
		VCS:                  ValueVCSKubernetes,
		Repository:           repo,
		KubernetesIssuer:     strings.TrimSpace(flagK8sIssuer),
		KubernetesClusterRef: strings.TrimSpace(flagK8sCluster),
	})
}

func init() {
	f := bootstrapKubernetesCmd.Flags()
	f.StringVar(&flagK8sNamespace, "namespace", "", "Kubernetes namespace for the workload service account (required)")
	f.StringVar(&flagK8sSAName, "service-account", "", "Kubernetes service account name (required)")
	f.StringVar(&flagK8sCluster, "cluster", "", "Thalassa cluster identity, slug, or name (uses federated identity provider with label kubernetes_cluster_id=<cluster identity>)")
	f.StringVar(&flagK8sIssuer, "issuer", "", "kube-apiserver service-account issuer URL (required without --cluster); with --cluster must match the cluster OIDC provider or omit")
	_ = bootstrapKubernetesCmd.RegisterFlagCompletionFunc("cluster", completion.CompleteKubernetesCluster)
}
