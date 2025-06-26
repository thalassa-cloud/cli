package context

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/fzf"
)

// organisationCmd represents the organisation command
var organisationCmd = &cobra.Command{
	Use:     "organisation <organisation>",
	Aliases: []string{"set-organisation", "use-organisation", "organisation", "org"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Set the organisation in the current-context",
	Long:    "Set the organisation in the current-context",
	Example: "tcloud context use-organisation <organisation>",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentContext, err := contextstate.GetContextConfiguration()
		if err != nil {
			return err
		}
		selectedOrganisation, err := getSelectedOrganisation(args)
		if err != nil {
			return err
		}
		currentContext.Organisation = selectedOrganisation
		err = contextstate.CombineConfigContext(currentContext)
		if err != nil {
			return err
		}
		fmt.Println(currentContext.Organisation)
		return contextstate.Save()
	},
}

func getSelectedOrganisation(args []string) (string, error) {
	if contextstate.OrganisationFlag != "" {
		return contextstate.OrganisationFlag, nil
	}

	if len(args) == 0 && fzf.IsInteractiveMode(os.Stdout) {
		command := fmt.Sprintf("%s me organisations --no-header --slug-only", os.Args[0])
		return fzf.InteractiveChoice(command)
	} else if len(args) == 1 {
		return args[0], nil
	} else {
		return "", ErrInvalidOrganisation
	}
}

func init() {
	ContextCmd.AddCommand(organisationCmd)
}

var (
	ErrInvalidOrganisation = errors.New("invalid organisation")
)
