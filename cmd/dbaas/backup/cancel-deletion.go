package backup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

// backupCancelDeletionCmd represents the backup cancel-deletion command
var backupCancelDeletionCmd = &cobra.Command{
	Use:     "cancel-deletion",
	Short:   "Cancel scheduled deletion of a backup",
	Long:    "Cancel the scheduled deletion of a database backup",
	Aliases: []string{"cancel", "undelete"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		backupIdentity := args[0]

		err = client.DBaaS().CancelDeleteDbBackup(cmd.Context(), backupIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("backup not found: %s", backupIdentity)
			}
			return fmt.Errorf("failed to cancel backup deletion: %w", err)
		}

		fmt.Printf("Backup deletion cancelled for %s\n", backupIdentity)
		return nil
	},
}

func init() {
	BackupCmd.AddCommand(backupCancelDeletionCmd)
}
