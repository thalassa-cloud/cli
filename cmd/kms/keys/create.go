package keys

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

var (
	createName                 string
	createDescription          string
	createLabels               []string
	createAnnotations          []string
	createKeyType              string
	createExportAllowed        bool
	createKeyRotationEnabled   bool
	createRotationPeriodInDays int
	createImportKeyMaterial    string
	createHashFunction         string
	createAllowRotation        bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a KMS key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		keyType := clientkms.KmsKeyType(createKeyType)
		if createKeyType != "" && !keyType.IsValid() {
			return fmt.Errorf("unsupported key type %q", createKeyType)
		}

		req := clientkms.CreateKmsKeyRequest{
			Name:               createName,
			Description:        createDescription,
			Labels:             shared.KeyValuePairsToMap(createLabels),
			Annotations:        shared.KeyValuePairsToMap(createAnnotations),
			KeyType:            keyType,
			ExportAllowed:      createExportAllowed,
			KeyRotationEnabled: createKeyRotationEnabled,
			ImportKeyMaterial:  createImportKeyMaterial,
			HashFunction:       createHashFunction,
			AllowRotation:      createAllowRotation,
		}
		if cmd.Flags().Changed("rotation-period-days") {
			period := createRotationPeriodInDays
			req.RotationPeriodInDays = &period
		}

		key, err := client.KMS().CreateKey(cmd.Context(), region, req)
		if err != nil {
			return fmt.Errorf("failed to create KMS key: %w", err)
		}

		body := [][]string{{
			key.Identity,
			key.Name,
			string(key.KeyType),
			string(key.Status),
			formattime.FormatTime(key.CreatedAt.Local(), showExactTime),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Type", "Status", "Created"}, body)
		}
		return nil
	},
}

func init() {
	KeysCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	createCmd.Flags().StringVar(&region, "region", "", "Region")
	createCmd.Flags().StringVar(&createName, "name", "", "Key name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Key description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	createCmd.Flags().StringVar(&createKeyType, "key-type", "", "Key type (aes128-gcm96, aes256-gcm96, chacha20-poly1305, ed25519, ecdsa-p256/384/521, rsa-2048/3072/4096, hmac, hmac-sha256, hmac-sha512)")
	createCmd.Flags().BoolVar(&createExportAllowed, "export-allowed", false, "Allow exporting key material")
	createCmd.Flags().BoolVar(&createKeyRotationEnabled, "rotation-enabled", false, "Enable automatic key rotation")
	createCmd.Flags().IntVar(&createRotationPeriodInDays, "rotation-period-days", 0, "Automatic rotation period in days")
	createCmd.Flags().StringVar(&createImportKeyMaterial, "import-key-material", "", "Wrapped key material for BYOK import")
	createCmd.Flags().StringVar(&createHashFunction, "hash-function", "", "Hash function for imported or HMAC keys")
	createCmd.Flags().BoolVar(&createAllowRotation, "allow-rotation", false, "Allow rotation for imported keys")
	_ = createCmd.MarkFlagRequired("region")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = createCmd.RegisterFlagCompletionFunc("key-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			"aes128-gcm96", "aes256-gcm96", "chacha20-poly1305",
			"ed25519", "ecdsa-p256", "ecdsa-p384", "ecdsa-p521",
			"rsa-2048", "rsa-3072", "rsa-4096",
			"hmac", "hmac-sha256", "hmac-sha512",
		}, cobra.ShellCompDirectiveNoFileComp
	})
}
