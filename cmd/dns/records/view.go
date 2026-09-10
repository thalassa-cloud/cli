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
)

var viewCmd = &cobra.Command{
	Use:     "view <record>",
	Aliases: []string{"get", "show"},
	Short:   "View a DNS record",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		record, err := client.DNS().GetRecord(cmd.Context(), zoneIdentity, args[0])
		if err != nil {
			return fmt.Errorf("failed to get DNS record: %w", err)
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
	RecordsCmd.AddCommand(viewCmd)
	viewCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	viewCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	viewCmd.Flags().StringVar(&zoneIdentity, "zone", "", "DNS zone identity")
	_ = viewCmd.MarkFlagRequired("zone")
	_ = viewCmd.RegisterFlagCompletionFunc("zone", completion.CompleteDnsZoneIdentity)
}
