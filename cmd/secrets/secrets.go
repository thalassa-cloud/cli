package secrets

import (
	"github.com/spf13/cobra"
)

// SecretsCmd manages Secrets Manager resources.
var SecretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"secret"},
	Short:   "Manage secrets (beta)",
	Long: `Manage Secrets Manager paths, versions, and access policies.

Note: This command is in beta. Commands that list or view secrets show
metadata only; use get-value when you intentionally need secret material.`,
}

var (
	noHeader      bool
	showExactTime bool
)
