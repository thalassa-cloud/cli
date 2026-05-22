package loadbalancers

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	noHeader          bool
	showExactTime     bool
	showLabels        bool
	listLabelSelector string
	listRegion        string
	listVpc           string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List load balancers",
	Long:    "List load balancers within your organisation.",
	Example: "tcloud networking loadbalancers list\ntcloud networking loadbalancers list --region us-west-1\ntcloud networking loadbalancers list --vpc vpc-123 --no-header",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := filters.Filters{}
		if listRegion != "" {
			f = append(f, &filters.FilterKeyValue{Key: "region", Value: listRegion})
		}
		if listVpc != "" {
			f = append(f, &filters.FilterKeyValue{Key: "vpc", Value: listVpc})
		}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		loadbalancers, err := client.IaaS().ListLoadbalancers(cmd.Context(), &iaas.ListLoadbalancersRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(loadbalancers))
		for _, lb := range loadbalancers {
			regionName := ""
			if lb.Vpc != nil && lb.Vpc.CloudRegion != nil {
				regionName = lb.Vpc.CloudRegion.Name
				if regionName == "" {
					regionName = lb.Vpc.CloudRegion.Identity
				}
				if regionName == "" {
					regionName = lb.Vpc.CloudRegion.Slug
				}
			}

			vpcName := ""
			if lb.Vpc != nil {
				vpcName = lb.Vpc.Name
			}

			ips := append(lb.ExternalIpAddresses, lb.InternalIpAddresses...)
			row := []string{
				lb.Identity,
				lb.Name,
				string(lb.Status),
				vpcName,
				regionName,
				joinStrings(ips),
				fmt.Sprintf("%d", len(lb.LoadbalancerListeners)),
				formattime.FormatTime(lb.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labelPairs := make([]string, 0, len(lb.Labels))
				for k, v := range lb.Labels {
					labelPairs = append(labelPairs, k+"="+v)
				}
				sort.Strings(labelPairs)
				row = append(row, joinStrings(labelPairs))
			}

			body = append(body, row)
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Status", "VPC", "Region", "IPs", "Listeners", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	LoadbalancersCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().StringVar(&listRegion, "region", "", "Filter by region")
	listCmd.Flags().StringVar(&listVpc, "vpc", "", "Filter by VPC")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector (format: key1=value1,key2=value2)")

	listCmd.RegisterFlagCompletionFunc("region", completeRegion)
	listCmd.RegisterFlagCompletionFunc("vpc", completeVPCID)
}
