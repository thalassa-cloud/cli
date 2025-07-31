package contextstate

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
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
	yamlFile, err := os.ReadFile(c.filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &c.config)
	if err != nil {
		return err
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
		},
	}
	c.replaceContext(contextRef)
	return nil
}

// Save saves the configuration to the file.
func (c *configFileContextManager) Save() error {
	c.config.ConfigVersion = "v1"

	data, err := yaml.Marshal(c.config)
	if err != nil {
		return err
	}
	return os.WriteFile(c.filename, data, 0600)
}

// Config returns the current configuration.
func (c *configFileContextManager) Config() Config {
	return c.config
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
