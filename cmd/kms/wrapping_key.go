package kms

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var wrappingKeyRegion string

var wrappingKeyCmd = &cobra.Command{
	Use:   "wrapping-key",
	Short: "Get the regional wrapping public key for BYOK import",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		wrappingKey, err := client.KMS().GetWrappingKey(cmd.Context(), wrappingKeyRegion)
		if err != nil {
			return fmt.Errorf("failed to get wrapping key: %w", err)
		}

		body := [][]string{{wrappingKey.Algorithm, wrappingKey.PublicKey}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Algorithm", "Public Key"}, body)
		}
		return nil
	},
}

func init() {
	KmsCmd.AddCommand(wrappingKeyCmd)
	wrappingKeyCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	wrappingKeyCmd.Flags().StringVar(&wrappingKeyRegion, "region", "", "Region")
	_ = wrappingKeyCmd.MarkFlagRequired("region")
	_ = wrappingKeyCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
