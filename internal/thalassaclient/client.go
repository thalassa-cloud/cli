package thalassaclient

import (
	"errors"
	"fmt"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/version"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

func GetThalassaClient() (thalassa.Client, error) {
	var org string
	context, err := contextstate.GetContext()
	if err != nil {
		if !errors.Is(err, contextstate.ErrContextNotFound) {
			return nil, fmt.Errorf("failed to get context: %w", err)
		}
	}
	endpoint := "https://api.thalassa.cloud"
	if context.Servers.API.Server != "" {
		endpoint = context.Servers.API.Server
	}

	opts := []client.Option{
		client.WithBaseURL(endpoint),
		client.WithUserAgent(version.UserAgent()),
	}
	org = context.Organisation
	if org != "" {
		opts = append(opts, client.WithOrganisation(org))
	}

	if contextstate.Debug() {
		fmt.Println("Debug mode enabled")
		fmt.Println("Context:", context)
		fmt.Println("Options:", opts)
	}

	token := contextstate.PersonalAccessToken()
	clientID := contextstate.ClientIdOrFlag()
	clientSecret := contextstate.ClientSecretOrFlag()
	accessToken := contextstate.AccessToken()
	if accessToken != "" {
		opts = append(opts, client.WithToken(accessToken))
	} else if clientID != "" && clientSecret != "" {
		opts = append(opts, client.WithAuthOIDC(clientID, clientSecret, fmt.Sprintf("%s/oidc/token", endpoint)))
	} else if token != "" {
		opts = append(opts, client.WithAuthPersonalToken(token))
	} else {
		return nil, errors.New("no authentication method provided")
	}

	client, err := thalassa.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return client, nil
}
