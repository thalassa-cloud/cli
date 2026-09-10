package kms

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	publicKeyRegion  string
	publicKeyKey     string
	publicKeyVersion int
)

var publicKeyCmd = &cobra.Command{
	Use:   "public-key",
	Short: "Get the public key material for an asymmetric KMS key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		var version *int
		if cmd.Flags().Changed("version") {
			v := publicKeyVersion
			version = &v
		}

		result, err := client.KMS().GetPublicKey(cmd.Context(), publicKeyRegion, publicKeyKey, version)
		if err != nil {
			return fmt.Errorf("failed to get public key: %w", err)
		}

		keys := make([]string, 0, len(result.Keys))
		for k := range result.Keys {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		body := make([][]string, 0, len(keys))
		for _, k := range keys {
			body = append(body, []string{k, result.Keys[k]})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Version", "Public Key"}, body)
		}
		return nil
	},
}

func init() {
	KmsCmd.AddCommand(publicKeyCmd)
	publicKeyCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	publicKeyCmd.Flags().StringVar(&publicKeyRegion, "region", "", "Region")
	publicKeyCmd.Flags().StringVar(&publicKeyKey, "key", "", "KMS key identity")
	publicKeyCmd.Flags().IntVar(&publicKeyVersion, "version", 0, "Key version")
	_ = publicKeyCmd.MarkFlagRequired("region")
	_ = publicKeyCmd.MarkFlagRequired("key")
	_ = publicKeyCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = publicKeyCmd.RegisterFlagCompletionFunc("key", completion.CompleteKmsKeyIdentity)
	_ = publicKeyCmd.RegisterFlagCompletionFunc("version", completion.CompleteKmsKeyVersion)
}
