package quicklaunch

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var getExactTime bool

var getCmd = &cobra.Command{
	Use:               "get <identity>",
	Aliases:           []string{"view", "show", "describe"},
	Short:             "Show a quick-launch job",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteQuickLaunchIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ql, err := client.QuickLaunch().GetQuickLaunch(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get quick-launch job: %w", err)
		}

		printQuickLaunchDetails(ql, getExactTime)
		return nil
	},
}

func init() {
	QuickLaunchCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&getExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
