package records

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientdns "github.com/thalassa-cloud/client-go/dns"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List DNS records in a zone",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		records, err := client.DNS().ListRecords(cmd.Context(), zoneIdentity, &clientdns.ListRecordsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list DNS records: %w", err)
		}

		body := make([][]string, 0, len(records))
		for _, r := range records {
			body = append(body, []string{
				r.Identity,
				r.Name,
				string(r.Type),
				fmt.Sprintf("%d", r.TTL),
				strings.Join(r.Values, ","),
				formattime.FormatTime(r.CreatedAt.Local(), showExactTime),
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Type", "TTL", "Values", "Created"}, body)
		}
		return nil
	},
}

func init() {
	RecordsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	listCmd.Flags().StringVar(&zoneIdentity, "zone", "", "DNS zone identity")
	_ = listCmd.MarkFlagRequired("zone")
	_ = listCmd.RegisterFlagCompletionFunc("zone", completion.CompleteDnsZoneIdentity)
}
