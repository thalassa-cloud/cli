package projects

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <project>",
	Short:             "Delete a project",
	Aliases:           []string{"rm", "remove"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to delete this project?\n  Project: %s\n", args[0]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		if err := client.Projects().DeleteProject(cmd.Context(), args[0]); err != nil {
			return fmt.Errorf("failed to delete project: %w", err)
		}
		fmt.Printf("Deleted project %s\n", args[0])
		return nil
	},
}

func init() {
	ProjectsCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
}
