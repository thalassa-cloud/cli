package roles

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/kubernetes/iam/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <role>",
	Short:             "Delete a custom Kubernetes cluster IAM role",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKubernetesClusterRoleIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		role, err := kuberesolve.ResolveKubernetesClusterRoleRef(ctx, client.Kubernetes(), args[0])
		if err != nil {
			return err
		}
		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to delete this Kubernetes cluster role?\n  Role: %s\n", role.Name))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := client.Kubernetes().DeleteClusterRole(ctx, role.Identity); err != nil {
			return fmt.Errorf("failed to delete role: %w", err)
		}
		fmt.Printf("Deleted role %s\n", role.Name)
		return nil
	},
}

func init() {
	RolesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
}
