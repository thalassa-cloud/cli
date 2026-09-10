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
	signKeyVersion    string
	signHashAlgorithm string
	signPrehashed     bool
	signContext       string
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign input with an asymmetric KMS key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		result, err := client.KMS().Sign(cmd.Context(), signRegion, signKey, clientkms.SignRequest{
			Input:         signInput,
			KeyVersion:    signKeyVersion,
			HashAlgorithm: signHashAlgorithm,
			Prehashed:     signPrehashed,
			Context:       signContext,
		})
		if err != nil {
			return fmt.Errorf("failed to sign: %w", err)
		}

		fmt.Println(result.Signature)
		return nil
	},
}

func init() {
	KmsCmd.AddCommand(signCmd)
	signCmd.Flags().StringVar(&signRegion, "region", "", "Region")
	signCmd.Flags().StringVar(&signKey, "key", "", "KMS key identity")
	signCmd.Flags().StringVar(&signInput, "input", "", "Input to sign (as required by the API)")
	signCmd.Flags().StringVar(&signKeyVersion, "key-version", "", "Key version")
	signCmd.Flags().StringVar(&signHashAlgorithm, "hash-algorithm", "", "Hash algorithm")
	signCmd.Flags().BoolVar(&signPrehashed, "prehashed", false, "Input is already hashed")
	signCmd.Flags().StringVar(&signContext, "context", "", "Optional signing context")
	_ = signCmd.MarkFlagRequired("region")
	_ = signCmd.MarkFlagRequired("key")
	_ = signCmd.MarkFlagRequired("input")
	_ = signCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = signCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
}
