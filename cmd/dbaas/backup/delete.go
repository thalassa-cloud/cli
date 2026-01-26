package backup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/filters"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupDeleteForce      bool
	backupDeleteAllFailed  bool
	backupDeleteLabelSelector string
)

// backupDeleteCmd represents the backup delete command
var backupDeleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Delete database backup(s)",
	Long:              "Delete database backup(s) by identity, label selector, or all failed backups",
	Aliases:           []string{"d", "rm", "del", "remove"},
	Args:              cobra.MinimumNArgs(0),
	ValidArgsFunction: completion.CompleteDbBackupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && !backupDeleteAllFailed && backupDeleteLabelSelector == "" {
			return fmt.Errorf("either backup identity(ies), --all-failed, or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect backups to delete
		backupsToDelete := []dbaas.DbClusterBackup{}

		// Build filters for listing backups
		f := []filters.Filter{}

		// Add label selector filter if provided
		if backupDeleteLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(backupDeleteLabelSelector),
			})
		}

		// Add status filter for --all-failed
		if backupDeleteAllFailed {
			f = append(f, &filters.FilterKeyValue{
				Key:   "status",
				Value: "failed",
			})
		}

		// If using --all-failed or --selector, list backups with filters
		if backupDeleteAllFailed || backupDeleteLabelSelector != "" {
			listRequest := &dbaas.ListDbBackupsRequest{
				Filters: f,
			}
			allBackups, err := client.DBaaS().ListDbBackupsForOrganisation(cmd.Context(), listRequest)
			if err != nil {
				return fmt.Errorf("failed to list backups: %w", err)
			}
			if len(allBackups) == 0 {
				if backupDeleteAllFailed && backupDeleteLabelSelector != "" {
					fmt.Println("No failed backups found matching the label selector")
				} else if backupDeleteAllFailed {
					fmt.Println("No failed backups found")
				} else {
					fmt.Println("No backups found matching the label selector")
				}
				return nil
			}
			backupsToDelete = append(backupsToDelete, allBackups...)
		}

		// Get backups by identity (if provided)
		for _, backupIdentity := range args {
			backup, err := client.DBaaS().GetDbBackup(cmd.Context(), backupIdentity)
			if err != nil {
				if tcclient.IsNotFound(err) {
					fmt.Printf("Backup %s not found\n", backupIdentity)
					continue
				}
				return fmt.Errorf("failed to get backup: %w", err)
			}
			// Check if already added (avoid duplicates when combining with filters)
			alreadyAdded := false
			for _, existing := range backupsToDelete {
				if existing.Identity == backup.Identity {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				backupsToDelete = append(backupsToDelete, *backup)
			}
		}

		if len(backupsToDelete) == 0 {
			fmt.Println("No backups to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !backupDeleteForce {
			if len(backupsToDelete) == 1 {
				fmt.Printf("Are you sure you want to delete backup %s?\n", backupsToDelete[0].Identity)
			} else {
				fmt.Printf("Are you sure you want to delete the following backup(s)?\n")
				for _, backup := range backupsToDelete {
					fmt.Printf("  %s\n", backup.Identity)
				}
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each backup
		for _, backup := range backupsToDelete {
			fmt.Printf("Deleting backup: %s\n", backup.Identity)
			err := client.DBaaS().DeleteDbBackup(cmd.Context(), backup.Identity)
			if err != nil {
				if tcclient.IsNotFound(err) {
					fmt.Printf("Backup %s not found\n", backup.Identity)
					continue
				}
				return fmt.Errorf("failed to delete backup: %w", err)
			}
			fmt.Printf("Backup %s deleted successfully\n", backup.Identity)
		}

		return nil
	},
}

func init() {
	BackupCmd.AddCommand(backupDeleteCmd)

	backupDeleteCmd.Flags().BoolVar(&backupDeleteForce, "force", false, "Force the deletion and skip the confirmation")
	backupDeleteCmd.Flags().BoolVar(&backupDeleteAllFailed, "all-failed", false, "Delete all failed backups")
	backupDeleteCmd.Flags().StringVarP(&backupDeleteLabelSelector, "selector", "l", "", "Label selector to filter backups (format: key1=value1,key2=value2)")
}
