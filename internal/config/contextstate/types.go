package contextstate

import (
	"errors"
	"time"

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

	refreshToken      string    `yaml:"-"`
	accessTokenExpiry time.Time `yaml:"-"`
}

func (u User) Secrets() credentials.Secrets {
	secrets := credentials.Secrets{
		Token:        u.Token,
		AccessToken:  u.AccessToken,
		RefreshToken: u.refreshToken,
		ClientID:     u.ClientID,
		ClientSecret: u.ClientSecret,
	}
	if !u.accessTokenExpiry.IsZero() {
		secrets.AccessTokenExpiry = u.accessTokenExpiry.UTC().Format(time.RFC3339)
	}
	return secrets
}

func (u *User) ApplySecrets(secrets credentials.Secrets) {
	u.Token = secrets.Token
	u.AccessToken = secrets.AccessToken
	u.refreshToken = secrets.RefreshToken
	u.ClientID = secrets.ClientID
	u.ClientSecret = secrets.ClientSecret
	u.accessTokenExpiry = time.Time{}
	if secrets.AccessTokenExpiry != "" {
		if parsed, err := time.Parse(time.RFC3339, secrets.AccessTokenExpiry); err == nil {
			u.accessTokenExpiry = parsed
		}
	}
}

func (u *User) SetBrowserTokens(accessToken, refreshToken string, expiry time.Time) {
	u.Token = ""
	u.AccessToken = accessToken
	u.refreshToken = refreshToken
	u.accessTokenExpiry = expiry
	u.ClientID = ""
	u.ClientSecret = ""
}

func (u User) RefreshToken() string {
	return u.refreshToken
}

func (u User) AccessTokenExpired(now time.Time, skew time.Duration) bool {
	if u.accessTokenExpiry.IsZero() {
		return false
	}
	return !u.accessTokenExpiry.After(now.Add(skew))
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
