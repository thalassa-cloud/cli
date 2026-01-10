package volumes

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	resizeSize          int
	resizeWait          bool
	resizeForce         bool
	resizeLabelSelector string
)

// resizeCmd represents the resize command
var resizeCmd = &cobra.Command{
	Use:               "resize [volume-id...]",
	Short:             "Resize volume(s)",
	Long:              "Resize volume(s) to a new size in GB. The new size must be larger than the current size. Can resize multiple volumes by identity or using a label selector.",
	Aliases:           []string{"res"},
	Args:              cobra.MinimumNArgs(0),
	ValidArgsFunction: completion.CompleteVolumeID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resizeSize <= 0 {
			return fmt.Errorf("--size must be greater than 0")
		}

		if len(args) == 0 && resizeLabelSelector == "" {
			return fmt.Errorf("either volume identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect volumes to resize
		volumesToResize := []iaas.Volume{}

		// If label selector is provided, filter by labels
		if resizeLabelSelector != "" {
			allVolumes, err := client.IaaS().ListVolumes(cmd.Context(), &iaas.ListVolumesRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(resizeLabelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list volumes: %w", err)
			}
			if len(allVolumes) == 0 {
				fmt.Println("No volumes found matching the label selector")
				return nil
			}
			volumesToResize = append(volumesToResize, allVolumes...)
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
				volumesToResize = append(volumesToResize, *volume)
			}
		}

		if len(volumesToResize) == 0 {
			fmt.Println("No volumes to resize")
			return nil
		}

		// Validate sizes for all volumes
		var invalidVolumes []string
		for _, volume := range volumesToResize {
			if resizeSize < volume.Size {
				invalidVolumes = append(invalidVolumes, fmt.Sprintf("%s (current: %dGB)", volume.Identity, volume.Size))
			} else if resizeSize == volume.Size {
				invalidVolumes = append(invalidVolumes, fmt.Sprintf("%s (already %dGB)", volume.Identity, volume.Size))
			}
		}
		if len(invalidVolumes) > 0 {
			return fmt.Errorf("cannot resize the following volume(s): %s", strings.Join(invalidVolumes, ", "))
		}

		// Ask for confirmation unless --force is provided
		if !resizeForce {
			fmt.Printf("Are you sure you want to resize the following volume(s) to %dGB?\n", resizeSize)
			for _, volume := range volumesToResize {
				fmt.Printf("  %s (%s) - current size: %dGB\n", volume.Name, volume.Identity, volume.Size)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Resize each volume
		var resizedVolumes []iaas.Volume
		for _, volume := range volumesToResize {
			req := iaas.UpdateVolume{
				Name:             volume.Name,
				Description:      volume.Description,
				Labels:           volume.Labels,
				Annotations:      volume.Annotations,
				DeleteProtection: volume.DeleteProtection,
				Size:             resizeSize,
			}
			updatedVolume, err := client.IaaS().UpdateVolume(cmd.Context(), volume.Identity, req)
			if err != nil {
				return fmt.Errorf("failed to resize volume %s: %w", volume.Identity, err)
			}

			if resizeWait {
				timeout := time.After(10 * time.Minute)
				tick := time.Tick(2 * time.Second)
				for {
					updatedVolume, err = client.IaaS().GetVolume(cmd.Context(), volume.Identity)
					if err != nil {
						return fmt.Errorf("failed to get volume: %w", err)
					}
					if strings.EqualFold(updatedVolume.Status, "available") || strings.EqualFold(updatedVolume.Status, "attached") {
						break
					}
					select {
					case <-timeout:
						return fmt.Errorf("timeout waiting for volume %s to be in ready state", volume.Identity)
					case <-tick:
						// continue looping
					}
				}
			}

			resizedVolumes = append(resizedVolumes, *updatedVolume)
		}

		// Output in table format
		body := make([][]string, 0, len(resizedVolumes))
		for _, updatedVolume := range resizedVolumes {
			volumeType := ""
			if updatedVolume.VolumeType != nil {
				volumeType = updatedVolume.VolumeType.Name
			}
			body = append(body, []string{
				updatedVolume.Identity,
				updatedVolume.Name,
				updatedVolume.Status,
				updatedVolume.Region.Name,
				volumeType,
				fmt.Sprintf("%dGB", updatedVolume.Size),
				formattime.FormatTime(updatedVolume.CreatedAt.Local(), false),
			})
		}
		table.Print([]string{"ID", "Name", "Status", "Region", "Type", "Size", "Age"}, body)

		return nil
	},
}

func init() {
	VolumesCmd.AddCommand(resizeCmd)

	resizeCmd.Flags().IntVar(&resizeSize, "size", 0, "New size in GB (required)")
	resizeCmd.Flags().BoolVar(&resizeWait, "wait", false, "Wait for the resize operation to complete")
	resizeCmd.Flags().BoolVar(&resizeForce, "force", false, "Force the resize and skip the confirmation")
	resizeCmd.Flags().StringVarP(&resizeLabelSelector, "selector", "l", "", "Label selector to filter volumes (format: key1=value1,key2=value2)")
	_ = resizeCmd.MarkFlagRequired("size")
}
