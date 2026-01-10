package objectstorage

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/objectstorage"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime bool
	showLabels    bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List object storage buckets",
	Aliases: []string{"l", "ls", "buckets", "bucket"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		var buckets []objectstorage.ObjectStorageBucket
		buckets, err = client.ObjectStorage().ListBuckets(cmd.Context())
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(buckets))
		for _, bucket := range buckets {
			regionName := "-"
			if bucket.Region != nil {
				regionName = bucket.Region.Name
				if regionName == "" {
					regionName = bucket.Region.Slug
				}
				if regionName == "" {
					regionName = bucket.Region.Identity
				}
			}

			row := []string{
				bucket.Name,
				bucket.Status,
				regionName,
				fmt.Sprintf("%.2f GB", bucket.Usage.TotalSizeGB),
				fmt.Sprintf("%d", bucket.Usage.TotalObjects),
				string(bucket.Versioning),
				fmt.Sprintf("%v", bucket.Public),
				formattime.FormatTime(bucket.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labelStrs := []string{}
				for k, v := range bucket.Labels {
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
			headers := []string{"Name", "Status", "Region", "Size", "Objects", "Versioning", "Public", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	ObjectStorageCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
}
