package members

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var removeForce bool

var removeCmd = &cobra.Command{
	Use:               "remove <team> <member>",
	Aliases:           []string{"rm"},
	Short:             "Remove a member from a team",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completion.CompleteIAMTeamThenTeamMember,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		ok, err := shared.PromptDestructiveUnlessForce(removeForce, fmt.Sprintf("Are you sure you want to remove this member from the team?\n  Team: %s\n  Member: %s\n", args[0], args[1]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := client.IAM().RemoveTeamMember(ctx, args[0], args[1]); err != nil {
			return fmt.Errorf("failed to remove team member: %w", err)
		}
		fmt.Printf("Removed member %s from team %s\n", args[1], args[0])
		return nil
	},
}

func init() {
	TeamMembersCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVar(&removeForce, shared.ForceKey, false, "Skip the confirmation prompt and remove the member")
	removeCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
