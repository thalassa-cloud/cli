package bindings

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
	Use:               "delete <role> <binding>",
	Short:             "Delete a Kubernetes cluster role binding",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completion.CompleteKubernetesClusterRoleThenBinding,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to delete this role binding?\n  Role: %s\n  Binding: %s\n", args[0], args[1]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		role, err := kuberesolve.ResolveKubernetesClusterRoleRef(ctx, client.Kubernetes(), args[0])
		if err != nil {
			return err
		}
		if err := client.Kubernetes().DeleteClusterRoleBinding(ctx, role.Identity, args[1]); err != nil {
			return fmt.Errorf("failed to delete binding: %w", err)
		}
		fmt.Printf("Deleted binding %s from role %s\n", args[1], role.Name)
		return nil
	},
}

func init() {
	BindingsCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
}
