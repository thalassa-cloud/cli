package vpcpeering

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
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List VPC peering connections",
	Aliases: []string{"l", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := []filters.Filter{}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		connections, err := client.IaaS().ListVpcPeeringConnections(cmd.Context(), &iaas.ListVpcPeeringConnectionsRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(connections))
		for _, conn := range connections {
			requesterVPC := "-"
			if conn.RequesterVpc != nil {
				requesterVPC = conn.RequesterVpc.Name
			}

			accepterVPC := "-"
			if conn.AccepterVpc != nil {
				accepterVPC = conn.AccepterVpc.Name
			}

			row := []string{
				conn.Identity,
				conn.Name,
				string(conn.Status),
				requesterVPC,
				accepterVPC,
				formattime.FormatTime(conn.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labelStrs := []string{}
				for k, v := range conn.Labels {
					labelStrs = append(labelStrs, k+"="+v)
				}
				sort.Strings(labelStrs)
				if len(labelStrs) == 0 {
					labelStrs = []string{"-"}
				}
				row = append(row, strings.Join(labelStrs, ","))
			}

			body = append(body, row)
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Status", "Requester VPC", "Accepter VPC", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	VpcPeeringCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter connections (format: key1=value1,key2=value2)")
}
