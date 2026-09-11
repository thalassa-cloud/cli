package dirconfig

import "strings"

const (
	// FileName is the project-local config file (or directory) discovered by walking up from the working directory.
	FileName = ".thalassa"
	// NestedConfigFile is used when FileName is a directory.
	NestedConfigFile = "config.yaml"
	// MaxSize is the maximum size of a directory config file.
	MaxSize = 64 * 1024

	// EnvDisable controls whether directory config is applied. Set to 0, false, no, or off to skip it.
	EnvDisable = "THALASSA_DIR_CONFIG"
)

// IgnoreDirConfigFlag disables directory config for this process when set (bound to --ignore-dir-config).
var IgnoreDirConfigFlag bool

// Config is the project-local overlay. It may only contain names and references — never credentials.
type Config struct {
	Context      string           `yaml:"context,omitempty"`
	Organisation string           `yaml:"organisation,omitempty"`
	Project      string           `yaml:"project,omitempty"`
	Kubernetes   KubernetesConfig `yaml:"kubernetes,omitempty"`
}

// KubernetesConfig holds Kubernetes defaults for shells in this directory tree.
type KubernetesConfig struct {
	Cluster string `yaml:"cluster,omitempty"`
}

// Loaded is a discovered directory config and the path it was read from.
type Loaded struct {
	Config Config
	Path   string
}

// Empty reports whether the overlay contains no values.
func (c Config) Empty() bool {
	return strings.TrimSpace(c.Context) == "" &&
		strings.TrimSpace(c.Organisation) == "" &&
		strings.TrimSpace(c.Project) == "" &&
		strings.TrimSpace(c.Kubernetes.Cluster) == ""
}

func (c *Config) normalize() {
	c.Context = strings.TrimSpace(c.Context)
	c.Organisation = strings.TrimSpace(c.Organisation)
	c.Project = strings.TrimSpace(c.Project)
	c.Kubernetes.Cluster = strings.TrimSpace(c.Kubernetes.Cluster)
}
