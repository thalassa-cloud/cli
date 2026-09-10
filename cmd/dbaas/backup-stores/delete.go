package backupstores

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupStoreDeleteForce bool
)

var backupStoreDeleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Delete a database backup store",
	Long:              "Delete a backup store by identity",
	Aliases:           []string{"d", "rm", "del", "remove"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDbBackupStoreID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		storeIdentity := args[0]
		storeName := storeIdentity
		if !backupStoreDeleteForce {
			store, err := client.DBaaS().GetDbObjectStore(cmd.Context(), storeIdentity)
			if err != nil {
				if tcclient.IsNotFound(err) {
					return fmt.Errorf("backup store not found: %s", storeIdentity)
				}
				return fmt.Errorf("failed to get backup store: %w", err)
			}
			storeName = store.Name
		}

		var summary strings.Builder
		fmt.Fprintf(&summary, "Are you sure you want to delete backup store %s (%s)?\n", storeName, storeIdentity)
		proceed, err := shared.PromptDestructiveUnlessForce(backupStoreDeleteForce, summary.String())
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}

		if err := client.DBaaS().DeleteDbObjectStore(cmd.Context(), storeIdentity); err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("backup store not found: %s", storeIdentity)
			}
			return fmt.Errorf("failed to delete backup store: %w", err)
		}

		fmt.Printf("Backup store %s deleted successfully\n", storeIdentity)
		return nil
	},
}

func init() {
	BackupStoresCmd.AddCommand(backupStoreDeleteCmd)

	backupStoreDeleteCmd.Flags().BoolVar(&backupStoreDeleteForce, "force", false, "Force the deletion and skip the confirmation")
}
