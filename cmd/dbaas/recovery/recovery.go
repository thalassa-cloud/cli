package recovery

import (
	"github.com/spf13/cobra"
)

const NoHeaderKey = "no-header"

// RecoveryCmd inspects backup recoverability and WAL archives.
var RecoveryCmd = &cobra.Command{
	Use:     "recovery",
	Aliases: []string{"pitr"},
	Short:   "Inspect database backup recovery and PITR windows",
	Long:    "Inspect backup store recovery overviews, WAL coverage, and point-in-time recovery windows",
}
