package context

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
)

// viewCmd represents the get command
var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "Shows current context",
	Long:    "Shows the current context",
	Example: "tcloud context view",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := yaml.Marshal(contextstate.GlobalConfigManager().Config())
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", data)
		return nil
	},
}

func init() {
	ContextCmd.AddCommand(viewCmd)
}
