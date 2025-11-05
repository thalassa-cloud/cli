package context

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// createCmd represents the create command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Thalassa Cloud",
	Long:  "Login to Thalassa Cloud using a personal access token, access token, or OIDC client id and secret, using the current context. Overrides the current context if --name is set.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := contextstate.PersonalAccessToken()
		accessToken := contextstate.AccessToken()
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
		} else if accessToken != "" {
			opts = append(opts, client.WithToken(accessToken))
		} else if token != "" {
			opts = append(opts, client.WithAuthPersonalToken(token))
		}
		if len(opts) == 0 {
			return errors.New("no authentication method provided")
		}
		opts = append(opts, client.WithBaseURL(apiURL))
		if contextstate.Organisation() != "" {
			opts = append(opts, client.WithOrganisation(contextstate.Organisation()))
		}

		if oidcClientID != "" && oidcClientSecret != "" {
			return contextstate.LoginWithAPIEndpointOidc(cmd.Context(), oidcClientID, oidcClientSecret, apiURL)
		}
		if accessToken != "" {
			if strings.HasPrefix(accessToken, "tc_pat_") {
				return errors.New("access token is a personal access token, use 'tcloud context login --token <token>' to login with a personal access token")
			}
			return contextstate.LoginWithAccessToken(cmd.Context(), accessToken, apiURL)
		}
		return contextstate.LoginWithAPIEndpoint(cmd.Context(), token, apiURL)
	},
}

func init() {
	ContextCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVar(&contextName, "name", "default", "name of the context")
}
