package machines

import "github.com/spf13/cobra"

// MachinesCmd manages Kubernetes cluster worker machines (nodes).
var MachinesCmd = &cobra.Command{
	Use:     "machines",
	Aliases: []string{"machine", "nodes", "node"},
	Short:   "List and manage Kubernetes cluster machines",
	Long:    "Commands for machines (nodes) that belong to Kubernetes node pools within a cluster.",
	Example: `  # List all machines in a cluster
  tcloud kubernetes machines list --cluster my-cluster

  # List machines in a specific node pool
  tcloud kubernetes machines list --cluster my-cluster --nodepool worker`,
}
