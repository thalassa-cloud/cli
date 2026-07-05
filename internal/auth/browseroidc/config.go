package browseroidc

import (
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

func RealmURL() string {
	if realm := strings.TrimSpace(os.Getenv(ThalassaOIDCRealmEnvVar)); realm != "" {
		return strings.TrimRight(realm, "/")
	}
	return DefaultRealmURL
}

func authURL() string {
	return RealmURL() + "/protocol/openid-connect/auth"
}

func tokenURL() string {
	return RealmURL() + "/protocol/openid-connect/token"
}

func DefaultRedirectURL() string {
	return RedirectURLForPort(CallbackPorts[0])
}

func RedirectURLForPort(port string) string {
	return "http://" + defaultCallbackHost + ":" + port + callbackPath
}
