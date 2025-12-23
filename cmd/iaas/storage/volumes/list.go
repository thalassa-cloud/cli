package volumes

import (
	"fmt"
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
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of volumes",
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

		volumes, err := client.IaaS().ListVolumes(cmd.Context(), &iaas.ListVolumesRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(volumes))
		for _, volume := range volumes {
			volumeType := ""
			if volume.VolumeType != nil {
				volumeType = volume.VolumeType.Name
			}

			item := []string{
				volume.Identity,
				volume.Name,
				volume.Status,
				volume.Region.Name,
				volumeType,
				fmt.Sprintf("%dGB", volume.Size),
				formattime.FormatTime(volume.CreatedAt.Local(), showExactTime),
			}
			if showLabels {
				labels := []string{}
				for k, v := range volume.Annotations {
					labels = append(labels, fmt.Sprintf("%s=%s", k, v))
				}
				item = append(item, strings.Join(labels, ", "))
			}
			body = append(body, item)
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Status", "Region", "Type", "Size", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	VolumesCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	getCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	getCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter volumes (format: key1=value1,key2=value2)")
}
