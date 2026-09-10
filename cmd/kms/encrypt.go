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
	encryptKeyVersion string
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt plaintext with a KMS key (plaintext must be base64-encoded)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		result, err := client.KMS().Encrypt(cmd.Context(), encryptRegion, encryptKey, clientkms.EncryptRequest{
			Plaintext:  encryptPlaintext,
			KeyVersion: encryptKeyVersion,
		})
		if err != nil {
			return fmt.Errorf("failed to encrypt: %w", err)
		}

		fmt.Println(result.Ciphertext)
		return nil
	},
}

func init() {
	KmsCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVar(&encryptRegion, "region", "", "Region")
	encryptCmd.Flags().StringVar(&encryptKey, "key", "", "KMS key identity")
	encryptCmd.Flags().StringVar(&encryptPlaintext, "plaintext", "", "Base64-encoded plaintext")
	encryptCmd.Flags().StringVar(&encryptKeyVersion, "key-version", "", "Key version")
	_ = encryptCmd.MarkFlagRequired("region")
	_ = encryptCmd.MarkFlagRequired("key")
	_ = encryptCmd.MarkFlagRequired("plaintext")
	_ = encryptCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = encryptCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
}
