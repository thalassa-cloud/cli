package natgateways

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime     bool
	showLabels        bool
	listLabelSelector string
	region            string
	vpc               string
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

		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
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

			row := []string{
				ngw.Identity,
				ngw.Name,
				string(ngw.Status),
				ngw.Vpc.Name,
				regionName,
				ngw.EndpointIP,
				formattime.FormatTime(ngw.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labels := []string{}
				for k, v := range ngw.Labels {
					labels = append(labels, k+"="+v)
				}
				sort.Strings(labels)
				if len(labels) == 0 {
					labels = []string{"-"}
				}
				row = append(row, strings.Join(labels, ","))
			}

			body = append(body, row)
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Status", "VPC", "Region", "IP", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
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
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter NAT gateways (format: key1=value1,key2=value2)")

	// Add completion
	listCmd.RegisterFlagCompletionFunc("region", completeRegion)
	listCmd.RegisterFlagCompletionFunc("vpc", completeVPCID)
}
