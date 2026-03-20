package teams

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	updateName        string
	updateDescription string
	updateLabels      []string
	updateAnnotations []string
)

var updateCmd = &cobra.Command{
	Use:               "update <team>",
	Short:             "Update a team (only flags you set are changed)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMTeamIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		team, err := client.IAM().GetTeam(ctx, args[0], &clientiam.GetTeamRequest{})
		if err != nil {
			return fmt.Errorf("failed to get team: %w", err)
		}
		update := clientiam.UpdateTeam{
			Name:        team.Name,
			Description: team.Description,
			Labels:      team.Labels,
			Annotations: team.Annotations,
		}
		if cmd.Flags().Changed("name") {
			update.Name = updateName
		}
		if cmd.Flags().Changed("description") {
			update.Description = updateDescription
		}
		if cmd.Flags().Changed("labels") {
			update.Labels = shared.KeyValuePairsToMap(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			update.Annotations = shared.KeyValuePairsToMap(updateAnnotations)
		}
		out, err := client.IAM().UpdateTeam(ctx, team.Identity, update)
		if err != nil {
			return fmt.Errorf("failed to update team: %w", err)
		}
		body := [][]string{{out.Identity, out.Name, out.Slug, formattime.FormatTime(out.CreatedAt.Local(), showExactTime)}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "Created"}, body)
		}
		return nil
	},
}

func init() {
	TeamsCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	updateCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	updateCmd.Flags().StringVar(&updateName, "name", "", "Team display name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Team description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Replace labels (key=value, repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Replace annotations (key=value, repeatable)")
}
