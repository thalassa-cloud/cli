package kms

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

var (
	verifyRegion        string
	verifyKey           string
	verifyInput         string
	verifySignature     string
	verifyHashAlgorithm string

	verifyHmacRegion        string
	verifyHmacKey           string
	verifyHmacInput         string
	verifyHmacValue         string
	verifyHmacHashAlgorithm string
)

func printValidity(valid bool) {
	fmt.Printf("%v\n", valid)
}

func runVerifySignature(ctx context.Context) error {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	result, err := client.KMS().VerifySignature(ctx, verifyRegion, verifyKey, clientkms.VerifySignatureRequest{
		Input:         verifyInput,
		Signature:     verifySignature,
		HashAlgorithm: verifyHashAlgorithm,
	})
	if err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}
	printValidity(result.Valid)
	return nil
}

func runVerifyHMAC(ctx context.Context) error {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	result, err := client.KMS().VerifyHMAC(ctx, verifyHmacRegion, verifyHmacKey, clientkms.VerifyHMACRequest{
		Input:         verifyHmacInput,
		HMAC:          verifyHmacValue,
		HashAlgorithm: verifyHmacHashAlgorithm,
	})
	if err != nil {
		return fmt.Errorf("failed to verify HMAC: %w", err)
	}
	printValidity(result.Valid)
	return nil
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a signature with an asymmetric KMS key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return runVerifySignature(cmd.Context())
	},
}

var verifyHmacCmd = &cobra.Command{
	Use:   "verify-hmac",
	Short: "Verify an HMAC with a KMS key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return runVerifyHMAC(cmd.Context())
	},
}

func init() {
	KmsCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVar(&verifyRegion, "region", "", "Region")
	verifyCmd.Flags().StringVar(&verifyKey, "key", "", "KMS key identity")
	verifyCmd.Flags().StringVar(&verifyInput, "input", "", "Input that was signed")
	verifyCmd.Flags().StringVar(&verifySignature, "signature", "", "Signature to verify")
	verifyCmd.Flags().StringVar(&verifyHashAlgorithm, "hash-algorithm", "", "Hash algorithm")
	_ = verifyCmd.MarkFlagRequired("region")
	_ = verifyCmd.MarkFlagRequired("key")
	_ = verifyCmd.MarkFlagRequired("input")
	_ = verifyCmd.MarkFlagRequired("signature")
	_ = verifyCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = verifyCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)

	KmsCmd.AddCommand(verifyHmacCmd)
	verifyHmacCmd.Flags().StringVar(&verifyHmacRegion, "region", "", "Region")
	verifyHmacCmd.Flags().StringVar(&verifyHmacKey, "key", "", "KMS key identity")
	verifyHmacCmd.Flags().StringVar(&verifyHmacInput, "input", "", "Input that was HMACed")
	verifyHmacCmd.Flags().StringVar(&verifyHmacValue, "hmac", "", "HMAC value to verify")
	verifyHmacCmd.Flags().StringVar(&verifyHmacHashAlgorithm, "hash-algorithm", "", "Hash algorithm")
	_ = verifyHmacCmd.MarkFlagRequired("region")
	_ = verifyHmacCmd.MarkFlagRequired("key")
	_ = verifyHmacCmd.MarkFlagRequired("input")
	_ = verifyHmacCmd.MarkFlagRequired("hmac")
	_ = verifyHmacCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = verifyHmacCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
}
