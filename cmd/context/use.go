package context

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/fzf"
)

// useContextCmd represents the use command
var useContextCmd = &cobra.Command{
	Use:     "use <context>",
	Aliases: []string{"set", "s", "select"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Set the current context",
	Long:    "Set the current context (or the context set with the --context flag)",
	Example: "tcloud context use <context>",
	RunE: func(cmd *cobra.Command, args []string) error {
		context, err := getSelectedContext(args)
		if err != nil {
			return err
		}

		fmt.Println(context)
		err = contextstate.Set(context)
		if err != nil {
			return err
		}
		return contextstate.Save()
	},
}

func getSelectedContext(args []string) (string, error) {
	if len(args) == 0 && fzf.IsInteractiveMode(os.Stdout) {
		return fzf.InteractiveChoice(fmt.Sprintf("%s context list --%s", os.Args[0], NoHeaderFlag))
	} else if len(args) == 1 {
		return args[0], nil
	} else {
		return "", errors.New("invalid context")
	}
}

func init() {
	ContextCmd.AddCommand(useContextCmd)
}
