package contextstate

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/thalassa-cloud/cli/internal/auth/browseroidc"
	"github.com/thalassa-cloud/cli/internal/shared"
)

const accessTokenRefreshSkew = 30 * time.Second

// ErrSessionExpired indicates the browser OIDC session can no longer be refreshed.
var ErrSessionExpired = errors.New("browser login session expired")

var (
	refreshAccessTokenFn = browseroidc.RefreshAccessToken
	browserLoginFn       = LoginWithBrowserOIDC
	isInteractiveAuthFn  = shared.IsInteractiveAuth
	promptYesNoFn        = shared.PromptYesNoUnlessForce

	// sessionReloginAttempted prevents repeated prompts in the same process
	// after a refresh failure (for example from request middleware).
	sessionReloginAttempted bool
)

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
	ApplyPreferredCredentialStoreForNewLogin(&currentContext.Users)
	if err := CombineConfigContext(currentContext); err != nil {
		return err
	}
	return Save()
}

// EnsureFreshAccessToken refreshes browser OIDC access tokens when they are close to expiring.
// When refresh fails for a browser session, it prompts for interactive re-login when possible.
func EnsureFreshAccessToken(ctx context.Context) error {
	currentContext, err := GetContextConfiguration()
	if err != nil {
		return nil
	}

	user := currentContext.Users.User
	if user.RefreshToken() == "" {
		return nil
	}
	if !user.AccessTokenExpired(time.Now(), accessTokenRefreshSkew) {
		return nil
	}

	tokens, err := refreshAccessTokenFn(ctx, user.RefreshToken())
	if err != nil {
		return handleExpiredBrowserSession(ctx)
	}

	currentContext.Users.User.SetBrowserTokens(tokens.AccessToken, tokens.RefreshToken, tokens.Expiry)
	if err := CombineConfigContext(currentContext); err != nil {
		return err
	}
	return Save()
}

func handleExpiredBrowserSession(ctx context.Context) error {
	if sessionReloginAttempted {
		return sessionExpiredError()
	}
	sessionReloginAttempted = true

	if authOverriddenByFlagOrEnv() {
		return sessionExpiredError()
	}

	if !isInteractiveAuthFn() {
		return sessionExpiredError()
	}

	fmt.Println("Your browser login session expired.")
	proceed, err := promptYesNoFn(false, "Open browser to log in again? [y/N]: ")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSessionExpired, err)
	}
	if !proceed {
		return sessionExpiredError()
	}

	if err := browserLoginFn(ctx, Server(), Name()); err != nil {
		return err
	}
	return nil
}

func sessionExpiredError() error {
	return fmt.Errorf("%w; run: tcloud context login --browser", ErrSessionExpired)
}

func authOverriddenByFlagOrEnv() bool {
	if AccessTokenFlag != "" || os.Getenv(ThalassaAccessTokenEnvVar) != "" {
		return true
	}
	if PersonalAccessTokenFlag != "" || os.Getenv(ThalassaPersonalAccessTokenEnvVar) != "" {
		return true
	}
	if OidcClientIDFlag != "" || os.Getenv(ThalassaOIDCClientIDEnvVar) != "" {
		return true
	}
	if OidcClientSecretFlag != "" || os.Getenv(ThalassaOIDCClientSecretEnvVar) != "" {
		return true
	}
	return false
}
