package thalassaclient

import (
	"fmt"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

func GetThalassaClient() (thalassa.Client, error) {
	context, err := contextstate.GetContext()
	if err != nil {
		return nil, err
	}

	opts := []client.Option{
		client.WithBaseURL(context.Servers.API.Server),
		client.WithOrganisation(context.Organisation),
	}

	if context.Users.User.Token != "" {
		opts = append(opts, client.WithAuthPersonalToken(context.Users.User.Token))
	}

	if context.Users.User.ClientID != "" && context.Users.User.ClientSecret != "" {
		opts = append(opts, client.WithAuthOIDC(context.Users.User.ClientID, context.Users.User.ClientSecret, fmt.Sprintf("%s/oidc/token", context.Servers.API.Server)))
	}

	client, err := thalassa.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}
