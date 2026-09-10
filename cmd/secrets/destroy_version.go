package secrets

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	destroyRegion  string
	destroyPath    string
	destroyVersion int
	destroyForce   bool
)

var destroyVersionCmd = &cobra.Command{
	Use:   "destroy-version",
	Short: "Permanently destroy a secret version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		proceed, err := shared.PromptDestructiveUnlessForce(destroyForce, fmt.Sprintf("Destroy version %d of secret %q?\n", destroyVersion, destroyPath))
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

		if err := client.Secrets().DestroySecretVersion(cmd.Context(), destroyRegion, destroyPath, destroyVersion); err != nil {
			return fmt.Errorf("failed to destroy secret version: %w", err)
		}
		fmt.Printf("Destroyed version %d of %s\n", destroyVersion, destroyPath)
		return nil
	},
}

func init() {
	SecretsCmd.AddCommand(destroyVersionCmd)
	destroyVersionCmd.Flags().StringVar(&destroyRegion, "region", "", "Region")
	destroyVersionCmd.Flags().StringVar(&destroyPath, "path", "", "Secret path")
	destroyVersionCmd.Flags().IntVar(&destroyVersion, "version", 0, "Version to destroy")
	destroyVersionCmd.Flags().BoolVar(&destroyForce, shared.ForceKey, false, "Skip confirmation")
	_ = destroyVersionCmd.MarkFlagRequired("region")
	_ = destroyVersionCmd.MarkFlagRequired("path")
	_ = destroyVersionCmd.MarkFlagRequired("version")
	_ = destroyVersionCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = destroyVersionCmd.RegisterFlagCompletionFunc("path", completion.CompleteSecretPath)
	_ = destroyVersionCmd.RegisterFlagCompletionFunc("version", completion.CompleteSecretVersion)
}
