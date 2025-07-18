package vpcs

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime bool
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of vpcs",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		vpcs, err := client.IaaS().ListVpcs(cmd.Context(), &iaas.ListVpcsRequest{})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(vpcs))
		for _, vpc := range vpcs {
			regionName := ""
			if vpc.CloudRegion != nil {
				regionName = vpc.CloudRegion.Name
				if regionName == "" {
					regionName = vpc.CloudRegion.Slug
				}
				if regionName == "" {
					regionName = vpc.CloudRegion.Identity
				}
			}

			body = append(body, []string{
				vpc.Identity,
				vpc.Name,
				regionName,

				strings.Join(vpc.CIDRs, ", "),
				formattime.FormatTime(vpc.CreatedAt.Local(), showExactTime),
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Region", "CIDRs", "Age"}, body)
		}
		return nil
	},
}

func init() {
	VpcsCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
}
