package kms

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

var (
	exportRegion     string
	exportKey        string
	exportKeyVersion string
	exportForce      bool
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export key material for a KMS key (when export is allowed)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ok, err := shared.PromptDestructiveUnlessForce(exportForce, fmt.Sprintf(
			"Exporting key material for %s. Ensure the output is handled securely.\n", exportKey,
		))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		result, err := client.KMS().ExportKey(cmd.Context(), exportRegion, exportKey, clientkms.ExportKeyRequest{
			KeyVersion: exportKeyVersion,
		})
		if err != nil {
			return fmt.Errorf("failed to export KMS key: %w", err)
		}

		fmt.Print(result.KeyMaterial)
		if result.KeyMaterial != "" && result.KeyMaterial[len(result.KeyMaterial)-1] != '\n' {
			fmt.Println()
		}
		return nil
	},
}

func init() {
	KmsCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&exportRegion, "region", "", "Region")
	exportCmd.Flags().StringVar(&exportKey, "key", "", "KMS key identity")
	exportCmd.Flags().StringVar(&exportKeyVersion, "key-version", "", "Key version to export")
	exportCmd.Flags().BoolVar(&exportForce, shared.ForceKey, false, "Skip the confirmation prompt")
	_ = exportCmd.MarkFlagRequired("region")
	_ = exportCmd.MarkFlagRequired("key")
	_ = exportCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = exportCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
	_ = exportCmd.RegisterFlagCompletionFunc("key-version", completion.CompleteKmsKeyVersion)
}
