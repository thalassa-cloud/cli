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
	createName   string
	createType   string
	createTTL    int
	createValues []string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DNS record",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		recordType := clientdns.DnsRecordType(strings.ToUpper(createType))
		switch recordType {
		case clientdns.DnsRecordTypeTXT,
			clientdns.DnsRecordTypeA,
			clientdns.DnsRecordTypeCNAME,
			clientdns.DnsRecordTypeCAA,
			clientdns.DnsRecordTypeAAAA,
			clientdns.DnsRecordTypeMX,
			clientdns.DnsRecordTypeNS,
			clientdns.DnsRecordTypeSRV:
		default:
			return fmt.Errorf("unsupported record type %q (supported: TXT, A, CNAME, CAA, AAAA, MX, NS, SRV)", createType)
		}

		if len(createValues) == 0 {
			return fmt.Errorf("--value is required")
		}

		record, err := client.DNS().CreateRecord(cmd.Context(), zoneIdentity, clientdns.CreateDnsRecordRequest{
			Name:   createName,
			Type:   recordType,
			TTL:    createTTL,
			Values: createValues,
		})
		if err != nil {
			return fmt.Errorf("failed to create DNS record: %w", err)
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
	RecordsCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	createCmd.Flags().StringVar(&zoneIdentity, "zone", "", "DNS zone identity")
	createCmd.Flags().StringVar(&createName, "name", "", "Record name")
	createCmd.Flags().StringVar(&createType, "type", "", "Record type (TXT, A, CNAME, CAA, AAAA, MX, NS, SRV)")
	createCmd.Flags().IntVar(&createTTL, "ttl", 300, "Record TTL in seconds")
	createCmd.Flags().StringSliceVar(&createValues, "value", nil, "Record value (repeatable)")
	_ = createCmd.MarkFlagRequired("zone")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("type")
	_ = createCmd.MarkFlagRequired("value")
	_ = createCmd.RegisterFlagCompletionFunc("zone", completion.CompleteDnsZoneIdentity)
	_ = createCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"TXT", "A", "CNAME", "CAA", "AAAA", "MX", "NS", "SRV"}, cobra.ShellCompDirectiveNoFileComp
	})
}
