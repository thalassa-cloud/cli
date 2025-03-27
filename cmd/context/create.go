package context

import (
	"context"
	"errors"
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

		token := contextstate.Token()
		if token == "" {
			return errors.New("no token set")
		}
		err := contextstate.LoginWithAPIEndpoint(ctx, token, contextstate.Server())
		if err != nil {
			return err
		}

		organisation := contextstate.Organisation()
		if organisation == "" {
			organisation, err = getSelectedOrganisation([]string{})
			if err != nil {
				return err
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
			Name: "token", // TODO: username from token
			User: contextstate.User{
				Token: contextstate.Token(),
			},
		},
	}
}
