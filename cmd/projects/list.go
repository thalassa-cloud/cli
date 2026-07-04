package projects

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/projectresolve"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

const (
	NoHeaderKey = "no-header"
)

var (
	noHeader    bool
	slugOnly    bool
	includeRoot bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects in the current organisation",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		projects, err := client.Projects().ListProjects(cmd.Context(), nil)
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}

		body := make([][]string, 0, len(projects)+1)
		if includeRoot {
			if slugOnly {
				body = append(body, []string{projectresolve.RootRef})
			} else {
				body = append(body, []string{
					projectresolve.RootRef,
					"Organisation root",
					"-",
				})
			}
		}
		for _, project := range projects {
			if slugOnly {
				body = append(body, []string{project.Slug})
			} else {
				body = append(body, []string{
					project.Identity,
					project.Name,
					project.Slug,
				})
			}
		}

		headers := []string{"ID", "Name", "Slug"}
		if slugOnly {
			headers = []string{"Slug"}
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print(headers, body)
		}

		return nil
	},
}

func init() {
	ProjectsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "do not print headers")
	listCmd.Flags().BoolVar(&slugOnly, "slug-only", false, "only print the slug")
	listCmd.Flags().BoolVar(&includeRoot, "include-root", false, "include organisation root (no project) as the first entry")
}
