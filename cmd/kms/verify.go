package kms

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

var (
	verifyRegion        string
	verifyKey           string
	verifyInput         string
	verifyFromFile      string
	verifySignature     string
	verifySignatureFile string
	verifyToFile        string
	verifyHashAlgorithm string

	verifyHmacRegion        string
	verifyHmacKey           string
	verifyHmacInput         string
	verifyHmacFromFile      string
	verifyHmacValue         string
	verifyHmacFromHMACFile  string
	verifyHmacToFile        string
	verifyHmacHashAlgorithm string
)

func resolveSignatureInput(flagValue, fromFile, flagName string) (string, error) {
	if flagValue == "" && fromFile == "" {
		return "", fmt.Errorf("provide --%s or --%s-file", flagName, flagName)
	}
	if fromFile == "" {
		return flagValue, nil
	}
	data, err := os.ReadFile(fromFile)
	if err != nil {
		return "", fmt.Errorf("read --%s-file: %w", flagName, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func writeValidity(toFile string, valid bool) error {
	return writeTextResult(toFile, fmt.Sprintf("%v", valid))
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a signature with an asymmetric KMS key",
	Long: `Verify a signature with an asymmetric KMS key.

Provide exactly one of --input (base64-encoded) or --from-file (raw file bytes),
and exactly one of --signature or --signature-file. Validity is written to
--to-file when set, otherwise to stdout.`,
	Example: `  tcloud kms verify --region nl-ams --key kms-123 --input "$(echo -n hello | base64)" --signature '...'
  tcloud kms verify --region nl-ams --key kms-123 --from-file message.txt --signature-file message.sig`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		input, err := resolveWireInput(verifyInput, verifyFromFile, "input")
		if err != nil {
			return err
		}
		signature, err := resolveSignatureInput(verifySignature, verifySignatureFile, "signature")
		if err != nil {
			return err
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		result, err := client.KMS().VerifySignature(cmd.Context(), verifyRegion, verifyKey, clientkms.VerifySignatureRequest{
			Input:         input,
			Signature:     signature,
			HashAlgorithm: verifyHashAlgorithm,
		})
		if err != nil {
			return fmt.Errorf("failed to verify signature: %w", err)
		}

		return writeValidity(verifyToFile, result.Valid)
	},
}

var verifyHmacCmd = &cobra.Command{
	Use:   "verify-hmac",
	Short: "Verify an HMAC with a KMS key",
	Long: `Verify an HMAC with a KMS key.

Provide exactly one of --input (base64-encoded) or --from-file (raw file bytes),
and exactly one of --hmac or --hmac-file. Validity is written to --to-file when
set, otherwise to stdout.`,
	Example: `  tcloud kms verify-hmac --region nl-ams --key kms-123 --from-file message.txt --hmac-file message.hmac`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		input, err := resolveWireInput(verifyHmacInput, verifyHmacFromFile, "input")
		if err != nil {
			return err
		}
		hmacValue, err := resolveSignatureInput(verifyHmacValue, verifyHmacFromHMACFile, "hmac")
		if err != nil {
			return err
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		result, err := client.KMS().VerifyHMAC(cmd.Context(), verifyHmacRegion, verifyHmacKey, clientkms.VerifyHMACRequest{
			Input:         input,
			HMAC:          hmacValue,
			HashAlgorithm: verifyHmacHashAlgorithm,
		})
		if err != nil {
			return fmt.Errorf("failed to verify HMAC: %w", err)
		}

		return writeValidity(verifyHmacToFile, result.Valid)
	},
}

func init() {
	KmsCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVar(&verifyRegion, "region", "", "Region")
	verifyCmd.Flags().StringVar(&verifyKey, "key", "", "KMS key identity")
	verifyCmd.Flags().StringVar(&verifyInput, "input", "", "Base64-encoded input that was signed")
	verifyCmd.Flags().StringVar(&verifyFromFile, "from-file", "", "Read raw input bytes from a file")
	verifyCmd.Flags().StringVar(&verifySignature, "signature", "", "Signature to verify")
	verifyCmd.Flags().StringVar(&verifySignatureFile, "signature-file", "", "Read signature from a file")
	verifyCmd.Flags().StringVar(&verifyToFile, "to-file", "", "Write validity (true/false) to a file instead of stdout")
	verifyCmd.Flags().StringVar(&verifyHashAlgorithm, "hash-algorithm", "", "Hash algorithm")
	_ = verifyCmd.MarkFlagRequired("region")
	_ = verifyCmd.MarkFlagRequired("key")
	verifyCmd.MarkFlagsMutuallyExclusive("input", "from-file")
	verifyCmd.MarkFlagsMutuallyExclusive("signature", "signature-file")
	_ = verifyCmd.MarkFlagFilename("from-file")
	_ = verifyCmd.MarkFlagFilename("signature-file")
	_ = verifyCmd.MarkFlagFilename("to-file")
	_ = verifyCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = verifyCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)

	KmsCmd.AddCommand(verifyHmacCmd)
	verifyHmacCmd.Flags().StringVar(&verifyHmacRegion, "region", "", "Region")
	verifyHmacCmd.Flags().StringVar(&verifyHmacKey, "key", "", "KMS key identity")
	verifyHmacCmd.Flags().StringVar(&verifyHmacInput, "input", "", "Base64-encoded input that was HMACed")
	verifyHmacCmd.Flags().StringVar(&verifyHmacFromFile, "from-file", "", "Read raw input bytes from a file")
	verifyHmacCmd.Flags().StringVar(&verifyHmacValue, "hmac", "", "HMAC value to verify")
	verifyHmacCmd.Flags().StringVar(&verifyHmacFromHMACFile, "hmac-file", "", "Read HMAC value from a file")
	verifyHmacCmd.Flags().StringVar(&verifyHmacToFile, "to-file", "", "Write validity (true/false) to a file instead of stdout")
	verifyHmacCmd.Flags().StringVar(&verifyHmacHashAlgorithm, "hash-algorithm", "", "Hash algorithm")
	_ = verifyHmacCmd.MarkFlagRequired("region")
	_ = verifyHmacCmd.MarkFlagRequired("key")
	verifyHmacCmd.MarkFlagsMutuallyExclusive("input", "from-file")
	verifyHmacCmd.MarkFlagsMutuallyExclusive("hmac", "hmac-file")
	_ = verifyHmacCmd.MarkFlagFilename("from-file")
	_ = verifyHmacCmd.MarkFlagFilename("hmac-file")
	_ = verifyHmacCmd.MarkFlagFilename("to-file")
	_ = verifyHmacCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = verifyHmacCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
}
