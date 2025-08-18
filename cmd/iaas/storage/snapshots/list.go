package snapshots

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime bool
	showLabels    bool
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
		snapshots, err := client.IaaS().ListSnapshots(cmd.Context(), &iaas.ListSnapshotsRequest{})
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
			headers := []string{"ID", "Name", "Region", "Size", "Age"}
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
}
