package namespaces

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/containerregistry"
	"github.com/thalassa-cloud/client-go/filters"
)

var (
	noHeader          bool
	showExactTime     bool
	showLabels        bool
	listLabelSelector string
	listRegion        string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List container registry namespaces",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := filters.Filters{}
		if listRegion != "" {
			f = append(f, &filters.FilterKeyValue{Key: "region", Value: listRegion})
		}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		namespaces, err := client.ContainerRegistry().ListContainerRegistryNamespaces(cmd.Context(), &containerregistry.ListContainerRegistryNamespacesRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(namespaces))
		for _, ns := range namespaces {
			regionName := "-"
			if ns.Region != nil {
				regionName = ns.Region.Name
				if regionName == "" {
					regionName = ns.Region.Identity
				}
			}

			row := []string{
				ns.Identity,
				ns.Namespace,
				regionName,
				fmt.Sprintf("%d", len(ns.Repositories)),
				formatBytes(ns.TotalSizeBytes),
				formattime.FormatTime(ns.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labelPairs := labelPairs(ns.Labels)
				row = append(row, joinLabelPairs(labelPairs))
			}

			body = append(body, row)
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Namespace", "Region", "Repositories", "Size", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	NamespacesCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().StringVar(&listRegion, "region", "", "Filter by region")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector (format: key1=value1,key2=value2)")

	listCmd.RegisterFlagCompletionFunc("region", completeRegion)
}

func labelPairs(labels map[string]string) []string {
	pairs := make([]string, 0, len(labels))
	for k, v := range labels {
		pairs = append(pairs, k+"="+v)
	}
	sort.Strings(pairs)
	return pairs
}

func joinLabelPairs(pairs []string) string {
	if len(pairs) == 0 {
		return "-"
	}
	return strings.Join(pairs, ",")
}

func formatBytes(b int64) string {
	if b == 0 {
		return "-"
	}
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
