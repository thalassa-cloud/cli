package projects

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var viewExactTime bool

var viewCmd = &cobra.Command{
	Use:               "view <project>",
	Short:             "View a project",
	Aliases:           []string{"show", "get", "describe"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		project, err := client.Projects().GetProject(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get project: %w", err)
		}

		parent := "-"
		if project.ParentProject != nil {
			parent = project.ParentProject.Name
			if project.ParentProject.Slug != "" {
				parent = fmt.Sprintf("%s (%s)", project.ParentProject.Name, project.ParentProject.Slug)
			}
		}

		fmt.Printf("ID:          %s\n", project.Identity)
		fmt.Printf("Name:        %s\n", project.Name)
		fmt.Printf("Slug:        %s\n", project.Slug)
		fmt.Printf("Description: %s\n", project.Description)
		fmt.Printf("Parent:      %s\n", parent)
		fmt.Printf("CreatedAt:   %s\n", formattime.FormatTime(project.CreatedAt.Local(), viewExactTime))
		return nil
	},
}

func init() {
	ProjectsCmd.AddCommand(viewCmd)
	viewCmd.Flags().BoolVar(&viewExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
