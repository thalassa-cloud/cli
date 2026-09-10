package secrets

import (
	"github.com/spf13/cobra"
)

// SecretsCmd manages Secrets Manager resources.
var SecretsCmd = &cobra.Command{
	Use:          "secrets",
	Aliases:      []string{"secret"},
	Short:        "Manage secrets (beta)",
	SilenceUsage: true,
	Long: `Manage Secrets Manager paths, versions, and access policies.

Note: This command is in beta. Commands that list or view secrets show
metadata only; use get-value when you intentionally need secret material.`,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		// Avoid dumping full flag help on runtime / API validation errors.
		cmd.SilenceUsage = true
	},
}

var (
	noHeader      bool
	showExactTime bool
)
