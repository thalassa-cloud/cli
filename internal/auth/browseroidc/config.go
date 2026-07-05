package browseroidc

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

const (
	DefaultRealmURL = "https://login.thalassa.cloud/realms/thalassa-cloud"
	ClientID        = "tcloud"

	ThalassaOIDCRealmEnvVar = "THALASSA_OIDC_REALM"

	defaultCallbackHost = "127.0.0.1"
	callbackPath        = "/auth/callback"
)

// CallbackPorts are tried in order when starting the browser login callback listener.
// Register each redirect URI in the OIDC client configuration.
var CallbackPorts = []string{"8765", "8766", "8767", "8768", "8769", "8770"}

func validatedRealmURL() (string, error) {
	realm := strings.TrimSpace(os.Getenv(ThalassaOIDCRealmEnvVar))
	if realm == "" {
		return DefaultRealmURL, nil
	}

	realm = strings.TrimRight(realm, "/")
	if err := validateRealmURL(realm); err != nil {
		return "", err
	}
	return realm, nil
}

func validateRealmURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse realm URL: %w", err)
	}
	if u.Scheme == "https" && u.Host != "" {
		return nil
	}
	if u.Scheme == "http" && isLoopbackHost(u.Hostname()) {
		return nil
	}
	return fmt.Errorf("realm URL must use https")
}

func validateLoopbackRedirectURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse redirect url: %w", err)
	}
	if u.Scheme != "http" || !isLoopbackHost(u.Hostname()) {
		return fmt.Errorf("redirect URL must use http://127.0.0.1 or http://localhost")
	}
	return nil
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func authURL() (string, error) {
	realm, err := validatedRealmURL()
	if err != nil {
		return "", err
	}
	return realm + "/protocol/openid-connect/auth", nil
}

func tokenURL() (string, error) {
	realm, err := validatedRealmURL()
	if err != nil {
		return "", err
	}
	return realm + "/protocol/openid-connect/token", nil
}

func DefaultRedirectURL() string {
	return RedirectURLForPort(CallbackPorts[0])
}

func RedirectURLForPort(port string) string {
	return "http://" + defaultCallbackHost + ":" + port + callbackPath
}
