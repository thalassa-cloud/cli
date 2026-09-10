package kubeconfigsessions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	deleteCluster string
	deleteForce   bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete [cluster] <session>",
	Short: "Delete a kubeconfig session",
	Long: `Revoke a kubeconfig session for a Kubernetes cluster.

Provide the cluster as the first argument or with --cluster, and the session identity as the final argument.`,
	Aliases:           []string{"rm", "remove"},
	Args:              cobra.RangeArgs(1, 2),
	ValidArgsFunction: completeDeleteArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		clusterRef := deleteCluster
		sessionIdentity := ""
		switch len(args) {
		case 1:
			sessionIdentity = args[0]
		case 2:
			if clusterRef != "" && clusterRef != args[0] {
				return fmt.Errorf("provide the cluster as an argument or --cluster, not both with different values")
			}
			clusterRef = args[0]
			sessionIdentity = args[1]
		}
		if clusterRef == "" {
			return fmt.Errorf("cluster is required (argument or --cluster)")
		}
		if sessionIdentity == "" {
			return fmt.Errorf("session identity is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		cluster, err := kuberesolve.ResolveKubernetesClusterRef(ctx, client.Kubernetes(), clusterRef)
		if err != nil {
			return err
		}

		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf(
			"Are you sure you want to delete this kubeconfig session?\n  Cluster: %s\n  Session: %s\n",
			cluster.Name, sessionIdentity,
		))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		if err := client.Kubernetes().DeleteKubeconfigSession(ctx, cluster.Identity, sessionIdentity); err != nil {
			return fmt.Errorf("failed to delete kubeconfig session: %w", err)
		}
		fmt.Printf("Deleted kubeconfig session %s\n", sessionIdentity)
		return nil
	},
}

func completeDeleteArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		if deleteCluster != "" {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completion.CompleteKubernetesCluster(cmd, args, toComplete)
	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

func init() {
	KubeconfigSessionsCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteCluster, ClusterFlag, "", "Cluster identity, name, or slug")
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
	_ = deleteCmd.RegisterFlagCompletionFunc(ClusterFlag, completion.CompleteKubernetesCluster)
}
