package zones

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientdns "github.com/thalassa-cloud/client-go/dns"
)

var (
	updateDescription string
	updateLabels      []string
	updateAnnotations []string
)

var updateCmd = &cobra.Command{
	Use:               "update <zone>",
	Short:             "Update a DNS zone",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDnsZoneIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req := clientdns.UpdateDnsZoneRequest{
			Description: updateDescription,
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = shared.KeyValuePairsToMap(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = shared.KeyValuePairsToMap(updateAnnotations)
		}

		zone, err := client.DNS().UpdateZone(cmd.Context(), args[0], req)
		if err != nil {
			return fmt.Errorf("failed to update DNS zone: %w", err)
		}

		body := [][]string{{
			zone.Identity,
			zone.Name,
			zone.Slug,
			zone.Description,
			formattime.FormatTime(zone.CreatedAt.Local(), showExactTime),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "Description", "Created"}, body)
		}
		return nil
	},
}

func init() {
	ZonesCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	updateCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Zone description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Labels as key=value (repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
}
