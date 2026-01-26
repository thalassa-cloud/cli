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
	"github.com/thalassa-cloud/client-go/dbaas"
)

var (
	backupScheduleListClusterFilter string
	backupScheduleListNoHeader      bool
	backupScheduleListShowExactTime bool
	backupScheduleListShowLabels    bool
)

// backupScheduleListCmd represents the backup-schedules list command
var backupScheduleListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List database backup schedules",
	Long:    "List database backup schedules for a specific cluster or all schedules in the organisation",
	Aliases: []string{"ls", "get", "clusters", "cluster"},
	Args:    cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return completion.CompleteDbClusterID(cmd, args, toComplete)
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		var schedules []dbaas.DbClusterBackupSchedule

		// If cluster identity is provided as argument, list schedules for that cluster
		if len(args) > 0 {
			clusterIdentity := args[0]
			schedules, err = client.DBaaS().ListDbBackupSchedules(cmd.Context(), clusterIdentity)
			if err != nil {
				return fmt.Errorf("failed to list backup schedules for cluster: %w", err)
			}
		} else {
			// Otherwise list all schedules for the organisation
			schedules, err = client.DBaaS().ListDbBackupSchedulesForOrganisation(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list backup schedules: %w", err)
			}
		}

		if len(schedules) == 0 {
			fmt.Println("No backup schedules found")
			return nil
		}

		body := make([][]string, 0, len(schedules))
		for _, schedule := range schedules {
			clusterName := "-"
			if schedule.DbCluster != nil {
				clusterName = schedule.DbCluster.Name
			}

			status := string(schedule.Status)
			if schedule.Suspended {
				status = fmt.Sprintf("%s (suspended)", status)
			}
			if schedule.DeleteScheduledAt != nil {
				status = fmt.Sprintf("%s (deletion scheduled)", status)
			}

			nextBackup := "-"
			if schedule.NextBackupAt != nil {
				nextBackup = formattime.FormatTime(schedule.NextBackupAt.Local(), backupScheduleListShowExactTime)
			}

			lastBackup := "-"
			if schedule.LastBackupAt != nil {
				lastBackup = formattime.FormatTime(schedule.LastBackupAt.Local(), backupScheduleListShowExactTime)
			}

			row := []string{
				schedule.Identity,
				schedule.Name,
				clusterName,
				string(schedule.Method),
				schedule.Schedule,
				schedule.RetentionPolicy,
				fmt.Sprintf("%d", schedule.BackupCount),
				nextBackup,
				lastBackup,
				status,
				formattime.FormatTime(schedule.CreatedAt.Local(), backupScheduleListShowExactTime),
			}

			if backupScheduleListShowLabels {
				labelStrs := []string{}
				for k, v := range schedule.Labels {
					labelStrs = append(labelStrs, k+"="+v)
				}
				sort.Strings(labelStrs)
				if len(labelStrs) == 0 {
					labelStrs = []string{"-"}
				}
				row = append(row, strings.Join(labelStrs, ","))
			}

			body = append(body, row)
		}

		if backupScheduleListNoHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Cluster", "Method", "Schedule", "Retention", "Backups", "Next Backup", "Last Backup", "Status", "Created"}
			if backupScheduleListShowLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}

		return nil
	},
}

func init() {
	BackupSchedulesCmd.AddCommand(backupScheduleListCmd)

	backupScheduleListCmd.Flags().BoolVar(&backupScheduleListNoHeader, NoHeaderKey, false, "Do not print the header")
	backupScheduleListCmd.Flags().BoolVar(&backupScheduleListShowExactTime, "exact-time", false, "Show exact time instead of relative time")
	backupScheduleListCmd.Flags().BoolVar(&backupScheduleListShowLabels, "show-labels", false, "Show labels")
	backupScheduleListCmd.Flags().StringVar(&backupScheduleListClusterFilter, "cluster", "", "Filter by database cluster identity, slug, or name")

	// Register completions
	backupScheduleListCmd.RegisterFlagCompletionFunc("cluster", completion.CompleteDbClusterID)
}
