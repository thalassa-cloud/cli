package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
)

// currentContextCmd represents the get command
var currentContextCmd = &cobra.Command{
	Use:     "current",
	Short:   "Shows the current context",
	Long:    "Shows the current context (or the context set with the --context flag)",
	Example: "tcloud context current",
	Args:    cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		currentContext, err := contextstate.GetContextConfiguration()
		if err != nil {
			return err
		}
		fmt.Println(currentContext.Name)
		return nil
	},
}

func init() {
	ContextCmd.AddCommand(currentContextCmd)
}
