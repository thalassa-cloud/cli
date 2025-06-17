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

	if contextstate.Debug() {
		fmt.Println("Debug mode enabled")
		fmt.Println("Context:", context)
		fmt.Println("Options:", opts)
	}

	token := contextstate.Token()
	clientID := contextstate.ClientIdOrFlag()
	clientSecret := contextstate.ClientSecretOrFlag()

	if token != "" {
		opts = append(opts, client.WithAuthPersonalToken(token))
	}

	if clientID != "" && clientSecret != "" {
		opts = append(opts, client.WithAuthOIDC(clientID, clientSecret, fmt.Sprintf("%s/oidc/token", context.Servers.API.Server)))
	}

	client, err := thalassa.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}
