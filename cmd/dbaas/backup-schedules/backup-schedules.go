package backupschedules

import (
	"github.com/spf13/cobra"
)

const NoHeaderKey = "no-header"

// BackupSchedulesCmd represents the backup-schedules command
var BackupSchedulesCmd = &cobra.Command{
	Use:     "backup-schedules",
	Aliases: []string{"backup-schedule", "schedules", "schedule"},
	Short:   "Manage database backup schedules",
	Long:    "Manage database backup schedules for database clusters",
}
