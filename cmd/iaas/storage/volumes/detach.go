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
	ConfirmDetach bool
)

// detachCmd represents the detach command
var detachCmd = &cobra.Command{
	Use:               "detach",
	Short:             "Detach a volume",
	Long:              "Detach a volume from any current attachment target by its identity.",
	Aliases:           []string{"det", "disassociate"},
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

		volumesToDetach := []*iaas.Volume{}
		for _, volumeIdentity := range args {
			volume, err := client.IaaS().GetVolume(cmd.Context(), volumeIdentity)
			if err != nil {
				if tcclient.IsNotFound(err) {
					fmt.Printf("Volume %s not found\n", volumeIdentity)
					continue
				}
				return fmt.Errorf("failed to get volume: %w", err)
			}
			volumesToDetach = append(volumesToDetach, volume)
		}

		if !ConfirmDetach {
			// ask for confirmation before deleting
			fmt.Printf("Are you sure you want to detach the following volumes?\n")
			for _, volume := range volumesToDetach {
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

		for _, volume := range volumesToDetach {
			fmt.Printf("Detaching volume: %s (%s)\n", volume.Name, volume.Identity)
			for _, attachment := range volume.Attachments {
				if err := client.IaaS().DetachVolume(cmd.Context(), volume.Identity, iaas.DetachVolumeRequest{
					ResourceIdentity: attachment.AttachedToIdentity,
					ResourceType:     attachment.AttachedToResourceType,
				}); err != nil {
					return fmt.Errorf("failed to detach volume: %w", err)
				}
			}

			fmt.Println("Volume detached successfully")
		}
		return nil
	},
}

func init() {
	detachCmd.Flags().BoolVar(&ConfirmDetach, "force", false, "Force the detachment and skip the confirmation")

	VolumesCmd.AddCommand(detachCmd)
}
