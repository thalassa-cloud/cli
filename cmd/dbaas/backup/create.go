package backup

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupCreateName            string
	backupCreateDescription     string
	backupCreateLabels          []string
	backupCreateAnnotations     []string
	backupCreateRetentionPolicy string
	backupCreateWait            bool
)

// backupCreateCmd represents the backup create command
var backupCreateCmd = &cobra.Command{
	Use:               "create",
	Short:             "Create a database backup",
	Long:              "Create a new backup for a database cluster",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDbClusterID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusterIdentity := args[0]

		if backupCreateName == "" {
			return fmt.Errorf("name is required")
		}

		// Parse labels from key=value format
		labels := make(map[string]string)
		for _, label := range backupCreateLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Parse annotations from key=value format
		annotations := make(map[string]string)
		for _, annotation := range backupCreateAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		createReq := dbaas.CreateDbClusterBackupRequest{
			Name:        backupCreateName,
			Labels:      labels,
			Annotations: annotations,
		}

		if backupCreateDescription != "" {
			createReq.Description = &backupCreateDescription
		}

		if backupCreateRetentionPolicy != "" {
			createReq.RetentionPolicy = &backupCreateRetentionPolicy
		}

		backup, err := client.DBaaS().CreateDbBackup(cmd.Context(), clusterIdentity, createReq)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("database cluster not found: %s", clusterIdentity)
			}
			return fmt.Errorf("failed to create backup: %w", err)
		}

		if backupCreateWait {
			// Poll until backup is completed
			for {
				backup, err = client.DBaaS().GetDbBackup(cmd.Context(), backup.Identity)
				if err != nil {
					return fmt.Errorf("failed to get backup: %w", err)
				}
				// Check if backup is completed (has StoppedAt timestamp)
				if backup.StoppedAt != nil {
					break
				}
				// Check for failed status
				if backup.Status == dbaas.ObjectStatusFailed {
					return fmt.Errorf("backup creation failed: %s", backup.StatusMessage)
				}
				// Simple polling with sleep
				select {
				case <-cmd.Context().Done():
					return cmd.Context().Err()
				case <-time.After(5 * time.Second):
					// Continue polling
				}
			}
		}

		// Output in table format
		clusterName := ""
		if backup.DbCluster != nil {
			clusterName = backup.DbCluster.Name
		}

		body := [][]string{
			{
				backup.Identity,
				clusterName,
				string(backup.EngineType),
				backup.EngineVersion,
				backup.BackupType,
				string(backup.Status),
			},
		}

		var backupCreateNoHeader bool
		if cmd.Flags().Changed(NoHeaderKey) {
			backupCreateNoHeader, _ = cmd.Flags().GetBool(NoHeaderKey)
		}
		if backupCreateNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Cluster", "Engine", "Version", "Type", "Status"}, body)
		}

		return nil
	},
}

func init() {
	BackupCmd.AddCommand(backupCreateCmd)

	backupCreateCmd.Flags().Bool(NoHeaderKey, false, "Do not print the header")
	backupCreateCmd.Flags().StringVar(&backupCreateName, "name", "", "Name of the backup (required)")
	backupCreateCmd.Flags().StringVar(&backupCreateDescription, "description", "", "Description of the backup")
	backupCreateCmd.Flags().StringSliceVar(&backupCreateLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	backupCreateCmd.Flags().StringSliceVar(&backupCreateAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	backupCreateCmd.Flags().StringVar(&backupCreateRetentionPolicy, "retention-policy", "", "Retention policy for the backup")
	backupCreateCmd.Flags().BoolVar(&backupCreateWait, "wait", false, "Wait for the backup to be completed before returning")

	_ = backupCreateCmd.MarkFlagRequired("name")
}
