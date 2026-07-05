package context

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Thalassa Cloud",
	Long:  "Login to Thalassa Cloud using browser OIDC (default when no credentials are provided), a personal access token, access token, or OIDC client credentials for the current context.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiURL := contextstate.Server()
		if apiURL == "" {
			return errors.New("api endpoint is required")
		}

		token := contextstate.PersonalAccessToken()
		accessToken := contextstate.AccessToken()
		oidcClientID := contextstate.ClientIdOrFlag()
		oidcClientSecret := contextstate.ClientSecretOrFlag()

		if contextstate.BrowserLoginFlag || (token == "" && oidcClientID == "" && oidcClientSecret == "" && accessToken == "") {
			return contextstate.LoginWithBrowserOIDC(cmd.Context(), apiURL, contextName)
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
	loginCmd.Flags().BoolVar(&contextstate.BrowserLoginFlag, "browser", false, "log in through the browser using OIDC (default when no other credentials are provided)")
}
