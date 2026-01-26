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
	backupScheduleCreateName            string
	backupScheduleCreateDescription     string
	backupScheduleCreateSchedule        string
	backupScheduleCreateRetentionPolicy string
	backupScheduleCreateMethod          string
	backupScheduleCreateLabels          []string
	backupScheduleCreateAnnotations     []string
)

// backupScheduleCreateCmd represents the backup-schedules create command
var backupScheduleCreateCmd = &cobra.Command{
	Use:               "create",
	Short:             "Create a database backup schedule",
	Long:              "Create a new backup schedule for a database cluster",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDbClusterID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusterIdentity := args[0]

		if backupScheduleCreateName == "" {
			return fmt.Errorf("name is required")
		}
		if backupScheduleCreateSchedule == "" {
			return fmt.Errorf("schedule is required")
		}
		if backupScheduleCreateRetentionPolicy == "" {
			return fmt.Errorf("retention-policy is required")
		}
		if backupScheduleCreateMethod == "" {
			return fmt.Errorf("method is required")
		}

		// Parse labels from key=value format
		labels := make(map[string]string)
		for _, label := range backupScheduleCreateLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Parse annotations from key=value format
		annotations := make(map[string]string)
		for _, annotation := range backupScheduleCreateAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		method := dbaas.DbClusterBackupScheduleMethod(backupScheduleCreateMethod)
		if method != dbaas.DbClusterBackupScheduleMethodSnapshot && method != dbaas.DbClusterBackupScheduleMethodBarman {
			return fmt.Errorf("method must be either 'snapshot' or 'barman'")
		}

		createReq := dbaas.CreateDbBackupScheduleRequest{
			Name:            backupScheduleCreateName,
			Schedule:        backupScheduleCreateSchedule,
			RetentionPolicy: backupScheduleCreateRetentionPolicy,
			Method:          method,
			Labels:          labels,
			Annotations:     annotations,
		}

		if backupScheduleCreateDescription != "" {
			createReq.Description = &backupScheduleCreateDescription
		}

		schedule, err := client.DBaaS().CreateDbBackupSchedule(cmd.Context(), clusterIdentity, createReq)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("database cluster not found: %s", clusterIdentity)
			}
			return fmt.Errorf("failed to create backup schedule: %w", err)
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

		backupScheduleCreateNoHeader, _ := cmd.Flags().GetBool(NoHeaderKey)
		if backupScheduleCreateNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Cluster", "Method", "Schedule", "Retention", "Next Backup", "Status"}, body)
		}

		return nil
	},
}

func init() {
	BackupSchedulesCmd.AddCommand(backupScheduleCreateCmd)

	backupScheduleCreateCmd.Flags().Bool(NoHeaderKey, false, "Do not print the header")
	backupScheduleCreateCmd.Flags().StringVarP(&backupScheduleCreateName, "name", "n", "", "Name of the backup schedule (required)")
	backupScheduleCreateCmd.Flags().StringVar(&backupScheduleCreateDescription, "description", "", "Description of the backup schedule")
	backupScheduleCreateCmd.Flags().StringVar(&backupScheduleCreateSchedule, "schedule", "", "Cron expression for the backup schedule (required)")
	backupScheduleCreateCmd.Flags().StringVar(&backupScheduleCreateRetentionPolicy, "retention-policy", "", "Retention policy for the backup schedule (required)")
	backupScheduleCreateCmd.Flags().StringVar(&backupScheduleCreateMethod, "method", "barman", "Backup method: 'barman' (default)")
	backupScheduleCreateCmd.Flags().StringSliceVar(&backupScheduleCreateLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	backupScheduleCreateCmd.Flags().StringSliceVar(&backupScheduleCreateAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")

	_ = backupScheduleCreateCmd.MarkFlagRequired("name")
	_ = backupScheduleCreateCmd.MarkFlagRequired("schedule")
	_ = backupScheduleCreateCmd.MarkFlagRequired("retention-policy")
}
