package thalassaclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/projectresolve"
	"github.com/thalassa-cloud/cli/internal/version"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

// Scope optionally overrides organisation and project from the active CLI context.
type Scope struct {
	Organisation string
	Project      string
}

// ScopeFromContext captures the organisation and project currently in effect.
func ScopeFromContext() Scope {
	return Scope{
		Organisation: contextstate.Organisation(),
		Project:      contextstate.Project(),
	}
}

func GetThalassaClient() (thalassa.Client, error) {
	return GetThalassaClientWithScope(Scope{})
}

func GetThalassaClientWithScope(scope Scope) (thalassa.Client, error) {
	if err := contextstate.EnsureFreshAccessToken(context.Background()); err != nil {
		return nil, err
	}

	endpoint := contextstate.Server()
	org := scope.Organisation
	if org == "" {
		org = contextstate.Organisation()
	}
	projectRef := scope.Project
	if projectRef == "" {
		projectRef = contextstate.Project()
	}

	opts := []client.Option{
		client.WithBaseURL(endpoint),
		client.WithUserAgent(version.UserAgent()),
	}
	if org != "" {
		opts = append(opts, client.WithOrganisation(org))
	}

	if contextstate.Debug() {
		fmt.Println("Debug mode enabled")
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

	if projectRef != "" {
		tempClient, err := thalassa.NewClient(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}
		projectIdentity, err := projectresolve.ResolveProjectIdentity(context.Background(), tempClient.Projects(), projectRef)
		if err != nil {
			return nil, err
		}
		opts = append(opts, client.WithProject(projectIdentity))
	}

	apiClient, err := thalassa.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return apiClient, nil
}
