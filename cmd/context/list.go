package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/table"
)

const (
	NoHeaderFlag    = "no-header"
	ShowCurrentFlag = "show-current"
)

var noHeader bool
var showCurrent bool

// listContextCmd represents the get command
var listContextCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"get", "list", "ls", "l", "g"},
	Short:   "List the contexts",
	Long:    "List the contexts from the config",
	Example: "tcloud context list",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentContext, err := contextstate.GetContextConfiguration()
		currentContextName := ""
		if err == nil {
			currentContextName = currentContext.Name
		}
		contexts := contextstate.GlobalConfigManager().Config().Contexts
		body := make([][]string, 0, len(contexts))
		for _, c := range contexts {
			name := c.Name
			if showCurrent && name == currentContextName {
				name = fmt.Sprintf("* %s", name)
			}
			body = append(body, []string{name, c.Context.Organisation, c.Context.User, c.Context.API})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Name", "Organisation", "User", "Endpoint"}, body)
		}
		return nil
	},
}

func init() {
	ContextCmd.AddCommand(listContextCmd)
	listContextCmd.Flags().BoolVar(&noHeader, NoHeaderFlag, false, "Do not print the header")
	listContextCmd.Flags().BoolVar(&showCurrent, ShowCurrentFlag, true, "Show the current context")
}
