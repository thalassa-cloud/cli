package federatedidentityproviders

import "github.com/spf13/cobra"

// FederatedIdentityProvidersCmd manages federated OIDC providers.
var FederatedIdentityProvidersCmd = &cobra.Command{
	Use:     "federated-identity-providers",
	Aliases: []string{"fed-providers", "federated-providers"},
	Short:   "Federated OIDC identity providers",
}
