package contextstate

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/mitchellh/go-homedir"
)

const (
	DefaultConfigFilename = ".tcloud"
	DefaultAPIURL         = "https://api.thalassa.cloud"
)

const (
	ThalassaConfigEnvVar  = "THALASSA_CONFIG"
	ThalassaCConfigEnvVar = "THALASSACONFIG"

	ThalassaAccessTokenEnvVar         = "THALASSA_ACCESS_TOKEN"
	ThalassaPersonalAccessTokenEnvVar = "THALASSA_PERSONAL_ACCESS_TOKEN"
	ThalassaOIDCClientIDEnvVar        = "THALASSA_CLIENT_ID"
	ThalassaOIDCClientSecretEnvVar    = "THALASSA_CLIENT_SECRET"
	ThalassaOrganisationIDEnvVar      = "THALASSA_ORGANISATION_ID"

	ThalassaAPIEndpointEnvVar = "THALASSA_API_ENDPOINT"
)

var (
	globalConfigManager     ConfigManager
	OrganisationFlag        string
	EndpointFlag            string
	PersonalAccessTokenFlag string

	AccessTokenFlag string

	OidcClientIDFlag     string
	OidcClientSecretFlag string

	DebugFlag   bool
	ContextFlag string
)

// ConfigManager defines an interface for managing contexts within the application.
// It provides methods to get, set, and manipulate contexts, as well as to load and save configurations.
type ConfigManager interface {
	// Get returns the current context.
	// It returns an error if there is an issue retrieving the context.
	Get() (Context, error)

	// Set sets the current context to the one specified by name.
	// It returns an error if there is an issue setting the context.
	Set(name string) error

	// AddOrMergeContext adds a new context or merges it with an existing one.
	// It returns an error if there is an issue adding or merging the context.
	AddOrMergeContext(context Context) error

	// Load loads the context configuration from a persistent storage.
	// It returns an error if there is an issue loading the configuration.
	Load() error

	// Save saves the current context configuration to a persistent storage.
	// It returns an error if there is an issue saving the configuration.
	Save() error

	// RemoveContext removes a context from the configuration.
	// It returns an error if there is an issue removing the context.
	RemoveContext(name string) error

	// RemoveContextUser removes a user from the configuration.
	// It returns an error if there is an issue removing the user.
	RemoveContextUser(name string) error

	// RemoveContextServer removes a server from the configuration.
	// It returns an error if there is an issue removing the server.
	RemoveContextServer(name string) error

	// Config returns the current configuration.
	Config() Config
}

func Init() {
	configFilename := getConfigFilename()
	globalConfigManager = NewConfigFileContextManager(configFilename)
	if err := globalConfigManager.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("Failed to initialize context: %v\n", err)
		os.Exit(1)
	}
}

func getConfigFilename() string {
	if configFilename := os.Getenv(ThalassaConfigEnvVar); configFilename != "" {
		return configFilename
	}
	if configFilename := os.Getenv(ThalassaCConfigEnvVar); configFilename != "" {
		return configFilename
	}
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Failed to get home directory: %v\n", err)
		os.Exit(1)
	}
	return fmt.Sprintf("%s/%s", home, DefaultConfigFilename)
}

func GlobalConfigManager() ConfigManager {
	return globalConfigManager
}

func GetContextConfiguration() (Context, error) {
	return globalConfigManager.Get()
}

func Set(name string) error {
	return globalConfigManager.Set(name)
}

func CombineConfigContext(context Context) error {
	return globalConfigManager.AddOrMergeContext(context)
}

func Load() error {
	return globalConfigManager.Load()
}

func Save() error {
	return globalConfigManager.Save()
}

func Organisation() string {
	if OrganisationFlag != "" {
		return OrganisationFlag
	}
	if organisation := os.Getenv(ThalassaOrganisationIDEnvVar); organisation != "" {
		return organisation
	}

	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		return ""
	}
	return currentcontext.Organisation
}

func Server() string {
	if EndpointFlag != "" {
		return EndpointFlag
	}
	if server := os.Getenv(ThalassaAPIEndpointEnvVar); server != "" {
		return server
	}

	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		return DefaultAPIURL
	}
	if currentcontext.Servers.API.Server != "" {
		return currentcontext.Servers.API.Server
	}
	return DefaultAPIURL
}

func AccessToken() string {
	if AccessTokenFlag != "" {
		return AccessTokenFlag
	}
	if accessToken := os.Getenv(ThalassaAccessTokenEnvVar); accessToken != "" {
		return accessToken
	}
	return ""
}

func PersonalAccessToken() string {
	if PersonalAccessTokenFlag != "" {
		return PersonalAccessTokenFlag
	}
	if personalAccessToken := os.Getenv(ThalassaPersonalAccessTokenEnvVar); personalAccessToken != "" {
		return personalAccessToken
	}
	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		return ""
	}

	return currentcontext.Users.User.Token
}

func ClientIdOrFlag() string {
	if OidcClientIDFlag != "" {
		return OidcClientIDFlag
	}
	if clientID := os.Getenv(ThalassaOIDCClientIDEnvVar); clientID != "" {
		return clientID
	}
	return ClientId()
}

func ClientId() string {
	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		return ""
	}
	return currentcontext.Users.User.ClientID
}

func ClientSecretOrFlag() string {
	if OidcClientSecretFlag != "" {
		return OidcClientSecretFlag
	}
	if clientSecret := os.Getenv(ThalassaOIDCClientSecretEnvVar); clientSecret != "" {
		return clientSecret
	}
	return ClientSecret()
}

func ClientSecret() string {
	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		return ""
	}
	return currentcontext.Users.User.ClientSecret
}

func GetContext() (Context, error) {
	return globalConfigManager.Get()
}

func Debug() bool {
	return DebugFlag
}
