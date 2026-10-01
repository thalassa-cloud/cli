package quicklaunch

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var logsCmd = &cobra.Command{
	Use:               "logs <identity>",
	Short:             "Show logs for a quick-launch job",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteQuickLaunchIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		logs, err := client.QuickLaunch().GetQuickLaunchLogs(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get quick-launch logs: %w", err)
		}
		if len(logs) == 0 {
			_, err = fmt.Fprintln(os.Stdout, "(no logs)")
			return err
		}
		if _, err = os.Stdout.Write(logs); err != nil {
			return err
		}
		if logs[len(logs)-1] != '\n' {
			_, err = fmt.Fprintln(os.Stdout)
			return err
		}
		return nil
	},
}

func init() {
	QuickLaunchCmd.AddCommand(logsCmd)
}
