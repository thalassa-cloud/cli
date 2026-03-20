package teams

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/teams/members"
)

// TeamsCmd represents the teams command
var TeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Manage organisation teams",
}

func init() {
	TeamsCmd.AddCommand(members.TeamMembersCmd)
}
