package dir

import (
	"github.com/spf13/cobra"
)

// DirCmd manages project-local directory configuration.
var DirCmd = &cobra.Command{
	Use:   "dir",
	Short: "Manage directory-local Thalassa defaults",
	Long: `Manage optional project-local defaults discovered from a .thalassa file (or .thalassa/config.yaml) by walking up from the current directory.

Directory config can set a preferred CLI context, organisation, project, and Kubernetes cluster.`,
}

func init() {
	DirCmd.AddCommand(showCmd)
	DirCmd.AddCommand(initCmd)
}
