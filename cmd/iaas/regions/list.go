package regions

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

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
	Short:   "Get a list of regions",
	Long:    "Get a list of regions to manage your regions within the Thalassa Cloud Platform. This command will list all the regions available to you.",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(regions))
		for _, region := range regions {
			zones := []string{}
			for _, zone := range region.Zones {
				zones = append(zones, zone.Name)
			}
			body = append(body, []string{
				region.Identity,
				region.Name,
				region.Slug,
				strings.Join(zones, ","),
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "Zones"}, body)
		}
		return nil
	},
}

func init() {
	RegionsCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
}
