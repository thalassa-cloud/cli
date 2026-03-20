package serviceaccounts

import "github.com/spf13/cobra"

// ServiceAccountsCmd manages organisation service accounts.
var ServiceAccountsCmd = &cobra.Command{
	Use:     "service-accounts",
	Aliases: []string{"sa", "service-account"},
	Short:   "Organisation service accounts",
}
