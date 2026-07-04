package contextstate

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/thalassa-cloud/cli/internal/config/securefile"
	"github.com/thalassa-cloud/cli/internal/credentials"
)

type configFileContextManager struct {
	filename string
	config   Config
}

// NewConfigFileContextManager creates a new context manager with the given filename.
func NewConfigFileContextManager(filename string) ConfigManager {
	return &configFileContextManager{
		filename: filename,
	}
}

// Get returns the current context.
func (c *configFileContextManager) Get() (Context, error) {
	if c.config.CurrentContext == "" {
		return Context{}, errors.New("no current context set in config")
	}

	contextRef, ok := c.getContextRef(c.config.CurrentContext)
	if !ok {
		return Context{}, fmt.Errorf("missing current context %q in config", c.config.CurrentContext)
	}

	api, ok := c.getAPI(contextRef.Context.API)
	if !ok {
		return Context{}, fmt.Errorf("missing api %q in config", contextRef.Context.API)
	}

	user, ok := c.getUser(contextRef.Context.User)
	if !ok {
		return Context{}, fmt.Errorf("missing user %q in config", contextRef.Context.User)
	}

	return Context{
		Name:         contextRef.Name,
		Organisation: contextRef.Context.Organisation,
		Project:      contextRef.Context.Project,
		Servers:      api,
		Users:        user,
	}, nil
}

// Set sets the current context to the given name.
func (c *configFileContextManager) Set(name string) error {
	if _, ok := c.getContextRef(name); !ok {
		return ErrContextNotFound
	}

	c.config.CurrentContext = name
	return nil
}

// Load loads the configuration from the file.
func (c *configFileContextManager) Load() error {
	if _, err := os.Stat(c.filename); err == nil {
		warnPermissiveConfigPermissions(c.filename)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	yamlFile, err := os.ReadFile(c.filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &c.config)
	if err != nil {
		return err
	}

	migrated, err := c.syncCredentialsAfterLoad()
	if err != nil {
		return err
	}
	if migrated {
		return c.Save()
	}
	return nil
}

// AddOrMergeContext adds or merges the given context.
func (c *configFileContextManager) AddOrMergeContext(context Context) error {
	c.setUser(context.Users)
	c.replaceAPI(context.Servers)
	contextRef := ContextReference{
		Name: context.Name,
		Context: ContextRef{
			API:          context.Servers.Name,
			User:         context.Users.Name,
			Organisation: context.Organisation,
			Project:      context.Project,
		},
	}
	c.replaceContext(contextRef)
	return nil
}

// Save saves the configuration to the file.
func (c *configFileContextManager) Save() error {
	c.config.ConfigVersion = "v1"

	configToSave, err := c.prepareConfigForSave()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(configToSave)
	if err != nil {
		return err
	}
	return securefile.Write(c.filename, data)
}

// Config returns the current configuration.
func (c *configFileContextManager) Config() Config {
	return c.config
}

func (c *configFileContextManager) SanitizedConfig() Config {
	return c.config.Sanitized()
}

// RemoveContext removes a context from the configuration.
func (c *configFileContextManager) RemoveContext(name string) error {
	fmt.Println("Removing context", name)
	c.removeContext(name)
	return c.Save()
}

// RemoveContextUser removes a user from the configuration.
func (c *configFileContextManager) RemoveContextUser(name string) error {
	// make sure there is no context using this user
	for _, context := range c.config.Contexts {
		if context.Context.User == name {
			return fmt.Errorf("cannot remove user %q as it is still in use by context %q", name, context.Name)
		}
	}
	fmt.Println("Removing user", name)
	user, ok := c.getUser(name)
	if ok {
		if err := c.deleteUserCredentials(user); err != nil {
			return err
		}
	}
	c.deleteUser(name)
	return c.Save()
}

// RemoveContextServer removes a server from the configuration.
func (c *configFileContextManager) RemoveContextServer(name string) error {
	// make sure there is no context using this server
	for _, context := range c.config.Contexts {
		if context.Context.API == name {
			return fmt.Errorf("cannot remove server %q as it is still in use by context %q", name, context.Name)
		}
	}
	fmt.Println("Removing server", name)
	c.removeAPI(name)
	return c.Save()
}

func (c *configFileContextManager) syncCredentialsAfterLoad() (bool, error) {
	preference := credentials.PreferredStore()
	migrated := false

	for i := range c.config.Users {
		user := &c.config.Users[i]
		if user.User.HasSecrets() && credentials.UseKeychain(preference, user.CredentialStore) {
			store := credentials.ResolveStore(credentials.StoreKeychain)
			if err := store.Set(user.Name, user.User.Secrets()); err != nil {
				if preference == credentials.StoreKeychain {
					return false, fmt.Errorf("store credentials in keychain: %w", err)
				}
				user.CredentialStore = credentials.StoreFile
				continue
			}
			user.CredentialStore = credentials.StoreKeychain
			migrated = true
			continue
		}

		if user.CredentialStore != credentials.StoreKeychain {
			continue
		}

		store := credentials.ResolveStore(credentials.StoreKeychain)
		secrets, err := store.Get(user.Name)
		if err != nil {
			if errors.Is(err, credentials.ErrNotAvailable) || preference == credentials.StoreFile {
				user.CredentialStore = credentials.StoreFile
				continue
			}
			return false, fmt.Errorf("load credentials for user %q from keychain: %w", user.Name, err)
		}
		user.User.ApplySecrets(secrets)
	}

	return migrated, nil
}

func (c *configFileContextManager) prepareConfigForSave() (Config, error) {
	preference := credentials.PreferredStore()
	configToSave := c.config.copyForSave()

	for i := range c.config.Users {
		inMemoryUser := c.config.Users[i]
		userToSave := &configToSave.Users[i]

		useKeychain := credentials.UseKeychain(preference, inMemoryUser.CredentialStore)
		if useKeychain {
			store := credentials.ResolveStore(credentials.StoreKeychain)
			if err := store.Set(inMemoryUser.Name, inMemoryUser.User.Secrets()); err != nil {
				if preference == credentials.StoreKeychain {
					return Config{}, fmt.Errorf("store credentials in keychain: %w", err)
				}
				userToSave.CredentialStore = credentials.StoreFile
				c.config.Users[i].CredentialStore = credentials.StoreFile
				continue
			}
			userToSave.CredentialStore = credentials.StoreKeychain
			userToSave.User.ClearSecrets()
			c.config.Users[i].CredentialStore = credentials.StoreKeychain
			continue
		}

		userToSave.CredentialStore = credentials.StoreFile
		c.config.Users[i].CredentialStore = credentials.StoreFile
	}

	return configToSave, nil
}

func (c *configFileContextManager) deleteUserCredentials(user Users) error {
	if user.CredentialStore != credentials.StoreKeychain {
		return nil
	}
	store := credentials.ResolveStore(credentials.StoreKeychain)
	return store.Delete(user.Name)
}

func (c *configFileContextManager) FixPermissions() error {
	if _, err := os.Stat(c.filename); err != nil {
		return err
	}
	return securefile.EnsurePermissions(c.filename)
}

func warnPermissiveConfigPermissions(filename string) {
	mode, permissive, err := securefile.CheckPermissions(filename)
	if err != nil || !permissive {
		return
	}
	fmt.Fprintf(
		os.Stderr,
		"warning: config file %s has permissive permissions (%#o); run 'tcloud context fix' to restrict access\n",
		filename,
		mode,
	)
}

// -----------

func (c *configFileContextManager) getContextRef(name string) (ContextReference, bool) {
	for _, context := range c.config.Contexts {
		if context.Name == name {
			return context, true
		}
	}
	return ContextReference{}, false
}

func (c *configFileContextManager) getAPI(name string) (Servers, bool) {
	for _, api := range c.config.Servers {
		if api.Name == name {
			return api, true
		}
	}
	return Servers{}, false
}

func (c *configFileContextManager) getUser(name string) (Users, bool) {
	for _, user := range c.config.Users {
		if user.Name == name {
			return user, true
		}
	}
	return Users{}, false
}

func (c *configFileContextManager) replaceContext(contextRef ContextReference) {
	c.removeContext(contextRef.Name)
	c.config.Contexts = append(c.config.Contexts, contextRef)
}

func (c *configFileContextManager) removeContext(name string) {
	for i, ctx := range c.config.Contexts {
		if ctx.Name == name {
			c.config.Contexts = removeItemFromSlice(c.config.Contexts, i)
			break
		}
	}
}

func (c *configFileContextManager) replaceAPI(api Servers) {
	c.removeAPI(api.Name)
	c.config.Servers = append(c.config.Servers, api)
}

func (c *configFileContextManager) removeAPI(name string) {
	for i, a := range c.config.Servers {
		if a.Name == name {
			c.config.Servers = removeItemFromSlice(c.config.Servers, i)
			break
		}
	}
}

func (c *configFileContextManager) setUser(user Users) {
	c.deleteUser(user.Name)
	c.config.Users = append(c.config.Users, user)
}

func (c *configFileContextManager) deleteUser(name string) {
	for i, u := range c.config.Users {
		if u.Name == name {
			c.config.Users = removeItemFromSlice(c.config.Users, i)
			break
		}
	}
}

func removeItemFromSlice[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
