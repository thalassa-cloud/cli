package browseroidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

var (
	ErrLoginTimeout  = errors.New("browser login timed out")
	ErrStateMismatch = errors.New("login state mismatch")
	ErrMissingCode   = errors.New("authorization code missing")
	ErrMissingTokens = errors.New("access token missing from token response")
)

// TokenPair holds OAuth tokens returned by browser OIDC login.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

// BrowserLoginOptions configures an interactive browser login.
type BrowserLoginOptions struct {
	RedirectURL string
	OpenURL     func(string) error
	Timeout     time.Duration
}

// BrowserLogin runs the authorization code flow with PKCE for the tcloud public client.
func BrowserLogin(ctx context.Context, opts BrowserLoginOptions) (TokenPair, error) {
	if opts.OpenURL == nil {
		opts.OpenURL = openBrowser
	}
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Minute
	}

	redirectURLString := opts.RedirectURL
	var listener net.Listener
	var err error
	if redirectURLString == "" {
		listener, redirectURLString, err = listenOnCallbackPorts(CallbackPorts)
	} else {
		redirectURL, parseErr := url.Parse(redirectURLString)
		if parseErr != nil {
			return TokenPair{}, fmt.Errorf("parse redirect url: %w", parseErr)
		}
		listener, err = net.Listen("tcp", net.JoinHostPort(defaultCallbackHost, redirectURL.Port()))
	}
	if err != nil {
		return TokenPair{}, fmt.Errorf("start callback listener: %w", err)
	}
	defer func() { _ = listener.Close() }()

	redirectURL, err := url.Parse(redirectURLString)
	if err != nil {
		return TokenPair{}, fmt.Errorf("parse redirect url: %w", err)
	}

	state, err := randomString(32)
	if err != nil {
		return TokenPair{}, err
	}
	verifier := oauth2.GenerateVerifier()

	conf := oauth2.Config{
		ClientID:    ClientID,
		Endpoint:    oauth2.Endpoint{AuthURL: authURL(), TokenURL: tokenURL()},
		RedirectURL: redirectURLString,
		Scopes:      []string{"openid", "offline_access", "profile", "email"},
	}

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	server := &http.Server{Handler: mux}

	mux.HandleFunc(redirectURL.Path, func(w http.ResponseWriter, r *http.Request) {
		if queryState := r.URL.Query().Get("state"); queryState != state {
			errCh <- ErrStateMismatch
			http.Error(w, "invalid login state", http.StatusBadRequest)
			return
		}
		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			desc := r.URL.Query().Get("error_description")
			errCh <- fmt.Errorf("authorization failed: %s: %s", errMsg, desc)
			http.Error(w, "login failed", http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- ErrMissingCode
			http.Error(w, "authorization code missing", http.StatusBadRequest)
			return
		}
		writeLoginSuccessPage(w)
		codeCh <- code
	})

	go func() {
		if serveErr := server.Serve(listener); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- serveErr
		}
	}()
	defer func() { _ = server.Shutdown(context.Background()) }()

	authRequest := conf.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)
	if err := opts.OpenURL(authRequest); err != nil {
		return TokenPair{}, fmt.Errorf("open browser: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	var code string
	select {
	case <-ctx.Done():
		return TokenPair{}, ErrLoginTimeout
	case err := <-errCh:
		return TokenPair{}, err
	case code = <-codeCh:
	}

	token, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return TokenPair{}, fmt.Errorf("exchange authorization code: %w", err)
	}
	if token.AccessToken == "" {
		return TokenPair{}, ErrMissingTokens
	}

	return TokenPair{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}, nil
}

// RefreshAccessToken uses a refresh token to obtain a new access token.
func RefreshAccessToken(ctx context.Context, refreshToken string) (TokenPair, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return TokenPair{}, errors.New("refresh token is required")
	}

	conf := oauth2.Config{
		ClientID: ClientID,
		Endpoint: oauth2.Endpoint{AuthURL: authURL(), TokenURL: tokenURL()},
	}
	tokenSource := conf.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	token, err := tokenSource.Token()
	if err != nil {
		return TokenPair{}, fmt.Errorf("refresh access token: %w", err)
	}
	if token.AccessToken == "" {
		return TokenPair{}, ErrMissingTokens
	}

	refresh := token.RefreshToken
	if refresh == "" {
		refresh = refreshToken
	}

	return TokenPair{
		AccessToken:  token.AccessToken,
		RefreshToken: refresh,
		Expiry:       token.Expiry,
	}, nil
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}

func randomString(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func listenOnCallbackPorts(ports []string) (net.Listener, string, error) {
	var lastErr error
	for _, port := range ports {
		listener, err := net.Listen("tcp", net.JoinHostPort(defaultCallbackHost, port))
		if err == nil {
			return listener, RedirectURLForPort(port), nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = errors.New("no callback ports configured")
	}
	return nil, "", lastErr
}
