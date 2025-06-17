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
)

var (
	globalConfigManager     ConfigManager
	OrganisationFlag        string
	EndpointFlag            string
	PersonalAccessTokenFlag string

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
	if configFilename := os.Getenv("THALASSA_CONFIG"); configFilename != "" {
		return configFilename
	}
	if configFilename := os.Getenv("THALASSACONFIG"); configFilename != "" {
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

	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	organisation := currentcontext.Organisation
	return organisation
}

func Server() string {
	if EndpointFlag != "" {
		return EndpointFlag
	}

	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return currentcontext.Servers.API.Server
}

func Token() string {
	if PersonalAccessTokenFlag != "" {
		return PersonalAccessTokenFlag
	}

	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return currentcontext.Users.User.Token
}

func ClientIdOrFlag() string {
	if OidcClientIDFlag != "" {
		return OidcClientIDFlag
	}
	return ClientId()
}

func ClientId() string {
	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return currentcontext.Users.User.ClientID
}

func ClientSecretOrFlag() string {
	if OidcClientSecretFlag != "" {
		return OidcClientSecretFlag
	}
	return ClientSecret()
}

func ClientSecret() string {
	currentcontext, err := globalConfigManager.Get()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return currentcontext.Users.User.ClientSecret
}

func GetContext() (Context, error) {
	return globalConfigManager.Get()
}

func Debug() bool {
	return DebugFlag
}
