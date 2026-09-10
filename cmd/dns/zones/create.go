package zones

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientdns "github.com/thalassa-cloud/client-go/dns"
)

var (
	createName        string
	createDescription string
	createLabels      []string
	createAnnotations []string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DNS zone",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		zone, err := client.DNS().CreateZone(cmd.Context(), clientdns.CreateDnsZoneRequest{
			ZoneName:    createName,
			Description: createDescription,
			Labels:      shared.KeyValuePairsToMap(createLabels),
			Annotations: shared.KeyValuePairsToMap(createAnnotations),
		})
		if err != nil {
			return fmt.Errorf("failed to create DNS zone: %w", err)
		}

		body := [][]string{{
			zone.Identity,
			zone.Name,
			zone.Slug,
			formattime.FormatTime(zone.CreatedAt.Local(), showExactTime),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "Created"}, body)
		}
		return nil
	},
}

func init() {
	ZonesCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	createCmd.Flags().StringVar(&createName, "name", "", "Zone name (e.g. example.com)")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Zone description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	_ = createCmd.MarkFlagRequired("name")
}
