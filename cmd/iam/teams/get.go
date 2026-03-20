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

var getCmd = &cobra.Command{
	Use:               "get <team>",
	Short:             "Show a team and its members",
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
		fmt.Printf("Identity:    %s\n", team.Identity)
		fmt.Printf("Name:        %s\n", team.Name)
		fmt.Printf("Slug:        %s\n", team.Slug)
		fmt.Printf("Description: %s\n", team.Description)
		fmt.Printf("Created:     %s\n", formattime.FormatTime(team.CreatedAt.Local(), showExactTime))
		if len(team.Members) > 0 {
			fmt.Println("\nMembers:")
			body := make([][]string, 0, len(team.Members))
			for _, m := range team.Members {
				body = append(body, []string{
					m.Identity,
					m.Role,
					shared.UserDisplay(m.User),
				})
			}
			table.Print([]string{"ID", "Role", "User"}, body)
		}
		return nil
	},
}

func init() {
	TeamsCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	getCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
