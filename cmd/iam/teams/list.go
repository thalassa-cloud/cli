package teams

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	noHeader               bool
	showExactTime          bool
	teamsListLabelSelector string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List teams",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		req := &clientiam.ListTeamsRequest{}
		if teamsListLabelSelector != "" {
			req.Filters = []filters.Filter{
				&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(teamsListLabelSelector)},
			}
		}
		teams, err := client.IAM().ListTeams(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list teams: %w", err)
		}
		body := make([][]string, 0, len(teams))
		for _, t := range teams {
			body = append(body, []string{
				t.Identity,
				t.Name,
				t.Slug,
				t.Description,
				fmt.Sprintf("%d", len(t.Members)),
				formattime.FormatTime(t.CreatedAt.Local(), showExactTime),
			})
		}
		headers := []string{"ID", "Name", "Slug", "Description", "Members", "Age"}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	TeamsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	listCmd.Flags().StringVar(&teamsListLabelSelector, "label-selector", "", "Filter by labels (key=value,key2=value2)")
}
