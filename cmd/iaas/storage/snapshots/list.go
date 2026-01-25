package snapshots

import (
	"fmt"
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
	listRegionFilter  string
	listStatusFilter  string
	listVolumeFilter  string
)

// getCmd represents the get command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of snapshots",
	Aliases: []string{"g", "get", "ls"},
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
		if listRegionFilter != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "region",
				Value: listRegionFilter,
			})
		}
		if listStatusFilter != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "status",
				Value: listStatusFilter,
			})
		}
		if listVolumeFilter != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "volume",
				Value: listVolumeFilter,
			})
		}

		snapshots, err := client.IaaS().ListSnapshots(cmd.Context(), &iaas.ListSnapshotsRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(snapshots))
		for _, snapshot := range snapshots {
			size := 0
			if snapshot.SizeGB != nil {
				size = *snapshot.SizeGB
			}

			item := []string{
				snapshot.Identity,
				snapshot.Name,
				string(snapshot.Status),
				snapshot.Region.Name,
				fmt.Sprintf("%dGB", size),
				formattime.FormatTime(snapshot.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labels := []string{}
				for k, v := range snapshot.Annotations {
					labels = append(labels, fmt.Sprintf("%s=%s", k, v))
				}
				item = append(item, strings.Join(labels, ", "))
			}
			body = append(body, item)
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Status", "Region", "Size", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	SnapshotsCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter snapshots (format: key1=value1,key2=value2)")
	listCmd.Flags().StringVar(&listRegionFilter, "region", "", "Region of the snapshot")
	listCmd.Flags().StringVar(&listStatusFilter, "status", "", "Status of the snapshot")
	listCmd.Flags().StringVar(&listVolumeFilter, "volume", "", "Source volume of the snapshot")

	// Register completions
	listCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	listCmd.RegisterFlagCompletionFunc("volume", completion.CompleteVolumeID)
}
