package federatedidentityproviders

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var getCmd = &cobra.Command{
	Use:               "get <identity>",
	Short:             "Show a federated identity provider",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMFederatedIdentityProviderIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		p, err := client.IAM().GetFederatedIdentityProvider(ctx, args[0])
		if err != nil {
			return fmt.Errorf("failed to get provider: %w", err)
		}
		fmt.Printf("Identity:    %s\n", p.Identity)
		fmt.Printf("Name:        %s\n", p.Name)
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("Issuer:      %s\n", p.ProviderIssuer)
		if p.ProviderJwksURI != nil {
			fmt.Printf("JWKS URI:    %s\n", *p.ProviderJwksURI)
		}
		fmt.Printf("Status:      %s\n", p.Status)
		fmt.Printf("Created:     %s\n", formattime.FormatTime(p.CreatedAt.Local(), showExactTime))
		return nil
	},
}

func init() {
	FederatedIdentityProvidersCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	getCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
