package snapshots

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	wait          bool
	force         bool
	labelSelector string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Delete snapshot(s)",
	Long:              "Delete snapshot(s) by identity or label selector.",
	Aliases:           []string{"d", "rm", "del", "remove"},
	Args:              cobra.MinimumNArgs(0),
	ValidArgsFunction: completion.CompleteSnapshotID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && labelSelector == "" {
			return fmt.Errorf("either snapshot identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect snapshots to delete
		snapshotsToDelete := []iaas.Snapshot{}

		// If label selector is provided, filter by labels
		if labelSelector != "" {
			allSnapshots, err := client.IaaS().ListSnapshots(cmd.Context(), &iaas.ListSnapshotsRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(labelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list snapshots: %w", err)
			}
			if len(allSnapshots) == 0 {
				fmt.Println("No snapshots found matching the label selector")
				return nil
			}
			snapshotsToDelete = append(snapshotsToDelete, allSnapshots...)
		} else {
			// Get snapshots by identity
			for _, snapshotIdentity := range args {
				snapshot, err := client.IaaS().GetSnapshot(cmd.Context(), snapshotIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("Snapshot %s not found\n", snapshotIdentity)
						continue
					}
					return fmt.Errorf("failed to get snapshot: %w", err)
				}
				snapshotsToDelete = append(snapshotsToDelete, *snapshot)
			}
		}

		if len(snapshotsToDelete) == 0 {
			fmt.Println("No snapshots to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !force {
			fmt.Printf("Are you sure you want to delete the following snapshot(s)?\n")
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

		// Delete each snapshot
		for _, snapshot := range snapshotsToDelete {
			fmt.Printf("Deleting snapshot: %s (%s)\n", snapshot.Name, snapshot.Identity)
			err := client.IaaS().DeleteSnapshot(cmd.Context(), snapshot.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete snapshot: %w", err)
			}

			if wait {
				if err := client.IaaS().WaitUntilSnapshotIsDeleted(cmd.Context(), snapshot.Identity); err != nil {
					return fmt.Errorf("failed to wait for snapshot to be deleted: %w", err)
				}
			}
			fmt.Printf("Snapshot %s deleted successfully\n", snapshot.Identity)
		}

		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&wait, "wait", false, "Wait for the snapshot(s) to be deleted")
	deleteCmd.Flags().BoolVar(&force, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVar(&labelSelector, "selector", "", "Label selector to filter snapshots (format: key1=value1,key2=value2)")

	SnapshotsCmd.AddCommand(deleteCmd)
}
