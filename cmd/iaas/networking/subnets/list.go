package subnets

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
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
	listVpcFilter     string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of subnets",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := []filters.Filter{}
		if listVpcFilter != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "vpc",
				Value: listVpcFilter,
			})
		}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		subnets, err := client.IaaS().ListSubnets(cmd.Context(), &iaas.ListSubnetsRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(subnets))
		for _, subnet := range subnets {

			vpcName := ""
			if subnet.Vpc != nil {
				vpcName = subnet.Vpc.Name
			} else {
				vpcName = subnet.VpcIdentity
			}

			row := []string{
				subnet.Identity,
				subnet.Name,
				string(subnet.Status),
				vpcName,
				subnet.Cidr,
				formattime.FormatTime(subnet.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labels := []string{}
				for k, v := range subnet.Labels {
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
			headers := []string{"ID", "Name", "Status", "VPC", "CIDR", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	SubnetsCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	getCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	getCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter subnets (format: key1=value1,key2=value2)")
	getCmd.Flags().StringVar(&listVpcFilter, "vpc", "", "Filter by VPC")

	// Add completion
	getCmd.RegisterFlagCompletionFunc("vpc", completion.CompleteVPCID)
}
