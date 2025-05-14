package contextstate

import "errors"

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
}

type Servers struct {
	Name string `yaml:"name"`
	API  API    `yaml:"api"`
}

type API struct {
	Server string `yaml:"server"`
}

type Users struct {
	Name string `yaml:"name"`
	User User   `yaml:"user"`
}

type User struct {
	Token        string `yaml:"token"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}
