package projects

import (
	"github.com/spf13/cobra"
)

var ProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects (beta)",
	Long: `Manage projects within the organisation selected in your context.

Note: This command is in beta and requires the project feature gate to be enabled on your organisation.`,
}
