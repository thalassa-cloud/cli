package dbaas

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/dbaas/backup"
	backupschedules "github.com/thalassa-cloud/cli/cmd/dbaas/backup-schedules"
)

var DbaasCmd = &cobra.Command{
	Use:     "dbaas",
	Aliases: []string{"db", "database"},
	Short:   "Manage database clusters and related services",
	Long:    "DBaaS commands to manage your database clusters and related services within the Thalassa Cloud Platform",
	Example: "tcloud dbaas list\ntcloud dbaas instance-types\ntcloud dbaas versions --engine postgres",
}

func init() {
	DbaasCmd.AddCommand(backup.BackupCmd)
	DbaasCmd.AddCommand(backupschedules.BackupSchedulesCmd)
}
