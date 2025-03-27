package version

import (
	"github.com/spf13/cobra"

	cliversion "github.com/thalassa-cloud/cli/internal/version"
)

var VersionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print version information",
	Long:    `Print version information about the Thalassa Cloud CLI. This command will display the version of the CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		cliversion.PrintVersion()
	},
}
