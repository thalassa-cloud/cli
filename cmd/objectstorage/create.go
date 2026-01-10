package objectstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/objectstorage"
)

const (
	CreateFlagName              = "name"
	CreateFlagRegion            = "region"
	CreateFlagPublic            = "public"
	CreateFlagPolicy            = "policy"
	CreateFlagVersioning        = "versioning"
	CreateFlagObjectLockEnabled = "object-lock-enabled"
	CreateFlagLabels            = "labels"
	CreateFlagAnnotations       = "annotations"
)

var (
	createName              string
	createRegion            string
	createPublic            bool
	createPolicy            string
	createVersioning        bool
	createObjectLockEnabled bool
	createLabels            []string
	createAnnotations       []string
	createWait              bool
	createTimeout           string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create an object storage bucket",
	Aliases: []string{"create-bucket"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if createName == "" {
			return fmt.Errorf("name is required")
		}
		if createRegion == "" {
			return fmt.Errorf("region is required")
		}

		// Parse policy if provided
		var policyDoc *objectstorage.PolicyDocument
		if createPolicy != "" {
			policyJSON := createPolicy

			// Check if it's a file path (contains / or starts with . or is absolute)
			if strings.Contains(createPolicy, "/") || strings.HasPrefix(createPolicy, ".") || strings.HasPrefix(createPolicy, string(os.PathSeparator)) {
				fileData, err := os.ReadFile(createPolicy)
				if err != nil {
					return fmt.Errorf("failed to read policy file: %w", err)
				}
				policyJSON = string(fileData)
			}

			policyDoc = &objectstorage.PolicyDocument{}
			if err := json.Unmarshal([]byte(policyJSON), policyDoc); err != nil {
				return fmt.Errorf("failed to parse policy JSON: %w", err)
			}
		}

		// Parse versioning - boolean flag means Enabled if true, Disabled if false
		versioning := objectstorage.ObjectStorageBucketVersioningDisabled
		if createVersioning {
			versioning = objectstorage.ObjectStorageBucketVersioningEnabled
		}

		// Parse labels from key=value format
		labels := make(map[string]string)
		for _, label := range createLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Parse annotations from key=value format
		annotations := make(map[string]string)
		for _, annotation := range createAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		req := objectstorage.CreateBucketRequest{
			BucketName:        createName,
			Region:            createRegion,
			Public:            createPublic,
			PolicyDocument:    policyDoc,
			Versioning:        versioning,
			ObjectLockEnabled: createObjectLockEnabled,
			Labels:            labels,
			Annotations:       annotations,
		}

		bucket, err := client.ObjectStorage().CreateBucket(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		if createWait {
			fmt.Println("Waiting for bucket to be ready...")

			// Parse timeout duration (default: 10 minutes)
			timeoutDuration := 10 * time.Minute
			if createTimeout != "" {
				parsedTimeout, err := time.ParseDuration(createTimeout)
				if err != nil {
					return fmt.Errorf("invalid timeout duration: %w", err)
				}
				timeoutDuration = parsedTimeout
			}

			timeout := time.After(timeoutDuration)
			tick := time.Tick(2 * time.Second)
			for {
				bucket, err = client.ObjectStorage().GetBucket(cmd.Context(), bucket.Name)
				if err != nil {
					return fmt.Errorf("failed to get bucket: %w", err)
				}
				// Check if bucket is ready (status is typically "ready" or "available" when ready)
				if strings.EqualFold(bucket.Status, "ready") || strings.EqualFold(bucket.Status, "available") || strings.EqualFold(bucket.Status, "active") {
					break
				}
				select {
				case <-timeout:
					return fmt.Errorf("timeout waiting for bucket %s to be ready (current status: %s)", bucket.Name, bucket.Status)
				case <-tick:
					// continue looping
				}
			}
			fmt.Println("Bucket is ready")
		}

		body := make([][]string, 0, 1)
		body = append(body, []string{
			bucket.Name,
			bucket.Status,
			fmt.Sprintf("%v", bucket.Public),
			formattime.FormatTime(bucket.CreatedAt.Local(), false),
		})
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Name", "Status", "Public", "Age"}, body)
		}
		return nil
	},
}

func init() {
	ObjectStorageCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createName, CreateFlagName, "", "Name of the bucket")
	createCmd.Flags().StringVar(&createRegion, CreateFlagRegion, "", "Region for the bucket")
	createCmd.Flags().StringVar(&createPolicy, CreateFlagPolicy, "", "Bucket policy as JSON string or path to a JSON file")
	createCmd.Flags().BoolVar(&createVersioning, CreateFlagVersioning, false, "Enable versioning for the bucket")
	createCmd.Flags().BoolVar(&createObjectLockEnabled, CreateFlagObjectLockEnabled, false, "Enable object lock")
	createCmd.Flags().StringSliceVar(&createLabels, CreateFlagLabels, []string{}, "Labels in key=value format")
	createCmd.Flags().StringSliceVar(&createAnnotations, CreateFlagAnnotations, []string{}, "Annotations in key=value format")
	createCmd.Flags().BoolVarP(&createWait, "wait", "w", false, "Wait for the bucket to be ready")
	createCmd.Flags().StringVar(&createTimeout, "timeout", "", "Timeout for waiting (e.g., 5m, 10m, 1h, default: 10m)")

	createCmd.MarkFlagRequired(CreateFlagName)
	createCmd.MarkFlagRequired(CreateFlagRegion)
}
