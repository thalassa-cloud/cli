package volumes

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
	Short:             "Delete volume(s)",
	Long:              "Delete volume(s) by identity or label selector.",
	Aliases:           []string{"d", "rm", "del", "remove"},
	Args:              cobra.MinimumNArgs(0),
	ValidArgsFunction: completion.CompleteVolumeID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && labelSelector == "" {
			return fmt.Errorf("either volume identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect volumes to delete
		volumesToDelete := []iaas.Volume{}

		// If label selector is provided, filter by labels
		if labelSelector != "" {
			allVolumes, err := client.IaaS().ListVolumes(cmd.Context(), &iaas.ListVolumesRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(labelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list volumes: %w", err)
			}
			if len(allVolumes) == 0 {
				fmt.Println("No volumes found matching the label selector")
				return nil
			}
			volumesToDelete = append(volumesToDelete, allVolumes...)
		} else {
			// Get volumes by identity
			for _, volumeIdentity := range args {
				volume, err := client.IaaS().GetVolume(cmd.Context(), volumeIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("Volume %s not found\n", volumeIdentity)
						continue
					}
					return fmt.Errorf("failed to get volume: %w", err)
				}
				volumesToDelete = append(volumesToDelete, *volume)
			}
		}

		if len(volumesToDelete) == 0 {
			fmt.Println("No volumes to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !force {
			fmt.Printf("Are you sure you want to delete the following volume(s)?\n")
			for _, volume := range volumesToDelete {
				fmt.Printf("  %s (%s)\n", volume.Name, volume.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each volume
		for _, volume := range volumesToDelete {
			fmt.Printf("Deleting volume: %s (%s)\n", volume.Name, volume.Identity)
			err := client.IaaS().DeleteVolume(cmd.Context(), volume.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete volume: %w", err)
			}

			if wait {
				if err := client.IaaS().WaitUntilVolumeIsDeleted(cmd.Context(), volume.Identity); err != nil {
					return fmt.Errorf("failed to wait for volume to be deleted: %w", err)
				}
			}
			fmt.Printf("Volume %s deleted successfully\n", volume.Identity)
		}

		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&wait, "wait", false, "Wait for the volume(s) to be deleted")
	deleteCmd.Flags().BoolVar(&force, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVar(&labelSelector, "selector", "", "Label selector to filter volumes (format: key1=value1,key2=value2)")

	VolumesCmd.AddCommand(deleteCmd)
}
