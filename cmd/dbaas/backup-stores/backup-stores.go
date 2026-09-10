package backupstores

import (
	"github.com/spf13/cobra"
)

const NoHeaderKey = "no-header"

// BackupStoresCmd manages DBaaS backup stores.
var BackupStoresCmd = &cobra.Command{
	Use:     "backup-stores",
	Aliases: []string{"backup-store"},
	Short:   "Manage database backup stores",
	Long:    "Manage backup stores used by database clusters for Barman backups and point-in-time recovery",
}
