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
	attachInstanceID string
)

// attachCmd represents the attach command
var attachCmd = &cobra.Command{
	Use:               "attach <volume-id> [<volume-id> ...]",
	Short:             "Attach volume(s) to a virtual machine",
	Long:              "Attach one or more volumes to a virtual machine instance by identity.",
	Aliases:           []string{"att"},
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: completion.CompleteVolumeID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if attachInstanceID == "" {
			return fmt.Errorf("--instance is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// verify that the instance exists
		vmi, err := client.IaaS().GetMachine(cmd.Context(), attachInstanceID)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("instance %s not found", attachInstanceID)
			}
			return fmt.Errorf("failed to get instance: %w", err)
		}

		for _, volumeIdentity := range args {
			volume, err := client.IaaS().GetVolume(cmd.Context(), volumeIdentity)
			if err != nil {
				return fmt.Errorf("failed to get volume: %w", err)
			}

			fmt.Printf("Attaching volume: %s (%s) to instance %s\n", volume.Name, volume.Identity, attachInstanceID)

			req := iaas.AttachVolumeRequest{ResourceIdentity: vmi.Identity, ResourceType: "cloud_machine"}
			_, err = client.IaaS().AttachVolume(cmd.Context(), volumeIdentity, req)
			if err != nil {
				return fmt.Errorf("failed to attach volume: %w", err)
			}

			fmt.Println("Volume attached successfully")
		}
		return nil
	},
}

func init() {
	VolumesCmd.AddCommand(attachCmd)
	attachCmd.Flags().StringVar(&attachInstanceID, "instance", "", "Virtual machine instance identity")
	_ = attachCmd.RegisterFlagCompletionFunc("instance", completion.CompleteMachineID)
}
