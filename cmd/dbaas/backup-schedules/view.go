package backupschedules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupScheduleViewShowExactTime bool
	backupScheduleViewNoHeader      bool
)

// backupScheduleViewCmd represents the backup-schedules view command
var backupScheduleViewCmd = &cobra.Command{
	Use:               "view",
	Short:             "View backup schedule details",
	Long:              "View detailed information about a database backup schedule",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completion.CompleteDbClusterID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusterIdentity := args[0]
		scheduleIdentity := args[1]

		schedule, err := client.DBaaS().GetDbBackupSchedule(cmd.Context(), clusterIdentity, scheduleIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("backup schedule not found: %s", scheduleIdentity)
			}
			return fmt.Errorf("failed to get backup schedule: %w", err)
		}

		clusterName := ""
		if schedule.DbCluster != nil {
			clusterName = schedule.DbCluster.Name
		}

		body := [][]string{
			{"ID", schedule.Identity},
			{"Name", schedule.Name},
			{"Cluster", clusterName},
			{"Status", string(schedule.Status)},
			{"Method", string(schedule.Method)},
			{"Schedule", schedule.Schedule},
			{"Retention Policy", schedule.RetentionPolicy},
			{"Backup Count", fmt.Sprintf("%d", schedule.BackupCount)},
			{"Suspended", fmt.Sprintf("%v", schedule.Suspended)},
			{"Created", formattime.FormatTime(schedule.CreatedAt.Local(), backupScheduleViewShowExactTime)},
		}

		if schedule.Description != nil && *schedule.Description != "" {
			body = append(body, []string{"Description", *schedule.Description})
		}

		if schedule.NextBackupAt != nil {
			body = append(body, []string{"Next Backup", formattime.FormatTime(schedule.NextBackupAt.Local(), backupScheduleViewShowExactTime)})
		}

		if schedule.LastBackupAt != nil {
			body = append(body, []string{"Last Backup", formattime.FormatTime(schedule.LastBackupAt.Local(), backupScheduleViewShowExactTime)})
		}

		if schedule.StatusMessage != "" {
			body = append(body, []string{"Status Message", schedule.StatusMessage})
		}

		if schedule.DeleteScheduledAt != nil {
			body = append(body, []string{"Delete Scheduled At", formattime.FormatTime(schedule.DeleteScheduledAt.Local(), backupScheduleViewShowExactTime)})
		}

		if len(schedule.Labels) > 0 {
			labelStrs := []string{}
			for k, v := range schedule.Labels {
				labelStrs = append(labelStrs, k+"="+v)
			}
			sort.Strings(labelStrs)
			body = append(body, []string{"Labels", strings.Join(labelStrs, ", ")})
		}

		if len(schedule.Annotations) > 0 {
			annotationStrs := []string{}
			for k, v := range schedule.Annotations {
				annotationStrs = append(annotationStrs, k+"="+v)
			}
			sort.Strings(annotationStrs)
			body = append(body, []string{"Annotations", strings.Join(annotationStrs, ", ")})
		}

		if backupScheduleViewNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Field", "Value"}, body)
		}

		return nil
	},
}

func init() {
	BackupSchedulesCmd.AddCommand(backupScheduleViewCmd)

	backupScheduleViewCmd.Flags().BoolVar(&backupScheduleViewNoHeader, NoHeaderKey, false, "Do not print the header")
	backupScheduleViewCmd.Flags().BoolVar(&backupScheduleViewShowExactTime, "exact-time", false, "Show exact time instead of relative time")
}
