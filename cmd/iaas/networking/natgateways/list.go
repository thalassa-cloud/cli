package natgateways

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime bool
	region        string
	vpc           string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of NAT gateways",
	Long:    "Get a list of NAT gateways within your organisation",
	Example: "tcloud networking natgateways list\ntcloud networking natgateways list --region us-west-1\ntcloud networking natgateways list --vpc vpc-123 --no-header",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := filters.Filters{}
		if region != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "region",
				Value: region,
			})
		}

		if vpc != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "vpc",
				Value: vpc,
			})
		}

		natgateways, err := client.IaaS().ListNatGateways(cmd.Context(), &iaas.ListNatGatewaysRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(natgateways))
		for _, ngw := range natgateways {

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
	NatGatewaysCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().StringVar(&region, "region", "", "Region of the NAT gateway")
	listCmd.Flags().StringVar(&vpc, "vpc", "", "VPC of the NAT gateway")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")

	// Add completion
	listCmd.RegisterFlagCompletionFunc("region", completeRegion)
	listCmd.RegisterFlagCompletionFunc("vpc", completeVPCID)
}
