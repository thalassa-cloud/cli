package nodepools

import (
	"github.com/spf13/cobra"
)

// KubernetesNodePoolsCmd represents the nodepools command
var KubernetesNodePoolsCmd = &cobra.Command{
	Use:     "nodepools",
	Aliases: []string{"nodepool", "np"},
	Short:   "Manage Kubernetes NodePools",
	Example: `  # List all nodepools in a cluster
  tcloud kubernetes nodepools list

  # Create a new nodepool
  tcloud kubernetes nodepools create --cluster my-cluster --name worker --machine-type pgp-medium --num-nodes 3 --availability-zone nl-01a

  # Delete a nodepool
  tcloud kubernetes nodepools delete --cluster my-cluster --name worker-pool`,
}

func init() {
	KubernetesNodePoolsCmd.AddCommand(listCmd)
	KubernetesNodePoolsCmd.AddCommand(createCmd)
	KubernetesNodePoolsCmd.AddCommand(updateCmd)
	KubernetesNodePoolsCmd.AddCommand(deleteCmd)
}
