package natgateways

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/pkg/thalassa"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime bool
	region        string
	vpc           string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of natgateways",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			return err
		}
		natgateways, err := client.IaaS().ListNatGateways(cmd.Context())
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(natgateways))
		for _, ngw := range natgateways {
			if region != "" && ngw.Vpc.CloudRegion != nil && (ngw.Vpc.CloudRegion.Name != region && ngw.Vpc.CloudRegion.Identity != region && ngw.Vpc.CloudRegion.Slug != region) {
				continue
			}
			if vpc != "" && (ngw.Vpc.Name != vpc && ngw.Vpc.Identity != vpc && ngw.Vpc.Slug != vpc) {
				continue
			}

			regionName := ""
			if ngw.Vpc.CloudRegion != nil {
				regionName = ngw.Vpc.CloudRegion.Name
				if regionName == "" {
					regionName = ngw.Vpc.CloudRegion.Identity
				}
				if regionName == "" {
					regionName = ngw.Vpc.CloudRegion.Slug
				}
			}

			body = append(body, []string{
				ngw.Identity,
				ngw.Name,
				ngw.Vpc.Name,
				regionName,
				ngw.EndpointIP,
				formattime.FormatTime(ngw.CreatedAt.Local(), showExactTime),
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "VPC", "Region", "IP", "Age"}, body)
		}
		return nil
	},
}

func init() {
	NatGatewaysCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	getCmd.Flags().StringVar(&region, "region", "", "Region of the natgateway")
	getCmd.Flags().StringVar(&vpc, "vpc", "", "VPC of the natgateway")
}
