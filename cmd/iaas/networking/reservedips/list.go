package reservedips

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

var (
	showExactTime     bool
	showLabels        bool
	listLabelSelector string
	listRegion        string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List reserved IP addresses",
	Long:    "List reserved IP addresses within your organisation.",
	Example: "tcloud networking reserved-ips list\ntcloud networking reserved-ips list --region nl-ams\ntcloud networking reserved-ips list --selector env=prod",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := filters.Filters{}
		if listRegion != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "region",
				Value: listRegion,
			})
		}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		reservedIPs, err := client.IaaS().ListReservedIPs(cmd.Context(), &iaas.ListReservedIPsRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(reservedIPs))
		for _, rip := range reservedIPs {
			row := []string{
				rip.Identity,
				rip.Name,
				regionName(rip),
				string(rip.Status),
				dashIfEmpty(rip.IPv4Address),
				dashIfEmpty(rip.IPv6Address),
				attachedTo(rip),
			}
			if showLabels {
				labelPairs := make([]string, 0, len(rip.Labels))
				for k, v := range rip.Labels {
					labelPairs = append(labelPairs, k+"="+v)
				}
				sort.Strings(labelPairs)
				if len(labelPairs) == 0 {
					labelPairs = []string{"-"}
				}
				row = append(row, strings.Join(labelPairs, ","))
			}
			if showExactTime {
				row = append(row, formattime.FormatTime(rip.CreatedAt.Local(), true))
			}
			body = append(body, row)
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Region", "Status", "IPv4", "IPv6", "AttachedTo"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			if showExactTime {
				headers = append(headers, "Created")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	ReservedIPsCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact creation time")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVar(&listRegion, "region", "", "Filter by region")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter reserved IPs (format: key1=value1,key2=value2)")

	_ = listCmd.RegisterFlagCompletionFunc("region", completeRegion)
}
