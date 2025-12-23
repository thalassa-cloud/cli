package securitygroups

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
	showExactTime    bool
	showLabels       bool
	listLabelSelector string
	listVpcFilter    string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of security groups",
	Long:    "Get a list of security groups within your organisation",
	Example: "tcloud networking security-groups list\ntcloud networking security-groups list --no-header\ntcloud networking security-groups list --exact-time",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := filters.Filters{}
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

		securityGroups, err := client.IaaS().ListSecurityGroups(cmd.Context(), &iaas.ListSecurityGroupsRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(securityGroups))
		for _, sg := range securityGroups {
			vpcName := ""
			if sg.Vpc != nil {
				vpcName = sg.Vpc.Name
			}

			ingressCount := fmt.Sprintf("%d", len(sg.IngressRules))
			egressCount := fmt.Sprintf("%d", len(sg.EgressRules))

			row := []string{
				sg.Identity,
				sg.Name,
				vpcName,
				string(sg.Status),
				ingressCount,
				egressCount,
				fmt.Sprintf("%t", sg.AllowSameGroupTraffic),
				formattime.FormatTime(sg.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labels := []string{}
				for k, v := range sg.Labels {
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
		if len(body) == 0 {
			fmt.Println("No security groups found")
			return nil
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "VPC", "Status", "Ingress Rules", "Egress Rules", "Allow Same Group", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	SecurityGroupsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter security groups (format: key1=value1,key2=value2)")
	listCmd.Flags().StringVar(&listVpcFilter, "vpc", "", "Filter by VPC")

	// Add completion
	listCmd.RegisterFlagCompletionFunc("vpc", completeVPCID)
}
