package projects

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientprojects "github.com/thalassa-cloud/client-go/projects"
)

var (
	createName        string
	createDescription string
	createLabels      []string
	createAnnotations []string
	createParent      string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project in the current organisation",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if createName == "" {
			return fmt.Errorf("--name is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req := clientprojects.CreateProjectRequest{
			Name:        createName,
			Description: createDescription,
			Labels:      shared.KeyValuePairsToMap(createLabels),
			Annotations: shared.KeyValuePairsToMap(createAnnotations),
		}
		if cmd.Flags().Changed("parent") {
			parent := createParent
			req.ParentProjectIdentity = &parent
		}

		project, err := client.Projects().CreateProject(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		body := [][]string{{
			project.Identity,
			project.Name,
			project.Slug,
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug"}, body)
		}
		return nil
	},
}

func init() {
	ProjectsCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "do not print headers")
	createCmd.Flags().StringVar(&createName, "name", "", "Project display name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Project description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	createCmd.Flags().StringVar(&createParent, "parent", "", "Parent project identity or slug")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.RegisterFlagCompletionFunc("parent", completion.CompleteProject)
}
