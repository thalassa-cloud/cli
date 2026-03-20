package members

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var noHeader bool

var listCmd = &cobra.Command{
	Use:               "list <team>",
	Short:             "List members of a team",
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
		body := make([][]string, 0, len(team.Members))
		for _, m := range team.Members {
			body = append(body, []string{m.Identity, m.Role, shared.UserDisplay(m.User)})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Role", "User"}, body)
		}
		return nil
	},
}

func init() {
	TeamMembersCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
