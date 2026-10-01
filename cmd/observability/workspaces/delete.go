package workspaces

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <workspace>",
	Aliases:           []string{"rm", "remove"},
	Short:             "Delete an observability workspace",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteObservabilityWorkspaceID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to delete this observability workspace?\n  Workspace: %s\n", args[0]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		if err := client.Observability().DeleteObservabilityWorkspace(cmd.Context(), args[0]); err != nil {
			return fmt.Errorf("failed to delete observability workspace: %w", err)
		}
		fmt.Printf("Deleted observability workspace %s\n", args[0])
		return nil
	},
}

func init() {
	WorkspacesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
}
