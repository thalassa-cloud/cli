package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/kubernetes/connect"
	"github.com/thalassa-cloud/cli/cmd/kubernetes/kubernetesversions"
	"github.com/thalassa-cloud/cli/cmd/kubernetes/nodepools"
)

var KubernetesCmd = &cobra.Command{
	Use:     "kubernetes",
	Aliases: []string{"kube", "k8s", "k"},
	Short:   "Manage Kubernetes clusters, node pools and more services related to Kubernetes",
	Long:    "Kubernetes commands to manage your Kubernetes clusters and node pools within the Thalassa Cloud Platform",
}

func init() {
	KubernetesCmd.AddCommand(kubernetesversions.KubernetesVersionsCmd)
	KubernetesCmd.AddCommand(nodepools.KubernetesNodePoolsCmd)
	KubernetesCmd.AddCommand(connect.KubernetesConnectCmd)
}
