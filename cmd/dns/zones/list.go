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

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List DNS zones",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		zones, err := client.DNS().ListZones(cmd.Context(), &clientdns.ListZonesRequest{})
		if err != nil {
			return fmt.Errorf("failed to list DNS zones: %w", err)
		}

		body := make([][]string, 0, len(zones))
		for _, z := range zones {
			body = append(body, []string{
				z.Identity,
				z.Name,
				z.Slug,
				z.Description,
				formattime.FormatTime(z.CreatedAt.Local(), showExactTime),
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "Description", "Created"}, body)
		}
		return nil
	},
}

func init() {
	ZonesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
