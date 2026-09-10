package kubeconfigsessions

import (
	"github.com/spf13/cobra"
)

// KubeconfigSessionsCmd manages Kubernetes kubeconfig sessions.
var KubeconfigSessionsCmd = &cobra.Command{
	Use:     "kubeconfig-sessions",
	Aliases: []string{"kubeconfigsessions", "kcs"},
	Short:   "Manage Kubernetes kubeconfig sessions",
	Long:    "List and revoke active kubeconfig sessions for a Kubernetes cluster.",
	Example: `  # List kubeconfig sessions for a cluster
  tcloud kubernetes kubeconfig-sessions list my-cluster

  # Delete a kubeconfig session
  tcloud kubernetes kubeconfig-sessions delete my-cluster sess-abc123 --force`,
}
