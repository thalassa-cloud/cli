package snapshots

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	createVolumeID         string
	createDescription      string
	createLabels           []string
	createAnnotations      []string
	createDeleteProtection bool

	waitForReady bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create <name>",
	Short:   "Create a snapshot",
	Long:    "Create a snapshot from a volume by providing a name and volume ID.",
	Aliases: []string{"c", "new"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if createVolumeID == "" {
			return fmt.Errorf("--volume is required")
		}

		snapshotName := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// verify that the volume exists
		volume, err := client.IaaS().GetVolume(cmd.Context(), createVolumeID)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("volume %s not found", createVolumeID)
			}
			return fmt.Errorf("failed to get volume: %w", err)
		}

		fmt.Printf("Creating snapshot: %s from volume %s (%s)\n", snapshotName, volume.Name, volume.Identity)

		// Parse labels from key=value format
		labels := make(map[string]string)
		for _, label := range createLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[parts[0]] = parts[1]
			}
		}

		// Parse annotations from key=value format
		annotations := make(map[string]string)
		for _, annotation := range createAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[parts[0]] = parts[1]
			}
		}

		req := iaas.CreateSnapshotRequest{
			Name:             snapshotName,
			Description:      createDescription,
			Labels:           labels,
			Annotations:      annotations,
			VolumeIdentity:   createVolumeID,
			DeleteProtection: createDeleteProtection,
		}
		snapshot, err := client.IaaS().CreateSnapshot(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create snapshot: %w", err)
		}

		fmt.Printf("Snapshot created successfully: %s (%s)\n", snapshot.Name, snapshot.Identity)

		if waitForReady {
			fmt.Println("Waiting for snapshot to be ready...")
			err = client.IaaS().WaitUntilSnapshotIsAvailable(cmd.Context(), snapshot.Identity)
			if err != nil {
				return fmt.Errorf("failed to wait for snapshot to be ready: %w", err)
			}
			fmt.Println("Snapshot is ready")
		}
		return nil
	},
}

func init() {
	SnapshotsCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createVolumeID, "volume", "", "Volume identity to create snapshot from")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description of the snapshot")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	createCmd.Flags().BoolVar(&createDeleteProtection, "delete-protection", false, "Enable delete protection for the snapshot")
	createCmd.Flags().BoolVar(&waitForReady, "wait", false, "Wait for the snapshot to be ready for use")

	_ = createCmd.RegisterFlagCompletionFunc("volume", completion.CompleteVolumeID)
}
