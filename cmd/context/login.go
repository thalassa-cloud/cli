package context

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

// createCmd represents the create command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Thalassa Cloud",
	Long:  "Login to Thalassa Cloud using a personal access token, using the current context. Overrides the current context if --name is set.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := contextstate.PersonalAccessToken()
		apiURL := contextstate.Server()

		if apiURL == "" {
			return errors.New("api endpoint is required")
		}

		oidcClientID := contextstate.ClientIdOrFlag()
		oidcClientSecret := contextstate.ClientSecretOrFlag()

		tokenURL := fmt.Sprintf("%s/oidc/token", apiURL)

		opts := []client.Option{}
		if oidcClientID != "" && oidcClientSecret != "" {
			opts = append(opts, client.WithAuthOIDC(oidcClientID, oidcClientSecret, tokenURL))
		} else {
			opts = append(opts, client.WithAuthPersonalToken(token))
		}
		if len(opts) == 0 {
			return errors.New("no authentication method provided")
		}
		opts = append(opts, client.WithBaseURL(apiURL))
		opts = append(opts, client.WithOrganisation(contextstate.Organisation()))
		client, err := thalassa.NewClient(opts...)

		// Test the token and api endpoint
		if err != nil {
			return err
		}
		_, err = client.Me().ListMyOrganisations(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to test token and api endpoint: %w", err)
		}

		if oidcClientID != "" && oidcClientSecret != "" {
			return contextstate.LoginWithAPIEndpointOidc(cmd.Context(), oidcClientID, oidcClientSecret, apiURL)
		}
		return contextstate.LoginWithAPIEndpoint(cmd.Context(), token, apiURL)
	},
}

func init() {
	ContextCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVar(&contextName, "name", "default", "name of the context")
}
