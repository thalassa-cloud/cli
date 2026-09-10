package secrets

import (
	"github.com/spf13/cobra"
)

// SecretsCmd manages Secrets Manager resources.
var SecretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"secret"},
	Short:   "Manage secrets",
}

var (
	noHeader      bool
	showExactTime bool
)
