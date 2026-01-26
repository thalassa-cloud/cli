package backupschedules

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupScheduleDeleteForce bool
)

// backupScheduleDeleteCmd represents the backup-schedules delete command
var backupScheduleDeleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Delete a database backup schedule",
	Long:              "Delete a database backup schedule",
	Aliases:           []string{"d", "rm", "del", "remove"},
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completion.CompleteDbClusterID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusterIdentity := args[0]
		scheduleIdentity := args[1]

		// Get schedule to show name in confirmation
		var scheduleName string
		if !backupScheduleDeleteForce {
			schedule, err := client.DBaaS().GetDbBackupSchedule(cmd.Context(), clusterIdentity, scheduleIdentity)
			if err != nil {
				if tcclient.IsNotFound(err) {
					return fmt.Errorf("backup schedule not found: %s", scheduleIdentity)
				}
				return fmt.Errorf("failed to get backup schedule: %w", err)
			}
			scheduleName = schedule.Name
		}

		// Ask for confirmation unless --force is provided
		if !backupScheduleDeleteForce {
			fmt.Printf("Are you sure you want to delete backup schedule %s (%s)?\n", scheduleName, scheduleIdentity)
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		err = client.DBaaS().DeleteDbBackupSchedule(cmd.Context(), clusterIdentity, scheduleIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("backup schedule not found: %s", scheduleIdentity)
			}
			return fmt.Errorf("failed to delete backup schedule: %w", err)
		}

		fmt.Printf("Backup schedule %s deleted successfully\n", scheduleIdentity)
		return nil
	},
}

func init() {
	BackupSchedulesCmd.AddCommand(backupScheduleDeleteCmd)

	backupScheduleDeleteCmd.Flags().BoolVar(&backupScheduleDeleteForce, "force", false, "Force the deletion and skip the confirmation")
}
