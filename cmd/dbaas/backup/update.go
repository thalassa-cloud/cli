package backup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupUpdateDeleteProtection bool
	backupUpdateNoHeader         bool
)

var backupUpdateCmd = &cobra.Command{
	Use:               "update",
	Short:             "Update a database backup",
	Long:              "Update backup properties such as delete protection",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDbBackupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("delete-protection") {
			return fmt.Errorf("--delete-protection is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		backupIdentity := args[0]
		deleteProtection := backupUpdateDeleteProtection
		backup, err := client.DBaaS().UpdateDbBackup(cmd.Context(), backupIdentity, dbaas.UpdateDbClusterBackupRequest{
			DeleteProtection: &deleteProtection,
		})
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("backup not found: %s", backupIdentity)
			}
			return fmt.Errorf("failed to update backup: %w", err)
		}

		clusterName := ""
		if backup.DbCluster != nil {
			clusterName = backup.DbCluster.Name
		}

		body := [][]string{{
			backup.Identity,
			clusterName,
			string(backup.Status),
			fmt.Sprintf("%v", backup.DeleteProtection),
		}}
		if backupUpdateNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Cluster", "Status", "Delete Protection"}, body)
		}

		return nil
	},
}

func init() {
	BackupCmd.AddCommand(backupUpdateCmd)

	backupUpdateCmd.Flags().BoolVar(&backupUpdateNoHeader, NoHeaderKey, false, "Do not print the header")
	backupUpdateCmd.Flags().BoolVar(&backupUpdateDeleteProtection, "delete-protection", false, "Enable or disable delete protection")
}
