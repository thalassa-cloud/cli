package kms

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

var (
	decryptRegion     string
	decryptKey        string
	decryptCiphertext string
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt ciphertext with a KMS key (prints base64 plaintext to stdout)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		result, err := client.KMS().Decrypt(cmd.Context(), decryptRegion, decryptKey, clientkms.DecryptRequest{
			Ciphertext: decryptCiphertext,
		})
		if err != nil {
			return fmt.Errorf("failed to decrypt: %w", err)
		}

		fmt.Println(result.Plaintext)
		return nil
	},
}

func init() {
	KmsCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVar(&decryptRegion, "region", "", "Region")
	decryptCmd.Flags().StringVar(&decryptKey, "key", "", "KMS key identity")
	decryptCmd.Flags().StringVar(&decryptCiphertext, "ciphertext", "", "Ciphertext from encrypt")
	_ = decryptCmd.MarkFlagRequired("region")
	_ = decryptCmd.MarkFlagRequired("key")
	_ = decryptCmd.MarkFlagRequired("ciphertext")
	_ = decryptCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = decryptCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
}
