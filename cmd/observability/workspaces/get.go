package workspaces

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var getExactTime bool

var getCmd = &cobra.Command{
	Use:               "get <workspace>",
	Aliases:           []string{"view", "show", "describe"},
	Short:             "Show an observability workspace",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteObservabilityWorkspaceID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ws, err := client.Observability().GetObservabilityWorkspace(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get observability workspace: %w", err)
		}

		printWorkspaceDetails(ws)
		fmt.Printf("Created:     %s\n", formattime.FormatTime(ws.CreatedAt.Local(), getExactTime))
		fmt.Printf("Updated:     %s\n", formattime.FormatTime(ws.UpdatedAt.Local(), getExactTime))
		return nil
	},
}

func init() {
	WorkspacesCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&getExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
