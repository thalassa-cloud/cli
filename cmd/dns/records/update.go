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

var (
	updateTTL    int
	updateValues []string
)

var updateCmd = &cobra.Command{
	Use:   "update <record>",
	Short: "Update a DNS record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if !cmd.Flags().Changed("value") {
			return fmt.Errorf("--value is required")
		}

		req := clientdns.UpdateDnsRecordRequest{
			Values: updateValues,
		}
		if cmd.Flags().Changed("ttl") {
			req.TTL = updateTTL
		}

		record, err := client.DNS().UpdateRecord(cmd.Context(), zoneIdentity, args[0], req)
		if err != nil {
			return fmt.Errorf("failed to update DNS record: %w", err)
		}

		body := [][]string{{
			record.Identity,
			record.Name,
			string(record.Type),
			fmt.Sprintf("%d", record.TTL),
			strings.Join(record.Values, ","),
			formattime.FormatTime(record.CreatedAt.Local(), showExactTime),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Type", "TTL", "Values", "Created"}, body)
		}
		return nil
	},
}

func init() {
	RecordsCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	updateCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	updateCmd.Flags().StringVar(&zoneIdentity, "zone", "", "DNS zone identity")
	updateCmd.Flags().IntVar(&updateTTL, "ttl", 0, "Record TTL in seconds")
	updateCmd.Flags().StringSliceVar(&updateValues, "value", nil, "Record value (repeatable)")
	_ = updateCmd.MarkFlagRequired("zone")
	_ = updateCmd.MarkFlagRequired("value")
	_ = updateCmd.RegisterFlagCompletionFunc("zone", completion.CompleteDnsZoneIdentity)
}
