package tfs

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
	"github.com/thalassa-cloud/client-go/tfs"
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
	Short:   "Get a list of TFS instances",
	Aliases: []string{"l", "ls", "get"},
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

		instances, err := client.Tfs().ListTfsInstances(cmd.Context(), &tfs.ListTfsInstancesRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(instances))
		for _, instance := range instances {
			regionName := ""
			if instance.Region != nil {
				regionName = instance.Region.Name
			}

			row := []string{
				instance.Identity,
				instance.Name,
				string(instance.Status),
				regionName,
				formattime.FormatTime(instance.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labelStrs := []string{}
				for k, v := range instance.Labels {
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
			headers := []string{"ID", "Name", "Status", "Region", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	TfsCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter TFS instances (format: key1=value1,key2=value2)")
}
