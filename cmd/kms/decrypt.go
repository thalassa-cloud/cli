package kms

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	decryptRegion     string
	decryptKey        string
	decryptCiphertext string
	decryptFromFile   string
	decryptToFile     string
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt ciphertext with a KMS key",
	Long: `Decrypt data with a KMS key.

Provide exactly one of --ciphertext or --from-file (file containing ciphertext).
When --to-file is set, decoded plaintext bytes are written with mode 0600.
Otherwise base64-encoded plaintext is printed to stdout.`,
	Example: `  tcloud kms decrypt --region nl-ams --key kms-123 --ciphertext 'thalassa:v1:...'
  tcloud kms decrypt --region nl-ams --key kms-123 --from-file secret.enc --to-file secret.txt`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ciphertext, err := resolveTextFileInput(decryptCiphertext, decryptFromFile, "ciphertext")
		if err != nil {
			return err
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		plaintext, err := client.KMS().DecryptBytes(cmd.Context(), decryptRegion, decryptKey, ciphertext)
		if err != nil {
			return fmt.Errorf("failed to decrypt: %w", err)
		}

		return writeBinaryResult(decryptToFile, plaintext)
	},
}

func init() {
	KmsCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVar(&decryptRegion, "region", "", "Region")
	decryptCmd.Flags().StringVar(&decryptKey, "key", "", "KMS key identity")
	decryptCmd.Flags().StringVar(&decryptCiphertext, "ciphertext", "", "Ciphertext from encrypt")
	decryptCmd.Flags().StringVar(&decryptFromFile, "from-file", "", "Read ciphertext from a file")
	decryptCmd.Flags().StringVar(&decryptToFile, "to-file", "", "Write decoded plaintext to a file (mode 0600) instead of stdout")
	_ = decryptCmd.MarkFlagRequired("region")
	_ = decryptCmd.MarkFlagRequired("key")
	decryptCmd.MarkFlagsMutuallyExclusive("ciphertext", "from-file")
	_ = decryptCmd.MarkFlagFilename("from-file")
	_ = decryptCmd.MarkFlagFilename("to-file")
	_ = decryptCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = decryptCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
}
