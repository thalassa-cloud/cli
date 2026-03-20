package roles

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <role>",
	Short:             "Delete a custom organisation role",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMOrganisationRoleIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to delete this organisation role?\n  Role: %s\n", args[0]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := client.IAM().DeleteOrganisationRole(ctx, args[0]); err != nil {
			return fmt.Errorf("failed to delete role: %w", err)
		}
		fmt.Printf("Deleted role %s\n", args[0])
		return nil
	},
}

func init() {
	RolesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
	deleteCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
