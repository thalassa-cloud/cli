package contextstate

import (
	"errors"

	"github.com/thalassa-cloud/cli/internal/credentials"
)

var (
	ErrContextNotFound = errors.New("context not found")
)

type Config struct {
	ConfigVersion  string             `yaml:"configVersion"`
	Contexts       []ContextReference `yaml:"contexts"`
	CurrentContext string             `yaml:"current-context"`
	Servers        []Servers          `yaml:"servers"`
	Users          []Users            `yaml:"users"`
}

type Context struct {
	Name         string
	Organisation string
	Project      string
	Servers      Servers
	Users        Users
}

type ContextReference struct {
	Name    string     `yaml:"name"`
	Context ContextRef `yaml:"context"`
}

type ContextRef struct {
	API          string `yaml:"api"`
	User         string `yaml:"user"`
	Organisation string `yaml:"organisation"`
	Project      string `yaml:"project,omitempty"`
}

type Servers struct {
	Name string `yaml:"name"`
	API  API    `yaml:"api"`
}

type API struct {
	Server string `yaml:"server"`
}

type Users struct {
	Name            string `yaml:"name"`
	CredentialStore string `yaml:"credentialStore,omitempty"`
	User            User   `yaml:"user"`
}

type User struct {
	Token        string `yaml:"token,omitempty"`
	AccessToken  string `yaml:"accessToken,omitempty"`
	ClientID     string `yaml:"clientID,omitempty"`
	ClientSecret string `yaml:"clientSecret,omitempty"`
}

func (u User) Secrets() credentials.Secrets {
	return credentials.Secrets{
		Token:        u.Token,
		AccessToken:  u.AccessToken,
		ClientID:     u.ClientID,
		ClientSecret: u.ClientSecret,
	}
}

func (u *User) ApplySecrets(secrets credentials.Secrets) {
	u.Token = secrets.Token
	u.AccessToken = secrets.AccessToken
	u.ClientID = secrets.ClientID
	u.ClientSecret = secrets.ClientSecret
}

func (u *User) ClearSecrets() {
	u.ApplySecrets(credentials.Secrets{})
}

func (u User) HasSecrets() bool {
	return !u.Secrets().Empty()
}

// Sanitized returns a copy of the config with sensitive fields redacted.
func (c Config) Sanitized() Config {
	sanitized := c.copyForSave()
	for i := range sanitized.Users {
		sanitized.Users[i].User.ClearSecrets()
	}
	return sanitized
}

func (c Config) copyForSave() Config {
	copied := c
	copied.Contexts = append([]ContextReference(nil), c.Contexts...)
	copied.Servers = append([]Servers(nil), c.Servers...)
	copied.Users = append([]Users(nil), c.Users...)
	return copied
}
