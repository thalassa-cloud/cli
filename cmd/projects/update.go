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
	updateName        string
	updateDescription string
	updateLabels      []string
	updateAnnotations []string
	updateParent      string
)

var updateCmd = &cobra.Command{
	Use:               "update <project>",
	Short:             "Update a project (only flags you set are changed)",
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

		req := clientprojects.UpdateProjectRequest{
			Name:        project.Name,
			Description: project.Description,
			Labels:      project.Labels,
			Annotations: project.Annotations,
		}
		if project.ParentProject != nil {
			parentIdentity := project.ParentProject.Identity
			req.ParentProjectIdentity = &parentIdentity
		}

		if cmd.Flags().Changed("name") {
			req.Name = updateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = shared.KeyValuePairsToMap(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = shared.KeyValuePairsToMap(updateAnnotations)
		}
		if cmd.Flags().Changed("parent") {
			parent := updateParent
			req.ParentProjectIdentity = &parent
		}

		out, err := client.Projects().UpdateProject(cmd.Context(), project.Identity, req)
		if err != nil {
			return fmt.Errorf("failed to update project: %w", err)
		}

		body := [][]string{{
			out.Identity,
			out.Name,
			out.Slug,
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
	ProjectsCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "do not print headers")
	updateCmd.Flags().StringVar(&updateName, "name", "", "Project display name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Project description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Replace labels (key=value, repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Replace annotations (key=value, repeatable)")
	updateCmd.Flags().StringVar(&updateParent, "parent", "", "Parent project identity or slug (empty clears parent)")
	_ = updateCmd.RegisterFlagCompletionFunc("parent", completion.CompleteProject)
}
