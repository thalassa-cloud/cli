package backup

import (
	"github.com/spf13/cobra"
)

const NoHeaderKey = "no-header"

// BackupCmd represents the backup command
var BackupCmd = &cobra.Command{
	Use:     "backup",
	Aliases: []string{"backups"},
	Short:   "Manage database backups",
	Long:    "Manage database backups for database clusters",
}
