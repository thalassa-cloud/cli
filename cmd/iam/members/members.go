package members

import "github.com/spf13/cobra"

// MembersCmd represents organisation members (owners and members).
var MembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Organisation members (owners and members)",
}
