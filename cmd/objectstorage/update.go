package objectstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/objectstorage"
)

const (
	UpdateFlagPublic            = "public"
	UpdateFlagPolicy            = "policy"
	UpdateFlagVersioning        = "versioning"
	UpdateFlagObjectLockEnabled = "object-lock-enabled"
	UpdateFlagLabels            = "labels"
	UpdateFlagAnnotations       = "annotations"
)

var (
	updatePublic            *bool
	updatePolicy            string
	updateVersioning        *bool
	updateObjectLockEnabled *bool
	updateLabels            []string
	updateAnnotations       []string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update an object storage bucket",
	Aliases: []string{"update-bucket"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketName := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get current bucket to preserve values if not provided
		current, err := client.ObjectStorage().GetBucket(cmd.Context(), bucketName)
		if err != nil {
			return fmt.Errorf("failed to get bucket: %w", err)
		}

		req := objectstorage.UpdateBucketRequest{
			Public:            current.Public,
			PolicyDocument:    &current.Policy,
			Versioning:        current.Versioning,
			ObjectLockEnabled: &current.ObjectLockEnabled,
			Labels:            current.Labels,
			Annotations:       current.Annotations,
		}

		// Update only provided fields
		if cmd.Flags().Changed(UpdateFlagPublic) {
			req.Public = *updatePublic
		}
		if cmd.Flags().Changed(UpdateFlagPolicy) {
			policyJSON := updatePolicy

			// Check if it's a file path (contains / or starts with . or is absolute)
			if strings.Contains(updatePolicy, "/") || strings.HasPrefix(updatePolicy, ".") || strings.HasPrefix(updatePolicy, string(os.PathSeparator)) {
				fileData, err := os.ReadFile(updatePolicy)
				if err != nil {
					return fmt.Errorf("failed to read policy file: %w", err)
				}
				policyJSON = string(fileData)
			}

			policyDoc := &objectstorage.PolicyDocument{}
			if err := json.Unmarshal([]byte(policyJSON), policyDoc); err != nil {
				return fmt.Errorf("failed to parse policy JSON: %w", err)
			}
			req.PolicyDocument = policyDoc
		}
		if cmd.Flags().Changed(UpdateFlagVersioning) {
			if *updateVersioning {
				req.Versioning = objectstorage.ObjectStorageBucketVersioningEnabled
			} else {
				req.Versioning = objectstorage.ObjectStorageBucketVersioningDisabled
			}
		}
		if cmd.Flags().Changed(UpdateFlagObjectLockEnabled) {
			req.ObjectLockEnabled = updateObjectLockEnabled
		}

		// Parse labels from key=value format
		if cmd.Flags().Changed(UpdateFlagLabels) {
			labels := make(map[string]string)
			for _, label := range updateLabels {
				parts := strings.SplitN(label, "=", 2)
				if len(parts) == 2 {
					labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Labels = labels
		}

		// Parse annotations from key=value format
		if cmd.Flags().Changed(UpdateFlagAnnotations) {
			annotations := make(map[string]string)
			for _, annotation := range updateAnnotations {
				parts := strings.SplitN(annotation, "=", 2)
				if len(parts) == 2 {
					annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Annotations = annotations
		}

		bucket, err := client.ObjectStorage().UpdateBucket(cmd.Context(), bucketName, req)
		if err != nil {
			return fmt.Errorf("failed to update bucket: %w", err)
		}

		body := make([][]string, 0, 1)
		body = append(body, []string{
			bucket.Name,
			bucket.Status,
			fmt.Sprintf("%v", bucket.Public),
			formattime.FormatTime(bucket.UpdatedAt.Local(), false),
		})
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Name", "Status", "Public", "Updated"}, body)
		}
		return nil
	},
}

func init() {
	ObjectStorageCmd.AddCommand(updateCmd)

	// Initialize pointer variables
	updatePublic = new(bool)
	updateVersioning = new(bool)
	updateObjectLockEnabled = new(bool)

	updateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	updateCmd.Flags().BoolVar(updatePublic, UpdateFlagPublic, false, "Make the bucket publicly accessible")
	updateCmd.Flags().StringVar(&updatePolicy, UpdateFlagPolicy, "", "Bucket policy as JSON string or path to a JSON file")
	updateCmd.Flags().BoolVar(updateVersioning, UpdateFlagVersioning, false, "Enable versioning for the bucket")
	updateCmd.Flags().BoolVar(updateObjectLockEnabled, UpdateFlagObjectLockEnabled, false, "Enable object lock")
	updateCmd.Flags().StringSliceVar(&updateLabels, UpdateFlagLabels, []string{}, "Labels in key=value format")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, UpdateFlagAnnotations, []string{}, "Annotations in key=value format")
}
