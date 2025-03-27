package context

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/pkg/thalassa"
)

// createCmd represents the create command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Thalassa Cloud",
	Long:  "Login to Thalassa Cloud using a personal access token, using the current context. Overrides the current context if --name is set.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := contextstate.Token()
		apiURL := contextstate.Server()

		if token == "" {
			return errors.New("token is required")
		}
		if apiURL == "" {
			return errors.New("api endpoint is required")
		}

		// Test the token and api endpoint
		client, err := thalassa.NewClient(
			client.WithBaseURL(apiURL),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(token),
		)
		if err != nil {
			return err
		}
		_, err = client.Me().ListMyOrganisations(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to test token and api endpoint: %w", err)
		}
		return contextstate.LoginWithAPIEndpoint(cmd.Context(), token, apiURL)
	},
}

func init() {
	ContextCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVar(&contextName, "name", "default", "name of the context")
}
