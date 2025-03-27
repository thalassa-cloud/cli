package context

import (
	"github.com/spf13/cobra"
)

// ContextCmd represents the get command
var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage context",
	Long:  "Manage context for the CLI. Contexts are used to manage multiple organisations and APIs.",
}

func init() {
}
