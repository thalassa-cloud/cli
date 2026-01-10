package volumes

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	resizeSize int
	resizeWait bool
)

// resizeCmd represents the resize command
var resizeCmd = &cobra.Command{
	Use:               "resize <volume-id>",
	Short:             "Resize a volume",
	Long:              "Resize a volume to a new size in GB. The new size must be larger than the current size.",
	Aliases:           []string{"res"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteVolumeID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resizeSize <= 0 {
			return fmt.Errorf("--size must be greater than 0")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		volumeIdentity := args[0]

		// Get the current volume to check its size
		volume, err := client.IaaS().GetVolume(cmd.Context(), volumeIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("volume %s not found", volumeIdentity)
			}
			return fmt.Errorf("failed to get volume: %w", err)
		}

		if resizeSize < volume.Size {
			return fmt.Errorf("new size (%dGB) must be greater than or equal to current size (%dGB)", resizeSize, volume.Size)
		}

		if resizeSize == volume.Size {
			return fmt.Errorf("volume is already %dGB", volume.Size)
		}

		req := iaas.UpdateVolume{
			Name:             volume.Name,
			Description:      volume.Description,
			Labels:           volume.Labels,
			Annotations:      volume.Annotations,
			DeleteProtection: volume.DeleteProtection,
			Size:             resizeSize,
		}
		updatedVolume, err := client.IaaS().UpdateVolume(cmd.Context(), volumeIdentity, req)
		if err != nil {
			return fmt.Errorf("failed to resize volume: %w", err)
		}

		if resizeWait {
			timeout := time.After(10 * time.Minute)
			tick := time.Tick(2 * time.Second)
			for {
				updatedVolume, err = client.IaaS().GetVolume(cmd.Context(), volumeIdentity)
				if err != nil {
					return fmt.Errorf("failed to get volume: %w", err)
				}
				if strings.EqualFold(updatedVolume.Status, "available") || strings.EqualFold(updatedVolume.Status, "attached") {
					break
				}
				select {
				case <-timeout:
					return fmt.Errorf("timeout waiting for volume to be in ready state")
				case <-tick:
					// continue looping
				}
			}
		}

		// Output in table format
		volumeType := ""
		if updatedVolume.VolumeType != nil {
			volumeType = updatedVolume.VolumeType.Name
		}
		body := [][]string{
			{
				updatedVolume.Identity,
				updatedVolume.Name,
				updatedVolume.Status,
				updatedVolume.Region.Name,
				volumeType,
				fmt.Sprintf("%dGB", updatedVolume.Size),
				formattime.FormatTime(updatedVolume.CreatedAt.Local(), false),
			},
		}
		table.Print([]string{"ID", "Name", "Status", "Region", "Type", "Size", "Age"}, body)

		return nil
	},
}

func init() {
	VolumesCmd.AddCommand(resizeCmd)

	resizeCmd.Flags().IntVar(&resizeSize, "size", 0, "New size in GB (required)")
	resizeCmd.Flags().BoolVar(&resizeWait, "wait", false, "Wait for the resize operation to complete")
	_ = resizeCmd.MarkFlagRequired("size")
}
