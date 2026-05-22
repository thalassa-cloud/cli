package targetgroups

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
	listVpc           string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List target groups",
	Long:    "List load balancer target groups within your organisation.",
	Example: "tcloud networking target-groups list\ntcloud networking target-groups list --vpc vpc-123",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := filters.Filters{}
		if listVpc != "" {
			f = append(f, &filters.FilterKeyValue{Key: "vpc", Value: listVpc})
		}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		targetGroups, err := client.IaaS().ListTargetGroups(cmd.Context(), &iaas.ListTargetGroupsRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(targetGroups))
		for _, tg := range targetGroups {
			vpcName := ""
			if tg.Vpc != nil {
				vpcName = tg.Vpc.Name
			}

			policy := "-"
			if tg.LoadbalancingPolicy != nil {
				policy = string(*tg.LoadbalancingPolicy)
			}

			row := []string{
				tg.Identity,
				tg.Name,
				vpcName,
				fmt.Sprintf("%d", tg.TargetPort),
				string(tg.Protocol),
				policy,
				fmt.Sprintf("%d", len(tg.LoadbalancerTargetGroupAttachments)),
				formattime.FormatTime(tg.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labelPairs := make([]string, 0, len(tg.Labels))
				for k, v := range tg.Labels {
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
			headers := []string{"ID", "Name", "VPC", "Port", "Protocol", "Policy", "Targets", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	TargetGroupsCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().StringVar(&listVpc, "vpc", "", "Filter by VPC")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector (format: key1=value1,key2=value2)")

	listCmd.RegisterFlagCompletionFunc("vpc", completeVPCID)
}
