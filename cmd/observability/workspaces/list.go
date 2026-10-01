package workspaces

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientobs "github.com/thalassa-cloud/client-go/observability"
)

var noHeader bool

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List observability workspaces",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		workspaces, err := client.Observability().ListObservabilityWorkspaces(cmd.Context(), &clientobs.ListObservabilityWorkspacesRequest{})
		if err != nil {
			return fmt.Errorf("failed to list observability workspaces: %w", err)
		}

		body := make([][]string, 0, len(workspaces))
		for _, ws := range workspaces {
			body = append(body, []string{
				ws.Identity,
				ws.Name,
				string(ws.Status),
				regionName(ws),
				retentionString(ws.RetentionDays),
				boolYesNo(ws.PrometheusEnabled),
				boolYesNo(ws.LokiEnabled),
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Status", "Region", "Retention", "Prometheus", "Loki"}, body)
		}
		return nil
	},
}

func init() {
	WorkspacesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
