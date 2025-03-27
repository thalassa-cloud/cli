package nodepools

import (
	"github.com/spf13/cobra"
)

// KubernetesNodePoolsCmd represents the incidents command
var KubernetesNodePoolsCmd = &cobra.Command{
	Use:     "nodepools",
	Aliases: []string{"nodepool", "np"},
	Short:   "Manage Kubernetes NodePools",
	Example: `  # List all nodepools in a cluster
  tcloud kubernetes nodepools list my-cluster

  # Create a new nodepool
  tcloud kubernetes nodepools create my-cluster --name worker-pool --size 3

  # Delete a nodepool
  tcloud kubernetes nodepools delete my-cluster worker-pool`,
}
