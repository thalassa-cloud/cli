package snapshotpolicies

import (
	"fmt"
	"sort"
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

var (
	showExactTime     bool
	showLabels        bool
	listLabelSelector string
	listRegion        string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List snapshot policies",
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

		policies, err := client.IaaS().ListSnapshotPolicies(cmd.Context(), &iaas.ListSnapshotPoliciesRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(policies))
		for _, policy := range policies {
			row := []string{
				policy.Identity,
				policy.Name,
				regionName(policy),
				fmt.Sprintf("%t", policy.Enabled),
				policy.Schedule,
				policy.Ttl.String(),
				keepCountString(policy.KeepCount),
				targetSummary(policy.Target),
				formattime.FormatTime(policy.CreatedAt.Local(), showExactTime),
			}
			if showLabels {
				labelPairs := make([]string, 0, len(policy.Labels))
				for k, v := range policy.Labels {
					labelPairs = append(labelPairs, k+"="+v)
				}
				sort.Strings(labelPairs)
				if len(labelPairs) == 0 {
					labelPairs = []string{"-"}
				}
				row = append(row, strings.Join(labelPairs, ","))
			}
			body = append(body, row)
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Region", "Enabled", "Schedule", "TTL", "Keep", "Target", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	SnapshotPoliciesCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVar(&listRegion, "region", "", "Filter by region")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter snapshot policies")

	_ = listCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
