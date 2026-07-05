package browseroidc_test

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalassa-cloud/cli/internal/auth/browseroidc"
)

func TestBrowserLogin(t *testing.T) {
	t.Setenv(browseroidc.ThalassaOIDCRealmEnvVar, "")

	var gotAuthQuery url.Values
	var gotTokenForm url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/protocol/openid-connect/auth":
			gotAuthQuery = r.URL.Query()
			redirectURL, err := url.Parse(gotAuthQuery.Get("redirect_uri"))
			require.NoError(t, err)
			q := redirectURL.Query()
			q.Set("code", "auth-code")
			q.Set("state", gotAuthQuery.Get("state"))
			redirectURL.RawQuery = q.Encode()
			http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		case "/protocol/openid-connect/token":
			require.NoError(t, r.ParseForm())
			gotTokenForm = r.Form
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token":  "access-token",
				"refresh_token": "refresh-token",
				"expires_in":    3600,
				"token_type":    "Bearer",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	t.Setenv(browseroidc.ThalassaOIDCRealmEnvVar, server.URL)

	redirectURL := "http://127.0.0.1:8765/auth/callback"
	tokens, err := browseroidc.BrowserLogin(context.Background(), browseroidc.BrowserLoginOptions{
		RedirectURL: redirectURL,
		OpenURL: func(authURL string) error {
			resp, err := http.Get(authURL)
			if err != nil {
				return err
			}
			_ = resp.Body.Close()
			return nil
		},
		Timeout: 10 * time.Second,
	})
	require.NoError(t, err)
	assert.Equal(t, "access-token", tokens.AccessToken)
	assert.Equal(t, "refresh-token", tokens.RefreshToken)
	assert.False(t, tokens.Expiry.IsZero())

	assert.Equal(t, browseroidc.ClientID, gotAuthQuery.Get("client_id"))
	assert.Equal(t, "S256", gotAuthQuery.Get("code_challenge_method"))
	assert.NotEmpty(t, gotAuthQuery.Get("code_challenge"))
	assert.Equal(t, "authorization_code", gotTokenForm.Get("grant_type"))
	assert.Equal(t, "auth-code", gotTokenForm.Get("code"))
	assert.NotEmpty(t, gotTokenForm.Get("code_verifier"))
}

func TestBrowserLoginUsesFallbackPortWhenDefaultIsBusy(t *testing.T) {
	blocker, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", browseroidc.CallbackPorts[0]))
	require.NoError(t, err)
	t.Cleanup(func() { _ = blocker.Close() })

	var gotRedirectURI string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/protocol/openid-connect/auth":
			gotRedirectURI = r.URL.Query().Get("redirect_uri")
			redirectURL, parseErr := url.Parse(gotRedirectURI)
			require.NoError(t, parseErr)
			q := redirectURL.Query()
			q.Set("code", "auth-code")
			q.Set("state", r.URL.Query().Get("state"))
			redirectURL.RawQuery = q.Encode()
			http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		case "/protocol/openid-connect/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token":  "access-token",
				"refresh_token": "refresh-token",
				"expires_in":    3600,
				"token_type":    "Bearer",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	t.Setenv(browseroidc.ThalassaOIDCRealmEnvVar, server.URL)

	_, err = browseroidc.BrowserLogin(context.Background(), browseroidc.BrowserLoginOptions{
		OpenURL: func(authURL string) error {
			resp, getErr := http.Get(authURL)
			if getErr != nil {
				return getErr
			}
			_ = resp.Body.Close()
			return nil
		},
		Timeout: 10 * time.Second,
	})
	require.NoError(t, err)
	assert.Equal(t, browseroidc.RedirectURLForPort(browseroidc.CallbackPorts[1]), gotRedirectURI)
}

func TestRefreshAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "refresh_token", r.Form.Get("grant_type"))
		assert.Equal(t, "old-refresh", r.Form.Get("refresh_token"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "new-access",
			"refresh_token": "new-refresh",
			"expires_in":    1800,
			"token_type":    "Bearer",
		})
	}))
	t.Cleanup(server.Close)
	t.Setenv(browseroidc.ThalassaOIDCRealmEnvVar, server.URL)

	tokens, err := browseroidc.RefreshAccessToken(context.Background(), "old-refresh")
	require.NoError(t, err)
	assert.Equal(t, "new-access", tokens.AccessToken)
	assert.Equal(t, "new-refresh", tokens.RefreshToken)
}
