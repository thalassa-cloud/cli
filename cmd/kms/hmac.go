package kms

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

var (
	hmacRegion     string
	hmacKey        string
	hmacInput      string
	hmacKeyVersion string
	hmacAlgorithm  string
)

var hmacCmd = &cobra.Command{
	Use:   "hmac",
	Short: "Compute an HMAC with a KMS key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		result, err := client.KMS().HMAC(cmd.Context(), hmacRegion, hmacKey, clientkms.HMACRequest{
			Input:      hmacInput,
			KeyVersion: hmacKeyVersion,
			Algorithm:  hmacAlgorithm,
		})
		if err != nil {
			return fmt.Errorf("failed to compute HMAC: %w", err)
		}

		fmt.Println(result.HMAC)
		return nil
	},
}

func init() {
	KmsCmd.AddCommand(hmacCmd)
	hmacCmd.Flags().StringVar(&hmacRegion, "region", "", "Region")
	hmacCmd.Flags().StringVar(&hmacKey, "key", "", "KMS key identity")
	hmacCmd.Flags().StringVar(&hmacInput, "input", "", "Input for HMAC")
	hmacCmd.Flags().StringVar(&hmacKeyVersion, "key-version", "", "Key version")
	hmacCmd.Flags().StringVar(&hmacAlgorithm, "algorithm", "", "HMAC algorithm")
	_ = hmacCmd.MarkFlagRequired("region")
	_ = hmacCmd.MarkFlagRequired("key")
	_ = hmacCmd.MarkFlagRequired("input")
	_ = hmacCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = hmacCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
	_ = hmacCmd.RegisterFlagCompletionFunc("key-version", completion.CompleteKmsKeyVersion)
}
