package oidc

import (
	"github.com/spf13/cobra"
)

// OidcCmd represents the oidc command
var OidcCmd = &cobra.Command{
	Use:   "oidc",
	Short: "OIDC token operations",
	Long:  "OIDC token operations for Thalassa Cloud, such as OIDC Federation and token exchange",
}

func init() {
	OidcCmd.AddCommand(tokenExchangeCmd)
}
