package members

import "github.com/spf13/cobra"

// TeamMembersCmd manages membership on a team (distinct from organisation members).
var TeamMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage team membership",
}
