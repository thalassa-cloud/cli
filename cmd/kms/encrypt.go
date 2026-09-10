package kms

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

var (
	encryptRegion     string
	encryptKey        string
	encryptPlaintext  string
	encryptFromFile   string
	encryptToFile     string
	encryptKeyVersion string
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt plaintext with a KMS key",
	Long: `Encrypt data with a KMS key.

Provide exactly one of --plaintext (base64-encoded) or --from-file (raw file
bytes). Ciphertext is written to --to-file when set, otherwise to stdout.`,
	Example: `  tcloud kms encrypt --region nl-ams --key kms-123 --plaintext "$(echo -n hello | base64)"
  tcloud kms encrypt --region nl-ams --key kms-123 --from-file secret.txt --to-file secret.enc`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		plaintext, err := resolveWireInput(encryptPlaintext, encryptFromFile, "plaintext")
		if err != nil {
			return err
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		result, err := client.KMS().Encrypt(cmd.Context(), encryptRegion, encryptKey, clientkms.EncryptRequest{
			Plaintext:  plaintext,
			KeyVersion: encryptKeyVersion,
		})
		if err != nil {
			return fmt.Errorf("failed to encrypt: %w", err)
		}

		return writeTextResult(encryptToFile, result.Ciphertext)
	},
}

func init() {
	KmsCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVar(&encryptRegion, "region", "", "Region")
	encryptCmd.Flags().StringVar(&encryptKey, "key", "", "KMS key identity")
	encryptCmd.Flags().StringVar(&encryptPlaintext, "plaintext", "", "Base64-encoded plaintext")
	encryptCmd.Flags().StringVar(&encryptFromFile, "from-file", "", "Read raw plaintext bytes from a file")
	encryptCmd.Flags().StringVar(&encryptToFile, "to-file", "", "Write ciphertext to a file (mode 0600) instead of stdout")
	encryptCmd.Flags().StringVar(&encryptKeyVersion, "key-version", "", "Key version")
	_ = encryptCmd.MarkFlagRequired("region")
	_ = encryptCmd.MarkFlagRequired("key")
	encryptCmd.MarkFlagsMutuallyExclusive("plaintext", "from-file")
	_ = encryptCmd.MarkFlagFilename("from-file")
	_ = encryptCmd.MarkFlagFilename("to-file")
	_ = encryptCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = encryptCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
}
