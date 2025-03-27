package context

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a context",
	Long:    "Delete a context from the config",
	Example: "tcloud context delete <context>",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return contextstate.Clear()
	},
}

func init() {
	ContextCmd.AddCommand(deleteCmd)
}
