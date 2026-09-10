package kms

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

var (
	signRegion        string
	signKey           string
	signInput         string
	signFromFile      string
	signToFile        string
	signKeyVersion    string
	signHashAlgorithm string
	signPrehashed     bool
	signContext       string
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign input with an asymmetric KMS key",
	Long: `Sign data with an asymmetric KMS key.

Provide exactly one of --input (base64-encoded) or --from-file (raw file bytes).
The signature is written to --to-file when set, otherwise to stdout.`,
	Example: `  tcloud kms sign --region nl-ams --key kms-123 --input "$(echo -n hello | base64)"
  tcloud kms sign --region nl-ams --key kms-123 --from-file message.txt --to-file message.sig`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		input, err := resolveWireInput(signInput, signFromFile, "input")
		if err != nil {
			return err
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		result, err := client.KMS().Sign(cmd.Context(), signRegion, signKey, clientkms.SignRequest{
			Input:         input,
			KeyVersion:    signKeyVersion,
			HashAlgorithm: signHashAlgorithm,
			Prehashed:     signPrehashed,
			Context:       signContext,
		})
		if err != nil {
			return fmt.Errorf("failed to sign: %w", err)
		}

		return writeTextResult(signToFile, result.Signature)
	},
}

func init() {
	KmsCmd.AddCommand(signCmd)
	signCmd.Flags().StringVar(&signRegion, "region", "", "Region")
	signCmd.Flags().StringVar(&signKey, "key", "", "KMS key identity")
	signCmd.Flags().StringVar(&signInput, "input", "", "Base64-encoded input to sign")
	signCmd.Flags().StringVar(&signFromFile, "from-file", "", "Read raw input bytes from a file")
	signCmd.Flags().StringVar(&signToFile, "to-file", "", "Write signature to a file (mode 0600) instead of stdout")
	signCmd.Flags().StringVar(&signKeyVersion, "key-version", "", "Key version")
	signCmd.Flags().StringVar(&signHashAlgorithm, "hash-algorithm", "", "Hash algorithm")
	signCmd.Flags().BoolVar(&signPrehashed, "prehashed", false, "Input is already hashed")
	signCmd.Flags().StringVar(&signContext, "context", "", "Optional signing context")
	_ = signCmd.MarkFlagRequired("region")
	_ = signCmd.MarkFlagRequired("key")
	signCmd.MarkFlagsMutuallyExclusive("input", "from-file")
	_ = signCmd.MarkFlagFilename("from-file")
	_ = signCmd.MarkFlagFilename("to-file")
	_ = signCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = signCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
}
