package context

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
)

var (
	// context creation options
	createContext bool
	contextName   string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new context with authentication and organisation",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		if createContext {
			err := createNewContext()
			if err != nil {
				return err
			}
		}

		token := contextstate.PersonalAccessToken()
		apiURL := contextstate.Server()
		if apiURL == "" {
			return errors.New("api endpoint is required")
		}

		oidcClientID := contextstate.ClientIdOrFlag()
		oidcClientSecret := contextstate.ClientSecretOrFlag()
		if token == "" && oidcClientID == "" && oidcClientSecret == "" {
			return errors.New("no token or oidc client id and secret set")
		}

		if oidcClientID != "" && oidcClientSecret != "" {
			if err := contextstate.LoginWithAPIEndpointOidc(ctx, oidcClientID, oidcClientSecret, apiURL); err != nil {
				return fmt.Errorf("failed to login with oidc: %w", err)
			}
		} else {
			if err := contextstate.LoginWithAPIEndpoint(ctx, token, apiURL); err != nil {
				return fmt.Errorf("failed to login with token: %w", err)
			}
		}

		var err error
		organisation := contextstate.Organisation()
		if organisation == "" {
			organisation, err = getSelectedOrganisation([]string{})
			if err != nil {
				return fmt.Errorf("failed to get selected organisation: %w", err)
			}
		}

		currentContext, err := contextstate.GetContextConfiguration()
		if err != nil {
			return err
		}
		currentContext.Organisation = organisation
		err = contextstate.CombineConfigContext(currentContext)
		if err != nil {
			return err
		}
		return contextstate.Save()
	},
}

func createNewContext() error {
	u, err := url.Parse(contextstate.Server())
	if err != nil {
		return err
	}

	defaultContext := newDefaultContext(contextName, contextstate.OrganisationFlag, u.Host)
	err = contextstate.GlobalConfigManager().AddOrMergeContext(defaultContext)
	if err != nil {
		return err
	}
	err = contextstate.Set(defaultContext.Name)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	ContextCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&createContext, "create-context", true, "creates a context")
	createCmd.Flags().StringVar(&contextName, "name", "default", "name of the context")
}

func newDefaultContext(contextName, organisation, apiName string) contextstate.Context {
	return contextstate.Context{
		Name:         contextName,
		Organisation: organisation,
		Servers: contextstate.Servers{
			Name: apiName,
			API: contextstate.API{
				Server: contextstate.Server(),
			},
		},
		Users: contextstate.Users{
			Name: contextName,
			User: contextstate.User{
				Token: contextstate.PersonalAccessToken(),
			},
		},
	}
}
