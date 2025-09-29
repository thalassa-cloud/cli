package snapshots

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	ConfirmDelete bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Delete a snapshot",
	Long:              "Delete a snapshot by its identity.",
	Aliases:           []string{"d", "rm", "del", "remove"},
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: completion.CompleteSnapshotID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("snapshot identity is required")
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		snapshotsToDelete := []*iaas.Snapshot{}
		for _, snapshotIdentity := range args {
			snapshot, err := client.IaaS().GetSnapshot(cmd.Context(), snapshotIdentity)
			if err != nil {
				return fmt.Errorf("failed to get snapshot: %w", err)
			}
			snapshotsToDelete = append(snapshotsToDelete, snapshot)
		}

		if !ConfirmDelete {
			// ask for confirmation before deleting
			fmt.Printf("Are you sure you want to delete the following snapshots?\n")
			for _, snapshot := range snapshotsToDelete {
				fmt.Printf("  %s (%s)\n", snapshot.Name, snapshot.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		for _, snapshot := range snapshotsToDelete {
			// Fetch the snapshot first to provide a friendly message
			fmt.Printf("Deleting snapshot: %s (%s)\n", snapshot.Name, snapshot.Identity)
			if err := client.IaaS().DeleteSnapshot(cmd.Context(), snapshot.Identity); err != nil {
				return fmt.Errorf("failed to delete snapshot: %w", err)
			}
			fmt.Println("Snapshot deleted successfully")
		}
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&ConfirmDelete, "force", false, "Force the deletion and skip the confirmation")

	SnapshotsCmd.AddCommand(deleteCmd)
}
