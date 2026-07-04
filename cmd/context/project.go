package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/fzf"
	"github.com/thalassa-cloud/cli/internal/projectresolve"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var projectCmd = &cobra.Command{
	Use:     "project <project>",
	Aliases: []string{"set-project", "use-project"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Set the project in the current context",
	Long:    "Set the project in the current context. Accepts a project identity or slug; the identity is stored in the context. Use \"root\" to clear the project and scope commands to the organisation root.",
	Example: "tcloud context project <project>",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentContext, err := contextstate.GetContextConfiguration()
		if err != nil {
			return err
		}
		selectedProject, err := getSelectedProject(args)
		if err != nil {
			return err
		}
		selectedProject = projectresolve.NormalizeProjectRef(selectedProject)
		if selectedProject != "" {
			client, err := thalassaclient.GetThalassaClient()
			if err != nil {
				return err
			}
			selectedProject, err = projectresolve.ResolveProjectIdentity(cmd.Context(), client.Projects(), selectedProject)
			if err != nil {
				return err
			}
		}
		currentContext.Project = selectedProject
		err = contextstate.CombineConfigContext(currentContext)
		if err != nil {
			return err
		}
		fmt.Println(currentContext.Project)
		return contextstate.Save()
	},
}

func getSelectedProject(args []string) (string, error) {
	if contextstate.ProjectFlag != "" {
		return contextstate.ProjectFlag, nil
	}

	if len(args) == 0 && fzf.IsInteractiveMode(os.Stdout) {
		command := fmt.Sprintf("%s projects list --no-header --include-root", os.Args[0])
		return fzf.InteractiveChoice(command)
	} else if len(args) == 1 {
		return args[0], nil
	}
	return "", nil
}

func init() {
	ContextCmd.AddCommand(projectCmd)
}
