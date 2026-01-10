package volumes

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	createVolumeName             string
	createVolumeDescription      string
	createVolumeRegion           string
	createVolumeSize             int
	createVolumeType             string
	createVolumeLabels           []string
	createVolumeAnnotations      []string
	createVolumeDeleteProtection bool
	createVolumeWait             bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a volume",
	Long:  "Create a new storage volume. The volume can be attached to machines after creation.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if createVolumeName == "" {
			return fmt.Errorf("name is required")
		}
		if createVolumeRegion == "" {
			return fmt.Errorf("region is required")
		}
		if createVolumeSize <= 0 {
			return fmt.Errorf("size must be greater than 0")
		}

		// Resolve region
		region, err := client.IaaS().GetRegion(cmd.Context(), createVolumeRegion)
		if err != nil {
			return fmt.Errorf("failed to get region: %w", err)
		}

		// Resolve volume type (by name, slug, or identity)
		volumeTypes, err := client.IaaS().ListVolumeTypes(cmd.Context(), &iaas.ListVolumeTypesRequest{})
		if err != nil {
			return fmt.Errorf("failed to get volume type: %w", err)
		}
		var volumeTypeObj *iaas.VolumeType
		for _, volumeType := range volumeTypes {
			if strings.EqualFold(volumeType.Name, createVolumeType) || strings.EqualFold(volumeType.Identity, createVolumeType) {
				volumeTypeObj = &volumeType
				break
			}
		}
		if volumeTypeObj == nil {
			return fmt.Errorf("volume type not found: %s", createVolumeType)
		}

		// Parse labels from key=value format
		labels := make(map[string]string)
		for _, label := range createVolumeLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Parse annotations from key=value format
		annotations := make(map[string]string)
		for _, annotation := range createVolumeAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		req := iaas.CreateVolume{
			Name:                createVolumeName,
			Description:         createVolumeDescription,
			CloudRegionIdentity: region.Identity,
			Size:                createVolumeSize,
			VolumeTypeIdentity:  volumeTypeObj.Identity,
			Labels:              labels,
			Annotations:         annotations,
			DeleteProtection:    createVolumeDeleteProtection,
		}

		volume, err := client.IaaS().CreateVolume(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create volume: %w", err)
		}

		if createVolumeWait {
			// Wait for volume to be available
			err = client.IaaS().WaitUntilVolumeIsAvailable(cmd.Context(), volume.Identity)
			if err != nil {
				return fmt.Errorf("failed to wait for volume to be available: %w", err)
			}
			// Refresh volume to get latest status
			volume, err = client.IaaS().GetVolume(cmd.Context(), volume.Identity)
			if err != nil {
				return fmt.Errorf("failed to get volume: %w", err)
			}
		}

		// Output in table format
		volumeType := ""
		if volume.VolumeType != nil {
			volumeType = volume.VolumeType.Name
		}
		body := [][]string{
			{
				volume.Identity,
				volume.Name,
				volume.Status,
				volume.Region.Name,
				volumeType,
				fmt.Sprintf("%dGB", volume.Size),
				formattime.FormatTime(volume.CreatedAt.Local(), false),
			},
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Status", "Region", "Type", "Size", "Age"}, body)
		}

		return nil
	},
}

func init() {
	VolumesCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createVolumeName, "name", "", "Name of the volume (required)")
	createCmd.Flags().StringVar(&createVolumeDescription, "description", "", "Description of the volume")
	createCmd.Flags().StringVar(&createVolumeRegion, "region", "", "Region of the volume (required)")
	createCmd.Flags().IntVar(&createVolumeSize, "size", 0, "Size of the volume in GB (required)")
	createCmd.Flags().StringVar(&createVolumeType, "type", "block", "Volume type")
	createCmd.Flags().StringSliceVar(&createVolumeLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	createCmd.Flags().StringSliceVar(&createVolumeAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	createCmd.Flags().BoolVar(&createVolumeDeleteProtection, "delete-protection", false, "Enable delete protection")
	createCmd.Flags().BoolVar(&createVolumeWait, "wait", false, "Wait for the volume to be available before returning")

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("region")
	_ = createCmd.MarkFlagRequired("size")
}
