package keys

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var cancelDeletionCmd = &cobra.Command{
	Use:               "cancel-deletion <key>",
	Short:             "Cancel a pending KMS key deletion",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKmsKeyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if err := client.KMS().CancelDeletion(cmd.Context(), region, args[0]); err != nil {
			return fmt.Errorf("failed to cancel KMS key deletion: %w", err)
		}
		fmt.Printf("Cancelled deletion for KMS key %s\n", args[0])
		return nil
	},
}

func init() {
	KeysCmd.AddCommand(cancelDeletionCmd)
	cancelDeletionCmd.Flags().StringVar(&region, "region", "", "Region")
	_ = cancelDeletionCmd.MarkFlagRequired("region")
	_ = cancelDeletionCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
