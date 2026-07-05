package contextstate

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/thalassa-cloud/cli/internal/auth/browseroidc"
)

const accessTokenRefreshSkew = 30 * time.Second

func LoginWithBrowserOIDC(ctx context.Context, apiEndpoint, contextName string) error {
	if err := EnsureLoginContext(contextName, apiEndpoint); err != nil {
		return err
	}

	fmt.Println("Opening browser to log in with Keycloak...")
	tokens, err := browseroidc.BrowserLogin(ctx, browseroidc.BrowserLoginOptions{})
	if err != nil {
		return err
	}
	if err := loginWithBrowserTokens(tokens, apiEndpoint); err != nil {
		return err
	}
	fmt.Println("Logged in successfully.")
	return nil
}

// EnsureLoginContext creates a named context when none is currently configured.
func EnsureLoginContext(name, apiEndpoint string) error {
	if _, err := GetContextConfiguration(); err == nil {
		return nil
	}

	u, err := url.Parse(apiEndpoint)
	if err != nil {
		return fmt.Errorf("invalid api endpoint: %w", err)
	}

	defaultContext := Context{
		Name: name,
		Servers: Servers{
			Name: u.Host,
			API:  API{Server: apiEndpoint},
		},
		Users: Users{
			Name: name,
			User: User{},
		},
	}
	if err := GlobalConfigManager().AddOrMergeContext(defaultContext); err != nil {
		return err
	}
	return Set(name)
}

func loginWithBrowserTokens(tokens browseroidc.TokenPair, apiEndpoint string) error {
	currentContext, err := GetContextConfiguration()
	if err != nil {
		return err
	}

	u, err := url.Parse(apiEndpoint)
	if err != nil {
		return fmt.Errorf("invalid api endpoint: %w", err)
	}

	currentContext.Users.User.SetBrowserTokens(tokens.AccessToken, tokens.RefreshToken, tokens.Expiry)
	currentContext.Servers.API.Server = u.String()
	if err := CombineConfigContext(currentContext); err != nil {
		return err
	}
	return Save()
}

// EnsureFreshAccessToken refreshes browser OIDC access tokens when they are close to expiring.
func EnsureFreshAccessToken(ctx context.Context) error {
	currentContext, err := GetContextConfiguration()
	if err != nil {
		return err
	}

	user := currentContext.Users.User
	if user.RefreshToken() == "" {
		return nil
	}
	if !user.AccessTokenExpired(time.Now(), accessTokenRefreshSkew) {
		return nil
	}

	tokens, err := browseroidc.RefreshAccessToken(ctx, user.RefreshToken())
	if err != nil {
		return fmt.Errorf("refresh browser login token: %w", err)
	}

	currentContext.Users.User.SetBrowserTokens(tokens.AccessToken, tokens.RefreshToken, tokens.Expiry)
	if err := CombineConfigContext(currentContext); err != nil {
		return err
	}
	return Save()
}
