package projects

import (
	"github.com/spf13/cobra"
)

var ProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects (private beta)",
	Long: `Manage projects within the organisation selected in your context.

Note: This command is in private beta and requires the project feature gate to be enabled on your organisation.`,
}
