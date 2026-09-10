package kms

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/kms/keys"
)

// KmsCmd manages Key Management Service resources and crypto operations.
var KmsCmd = &cobra.Command{
	Use:   "kms",
	Short: "Manage KMS keys and cryptographic operations (beta)",
	Long: `Manage Key Management Service keys and cryptographic operations.

Note: This command is in beta.`,
}

var (
	noHeader bool
)

func init() {
	KmsCmd.AddCommand(keys.KeysCmd)
}
