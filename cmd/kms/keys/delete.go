package keys

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <key>",
	Aliases:           []string{"rm", "remove"},
	Short:             "Schedule a KMS key for deletion",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKmsKeyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		key, err := client.KMS().GetKey(cmd.Context(), region, args[0])
		if err != nil {
			return fmt.Errorf("failed to get KMS key: %w", err)
		}

		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf(
			"Are you sure you want to schedule this KMS key for deletion?\n  ID: %s\n  Name: %s\n",
			key.Identity, key.Name,
		))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		if err := client.KMS().DeleteKey(cmd.Context(), region, key.Identity); err != nil {
			return fmt.Errorf("failed to delete KMS key: %w", err)
		}
		fmt.Printf("Scheduled KMS key %s for deletion\n", key.Name)
		return nil
	},
}

func init() {
	KeysCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
	deleteCmd.Flags().StringVar(&region, "region", "", "Region")
	_ = deleteCmd.MarkFlagRequired("region")
	_ = deleteCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
