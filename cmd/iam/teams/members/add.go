package members

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var addUser string
var addRole string

var addCmd = &cobra.Command{
	Use:   "add <team>",
	Short: "Add a user to a team",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if addUser == "" || addRole == "" {
			return fmt.Errorf("--user and --role are required")
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		if err := client.IAM().AddTeamMember(ctx, args[0], clientiam.AddTeamMemberRequest{
			UserIdentity: addUser,
			Role:         addRole,
		}); err != nil {
			return fmt.Errorf("failed to add team member: %w", err)
		}
		fmt.Printf("Added user %s to team %s with role %s\n", addUser, args[0], addRole)
		return nil
	},
}

func init() {
	TeamMembersCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	addCmd.Flags().StringVar(&addUser, "user", "", "User identity to add")
	addCmd.Flags().StringVar(&addRole, "role", "", "Team role for the user")
	_ = addCmd.RegisterFlagCompletionFunc("user", completion.CompleteIAMAppUserSubject)
	_ = addCmd.MarkFlagRequired("user")
	_ = addCmd.MarkFlagRequired("role")
}
