package backupschedules

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupScheduleUpdateName            string
	backupScheduleUpdateDescription     string
	backupScheduleUpdateSchedule        string
	backupScheduleUpdateRetentionPolicy string
	backupScheduleUpdateLabels          []string
	backupScheduleUpdateAnnotations     []string
)

// backupScheduleUpdateCmd represents the backup-schedules update command
var backupScheduleUpdateCmd = &cobra.Command{
	Use:               "update",
	Short:             "Update a database backup schedule",
	Long:              "Update properties of an existing backup schedule",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completion.CompleteDbClusterID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusterIdentity := args[0]
		scheduleIdentity := args[1]

		// Get current schedule
		current, err := client.DBaaS().GetDbBackupSchedule(cmd.Context(), clusterIdentity, scheduleIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("backup schedule not found: %s", scheduleIdentity)
			}
			return fmt.Errorf("failed to get backup schedule: %w", err)
		}

		req := dbaas.UpdateDbBackupScheduleRequest{
			Name:            current.Name,
			Description:     "",
			Schedule:        current.Schedule,
			RetentionPolicy: current.RetentionPolicy,
			Labels:          current.Labels,
			Annotations:     current.Annotations,
		}

		if current.Description != nil {
			req.Description = *current.Description
		}

		// Update name if provided
		if cmd.Flags().Changed("name") {
			req.Name = backupScheduleUpdateName
		}

		// Update description if provided
		if cmd.Flags().Changed("description") {
			req.Description = backupScheduleUpdateDescription
		}

		// Update schedule if provided
		if cmd.Flags().Changed("schedule") {
			req.Schedule = backupScheduleUpdateSchedule
		}

		// Update retention policy if provided
		if cmd.Flags().Changed("retention-policy") {
			req.RetentionPolicy = backupScheduleUpdateRetentionPolicy
		}

		// Parse labels from key=value format
		if cmd.Flags().Changed("labels") {
			labels := make(map[string]string)
			for _, label := range backupScheduleUpdateLabels {
				parts := strings.SplitN(label, "=", 2)
				if len(parts) == 2 {
					labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Labels = labels
		}

		// Parse annotations from key=value format
		if cmd.Flags().Changed("annotations") {
			annotations := make(map[string]string)
			for _, annotation := range backupScheduleUpdateAnnotations {
				parts := strings.SplitN(annotation, "=", 2)
				if len(parts) == 2 {
					annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Annotations = annotations
		}

		schedule, err := client.DBaaS().UpdateDbBackupSchedule(cmd.Context(), clusterIdentity, scheduleIdentity, req)
		if err != nil {
			return fmt.Errorf("failed to update backup schedule: %w", err)
		}

		// Output in table format
		clusterName := ""
		if schedule.DbCluster != nil {
			clusterName = schedule.DbCluster.Name
		}

		nextBackup := "-"
		if schedule.NextBackupAt != nil {
			nextBackup = schedule.NextBackupAt.Format("2006-01-02 15:04:05")
		}

		body := [][]string{
			{
				schedule.Identity,
				schedule.Name,
				clusterName,
				string(schedule.Method),
				schedule.Schedule,
				schedule.RetentionPolicy,
				nextBackup,
				string(schedule.Status),
			},
		}

		backupScheduleUpdateNoHeader, _ := cmd.Flags().GetBool(NoHeaderKey)
		if backupScheduleUpdateNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Cluster", "Method", "Schedule", "Retention", "Next Backup", "Status"}, body)
		}

		return nil
	},
}

func init() {
	BackupSchedulesCmd.AddCommand(backupScheduleUpdateCmd)

	backupScheduleUpdateCmd.Flags().Bool(NoHeaderKey, false, "Do not print the header")
	backupScheduleUpdateCmd.Flags().StringVar(&backupScheduleUpdateName, "name", "", "Name of the backup schedule")
	backupScheduleUpdateCmd.Flags().StringVar(&backupScheduleUpdateDescription, "description", "", "Description of the backup schedule")
	backupScheduleUpdateCmd.Flags().StringVar(&backupScheduleUpdateSchedule, "schedule", "", "Cron expression for the backup schedule")
	backupScheduleUpdateCmd.Flags().StringVar(&backupScheduleUpdateRetentionPolicy, "retention-policy", "", "Retention policy for the backup schedule")
	backupScheduleUpdateCmd.Flags().StringSliceVar(&backupScheduleUpdateLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	backupScheduleUpdateCmd.Flags().StringSliceVar(&backupScheduleUpdateAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
}
