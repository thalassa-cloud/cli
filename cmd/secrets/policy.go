package secrets

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientsecrets "github.com/thalassa-cloud/client-go/secrets"
)

var (
	policyRegion string
	policyPath   string
	policyFile   string
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Replace a secret access policy from a JSON file",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		policy, err := readSecretPolicy(policyFile)
		if err != nil {
			return err
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		secret, err := client.Secrets().UpdateAccessPolicy(cmd.Context(), policyRegion, policyPath, clientsecrets.UpdateAccessPolicyRequest{
			AccessPolicy: *policy,
		})
		if err != nil {
			return fmt.Errorf("failed to update access policy: %w", err)
		}
		fmt.Printf("Updated access policy for %s\n", secret.Path)
		return nil
	},
}

func init() {
	SecretsCmd.AddCommand(policyCmd)
	policyCmd.Flags().StringVar(&policyRegion, "region", "", "Region")
	policyCmd.Flags().StringVar(&policyPath, "path", "", "Secret path")
	policyCmd.Flags().StringVar(&policyFile, "file", "", "JSON file containing the access policy")
	_ = policyCmd.MarkFlagRequired("region")
	_ = policyCmd.MarkFlagRequired("path")
	_ = policyCmd.MarkFlagRequired("file")
	_ = policyCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
