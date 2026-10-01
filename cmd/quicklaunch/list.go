package quicklaunch

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientql "github.com/thalassa-cloud/client-go/quicklaunch"
)

var noHeader bool

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List quick-launch jobs",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		jobs, err := client.QuickLaunch().ListQuickLaunches(cmd.Context(), &clientql.ListQuickLaunchesRequest{})
		if err != nil {
			return fmt.Errorf("failed to list quick-launch jobs: %w", err)
		}

		body := make([][]string, 0, len(jobs))
		for _, ql := range jobs {
			body = append(body, []string{
				ql.Identity,
				ql.Name,
				string(ql.Template),
				ql.Status,
				fmt.Sprintf("%d", len(ql.Resources)),
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Template", "Status", "Resources"}, body)
		}
		return nil
	},
}

func init() {
	QuickLaunchCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
