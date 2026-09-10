package secrets

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	deleteRegion string
	deletePath   string
	deleteForce  bool
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del"},
	Short:   "Delete a secret",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		proceed, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Delete secret %q in region %s?\n", deletePath, deleteRegion))
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if err := client.Secrets().DeleteSecret(cmd.Context(), deleteRegion, deletePath); err != nil {
			return fmt.Errorf("failed to delete secret: %w", err)
		}
		fmt.Printf("Secret %s deleted\n", deletePath)
		return nil
	},
}

func init() {
	SecretsCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteRegion, "region", "", "Region")
	deleteCmd.Flags().StringVar(&deletePath, "path", "", "Secret path")
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip confirmation")
	_ = deleteCmd.MarkFlagRequired("region")
	_ = deleteCmd.MarkFlagRequired("path")
	_ = deleteCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
