package zones

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var exportCmd = &cobra.Command{
	Use:               "export <zone>",
	Short:             "Export a DNS zone as a BIND zone file to stdout",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDnsZoneIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		exported, err := client.DNS().ExportZoneFile(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to export DNS zone: %w", err)
		}

		fmt.Print(exported.ZoneFile)
		if exported.ZoneFile != "" && exported.ZoneFile[len(exported.ZoneFile)-1] != '\n' {
			fmt.Println()
		}
		return nil
	},
}

func init() {
	ZonesCmd.AddCommand(exportCmd)
}
