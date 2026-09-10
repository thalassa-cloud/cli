package zones

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientdns "github.com/thalassa-cloud/client-go/dns"
)

var (
	importFile    string
	importReplace bool
)

var importCmd = &cobra.Command{
	Use:               "import <zone>",
	Short:             "Import DNS records from a BIND zone file",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDnsZoneIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		data, err := os.ReadFile(importFile)
		if err != nil {
			return fmt.Errorf("failed to read zone file: %w", err)
		}

		result, err := client.DNS().ImportZoneFile(cmd.Context(), args[0], clientdns.ImportDnsZoneFileRequest{
			ZoneFile:        string(data),
			ReplaceExisting: importReplace,
		})
		if err != nil {
			return fmt.Errorf("failed to import DNS zone: %w", err)
		}

		body := [][]string{{
			fmt.Sprintf("%d", result.Created),
			fmt.Sprintf("%d", result.Updated),
			fmt.Sprintf("%d", result.Deleted),
			fmt.Sprintf("%d", result.Skipped),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Created", "Updated", "Deleted", "Skipped"}, body)
		}
		return nil
	},
}

func init() {
	ZonesCmd.AddCommand(importCmd)
	importCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	importCmd.Flags().StringVar(&importFile, "file", "", "Path to BIND zone file")
	importCmd.Flags().BoolVar(&importReplace, "replace", false, "Replace existing records that conflict")
	_ = importCmd.MarkFlagRequired("file")
}
