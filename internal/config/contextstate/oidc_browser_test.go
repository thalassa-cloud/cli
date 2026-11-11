package contextstate

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/thalassa-cloud/cli/internal/auth/browseroidc"
	"github.com/thalassa-cloud/cli/internal/credentials"
)

func TestEnsureFreshAccessTokenWithoutContext(t *testing.T) {
	path := t.TempDir() + "/.tcloud"
	globalConfigManager = NewConfigFileContextManager(path)

	err := EnsureFreshAccessToken(context.Background())
	require.NoError(t, err)
}

func TestEnsureFreshAccessTokenRefreshFailure(t *testing.T) {
	tests := []struct {
		name             string
		interactive      bool
		promptProceed    bool
		promptErr        error
		overrideAccess   bool
		wantErrIs        error
		wantLoginCalls   int
		wantErrContains  string
		setupRefreshErr  bool
		skipExpiredToken bool
	}{
		{
			name:            "non-interactive returns session expired hint",
			interactive:     false,
			setupRefreshErr: true,
			wantErrIs:       ErrSessionExpired,
			wantErrContains: "tcloud context login --browser",
			wantLoginCalls:  0,
		},
		{
			name:            "interactive decline returns session expired",
			interactive:     true,
			promptProceed:   false,
			setupRefreshErr: true,
			wantErrIs:       ErrSessionExpired,
			wantErrContains: "tcloud context login --browser",
			wantLoginCalls:  0,
		},
		{
			name:            "interactive accept calls browser login",
			interactive:     true,
			promptProceed:   true,
			setupRefreshErr: true,
			wantLoginCalls:  1,
		},
		{
			name:            "auth override skips prompt and browser login",
			interactive:     true,
			promptProceed:   true,
			overrideAccess:  true,
			setupRefreshErr: true,
			wantErrIs:       ErrSessionExpired,
			wantErrContains: "tcloud context login --browser",
			wantLoginCalls:  0,
		},
		{
			name:             "no refresh needed when token still valid",
			skipExpiredToken: true,
			setupRefreshErr:  true,
			wantLoginCalls:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetOIDCBrowserTestHooks(t)
			setupBrowserContext(t, !tt.skipExpiredToken)

			refreshCalls := 0
			refreshAccessTokenFn = func(ctx context.Context, refreshToken string) (browseroidc.TokenPair, error) {
				refreshCalls++
				if tt.setupRefreshErr {
					return browseroidc.TokenPair{}, errors.New("invalid_grant")
				}
				return browseroidc.TokenPair{
					AccessToken:  "new-access",
					RefreshToken: refreshToken,
					Expiry:       time.Now().Add(time.Hour),
				}, nil
			}

			loginCalls := 0
			browserLoginFn = func(ctx context.Context, apiEndpoint, contextName string) error {
				loginCalls++
				return nil
			}

			isInteractiveAuthFn = func() bool { return tt.interactive }
			promptYesNoFn = func(force bool, question string) (bool, error) {
				if tt.promptErr != nil {
					return false, tt.promptErr
				}
				return tt.promptProceed, nil
			}

			if tt.overrideAccess {
				AccessTokenFlag = "override-token"
				t.Cleanup(func() { AccessTokenFlag = "" })
			}

			err := EnsureFreshAccessToken(context.Background())
			if tt.wantErrIs != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErrIs)
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantLoginCalls, loginCalls)
			if tt.skipExpiredToken {
				assert.Equal(t, 0, refreshCalls)
			}
		})
	}
}

func TestEnsureFreshAccessTokenWithoutRefreshToken(t *testing.T) {
	resetOIDCBrowserTestHooks(t)
	keyring.MockInit()
	t.Cleanup(credentials.ResetKeyring)
	t.Setenv(credentials.ThalassaCredentialStoreEnvVar, credentials.StoreKeychain)

	path := filepath.Join(t.TempDir(), ".tcloud")
	globalConfigManager = NewConfigFileContextManager(path)
	require.NoError(t, globalConfigManager.AddOrMergeContext(Context{
		Name: "default",
		Users: Users{
			Name: "default",
			User: User{AccessToken: "access-only"},
		},
		Servers: Servers{Name: "api", API: API{Server: "https://api.example.com"}},
	}))
	require.NoError(t, globalConfigManager.Set("default"))

	refreshCalled := false
	refreshAccessTokenFn = func(ctx context.Context, refreshToken string) (browseroidc.TokenPair, error) {
		refreshCalled = true
		return browseroidc.TokenPair{}, errors.New("should not be called")
	}

	err := EnsureFreshAccessToken(context.Background())
	require.NoError(t, err)
	assert.False(t, refreshCalled)
}

func TestEnsureFreshAccessTokenSkipsRepeatedReloginPrompt(t *testing.T) {
	resetOIDCBrowserTestHooks(t)
	setupBrowserContext(t, true)

	refreshAccessTokenFn = func(ctx context.Context, refreshToken string) (browseroidc.TokenPair, error) {
		return browseroidc.TokenPair{}, errors.New("invalid_grant")
	}
	isInteractiveAuthFn = func() bool { return true }
	promptCalls := 0
	promptYesNoFn = func(force bool, question string) (bool, error) {
		promptCalls++
		return false, nil
	}
	browserLoginFn = func(ctx context.Context, apiEndpoint, contextName string) error {
		t.Fatal("browser login should not be called")
		return nil
	}

	err := EnsureFreshAccessToken(context.Background())
	require.ErrorIs(t, err, ErrSessionExpired)
	assert.Equal(t, 1, promptCalls)

	err = EnsureFreshAccessToken(context.Background())
	require.ErrorIs(t, err, ErrSessionExpired)
	assert.Equal(t, 1, promptCalls, "should not prompt again in the same process")
}

func setupBrowserContext(t *testing.T, expired bool) {
	t.Helper()
	keyring.MockInit()
	t.Cleanup(credentials.ResetKeyring)
	t.Setenv(credentials.ThalassaCredentialStoreEnvVar, credentials.StoreKeychain)
	t.Setenv(ThalassaAccessTokenEnvVar, "")
	t.Setenv(ThalassaPersonalAccessTokenEnvVar, "")
	t.Setenv(ThalassaOIDCClientIDEnvVar, "")
	t.Setenv(ThalassaOIDCClientSecretEnvVar, "")
	AccessTokenFlag = ""
	PersonalAccessTokenFlag = ""
	OidcClientIDFlag = ""
	OidcClientSecretFlag = ""

	path := filepath.Join(t.TempDir(), ".tcloud")
	globalConfigManager = NewConfigFileContextManager(path)

	require.NoError(t, globalConfigManager.AddOrMergeContext(Context{
		Name: "default",
		Users: Users{
			Name: "default",
			User: User{},
		},
		Servers: Servers{Name: "api", API: API{Server: "https://api.example.com"}},
	}))
	require.NoError(t, globalConfigManager.Set("default"))

	ctx, err := globalConfigManager.Get()
	require.NoError(t, err)
	expiry := time.Now().Add(time.Hour)
	if expired {
		expiry = time.Now().Add(-time.Minute)
	}
	ctx.Users.User.SetBrowserTokens("access-token", "refresh-token", expiry)
	require.NoError(t, CombineConfigContext(ctx))
	require.NoError(t, Save())
}

func resetOIDCBrowserTestHooks(t *testing.T) {
	t.Helper()
	prevRefresh := refreshAccessTokenFn
	prevLogin := browserLoginFn
	prevInteractive := isInteractiveAuthFn
	prevPrompt := promptYesNoFn
	prevAttempted := sessionReloginAttempted

	t.Cleanup(func() {
		refreshAccessTokenFn = prevRefresh
		browserLoginFn = prevLogin
		isInteractiveAuthFn = prevInteractive
		promptYesNoFn = prevPrompt
		sessionReloginAttempted = prevAttempted
		AccessTokenFlag = ""
		PersonalAccessTokenFlag = ""
		OidcClientIDFlag = ""
		OidcClientSecretFlag = ""
	})

	sessionReloginAttempted = false
}
