package securitygroups

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
	listVpcFilter string
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

			body = append(body, []string{
				sg.Identity,
				sg.Name,
				vpcName,
				string(sg.Status),
				ingressCount,
				egressCount,
				fmt.Sprintf("%t", sg.AllowSameGroupTraffic),
				formattime.FormatTime(sg.CreatedAt.Local(), showExactTime),
			})
		}
		if len(body) == 0 {
			fmt.Println("No security groups found")
			return nil
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "VPC", "Status", "Ingress Rules", "Egress Rules", "Allow Same Group", "Age"}, body)
		}
		return nil
	},
}

func init() {
	SecurityGroupsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().StringVar(&listVpcFilter, "vpc", "", "Filter by VPC")

	// Add completion
	listCmd.RegisterFlagCompletionFunc("vpc", completeVPCID)
}
