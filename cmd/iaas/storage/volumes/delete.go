package volumes

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	ConfirmDelete bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Delete a volume",
	Long:              "Delete a volume by its identity.",
	Aliases:           []string{"d", "rm", "del", "remove"},
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: completion.CompleteVolumeID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("volume identity is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		volumesToDelete := []*iaas.Volume{}
		for _, volumeIdentity := range args {
			volume, err := client.IaaS().GetVolume(cmd.Context(), volumeIdentity)
			if err != nil {
				if tcclient.IsNotFound(err) {
					fmt.Printf("Volume %s not found\n", volumeIdentity)
					continue
				}
				return fmt.Errorf("failed to get volume: %w", err)
			}
			volumesToDelete = append(volumesToDelete, volume)
		}

		if !ConfirmDelete {
			// ask for confirmation before deleting
			fmt.Printf("Are you sure you want to delete the following volumes?\n")
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

		for _, volume := range volumesToDelete {
			fmt.Printf("Deleting volume: %s (%s)\n", volume.Name, volume.Identity)
			if err := client.IaaS().DeleteVolume(cmd.Context(), volume.Identity); err != nil {
				return fmt.Errorf("failed to delete volume: %w", err)
			}
			fmt.Println("Volume deleted successfully")
		}
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&ConfirmDelete, "force", false, "Force the deletion and skip the confirmation")

	VolumesCmd.AddCommand(deleteCmd)
}
