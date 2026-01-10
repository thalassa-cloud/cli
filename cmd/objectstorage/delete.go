package objectstorage

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	deleteForce   bool
	deleteWait    bool
	deleteTimeout string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete an object storage bucket",
	Long:    "Delete an object storage bucket by name. This will permanently delete the bucket and all its contents.",
	Example: "tcloud storage object-storage delete my-bucket\ntcloud storage object-storage delete my-bucket --force",
	Aliases: []string{"d", "del", "remove", "rm", "delete-bucket"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketName := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get bucket details for confirmation
		bucket, err := client.ObjectStorage().GetBucket(cmd.Context(), bucketName)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("bucket %s not found", bucketName)
			}
			return fmt.Errorf("failed to get bucket: %w", err)
		}

		// Ask for confirmation unless --force is provided
		if !deleteForce {
			fmt.Printf("Are you sure you want to delete the following bucket?\n")
			fmt.Printf("  Name: %s\n", bucket.Name)
			fmt.Printf("  Status: %s\n", bucket.Status)
			fmt.Printf("  Total Size: %.2f GB\n", bucket.Usage.TotalSizeGB)
			fmt.Printf("  Total Objects: %d\n", bucket.Usage.TotalObjects)
			fmt.Printf("\nWARNING: This will permanently delete the bucket and all its contents!\n")
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		err = client.ObjectStorage().DeleteBucket(cmd.Context(), bucketName)
		if err != nil {
			return fmt.Errorf("failed to delete bucket: %w", err)
		}

		if deleteWait {
			fmt.Println("Waiting for bucket to be deleted...")

			// Parse timeout duration (default: 10 minutes)
			timeoutDuration := 10 * time.Minute
			if deleteTimeout != "" {
				parsedTimeout, err := time.ParseDuration(deleteTimeout)
				if err != nil {
					return fmt.Errorf("invalid timeout duration: %w", err)
				}
				timeoutDuration = parsedTimeout
			}

			timeout := time.After(timeoutDuration)
			tick := time.Tick(2 * time.Second)
			for {
				_, err := client.ObjectStorage().GetBucket(cmd.Context(), bucketName)
				if err != nil {
					if tcclient.IsNotFound(err) {
						// Bucket is deleted
						break
					}
					return fmt.Errorf("failed to check bucket status: %w", err)
				}
				select {
				case <-timeout:
					return fmt.Errorf("timeout waiting for bucket %s to be deleted", bucketName)
				case <-tick:
					// continue looping
				}
			}
			fmt.Println("Bucket deleted successfully")
		} else {
			fmt.Printf("Bucket %s deleted successfully\n", bucketName)
		}
		return nil
	},
}

func init() {
	ObjectStorageCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().BoolVarP(&deleteWait, "wait", "w", false, "Wait for the bucket to be deleted")
	deleteCmd.Flags().StringVar(&deleteTimeout, "timeout", "", "Timeout for waiting (e.g., 5m, 10m, 1h, default: 10m)")
}
