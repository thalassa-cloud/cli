package federatedidentities

import "github.com/spf13/cobra"

// FederatedIdentitiesCmd manages federated OIDC identities.
var FederatedIdentitiesCmd = &cobra.Command{
	Use:     "federated-identities",
	Aliases: []string{"fed-ids", "federated-identity"},
	Short:   "Federated identities (OIDC subject bindings)",
}
