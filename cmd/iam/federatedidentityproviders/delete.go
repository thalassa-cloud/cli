package federatedidentityproviders

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <identity>",
	Short:             "Delete a federated identity provider",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMFederatedIdentityProviderIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to delete this federated identity provider?\n  Identity: %s\n", args[0]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := client.IAM().DeleteFederatedIdentityProvider(ctx, args[0]); err != nil {
			return fmt.Errorf("failed to delete provider: %w", err)
		}
		fmt.Printf("Deleted federated identity provider %s\n", args[0])
		return nil
	},
}

func init() {
	FederatedIdentityProvidersCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
	deleteCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
