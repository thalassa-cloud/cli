package invites

import "github.com/spf13/cobra"

// InvitesCmd lists organisation member invitations.
var InvitesCmd = &cobra.Command{
	Use:   "invites",
	Short: "Organisation member invitations",
}
