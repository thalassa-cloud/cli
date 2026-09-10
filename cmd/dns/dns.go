package dns

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/dns/records"
	"github.com/thalassa-cloud/cli/cmd/dns/zones"
)

// DnsCmd manages DNS zones and records.
var DnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage DNS zones and records (beta)",
	Long: `Manage DNS zones and records within the Thalassa Cloud Platform.

Note: This command is in beta.`,
}

func init() {
	DnsCmd.AddCommand(zones.ZonesCmd)
	DnsCmd.AddCommand(records.RecordsCmd)
}
