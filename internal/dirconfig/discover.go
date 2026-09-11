package dirconfig

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	mu     sync.Mutex
	cached Loaded
	loaded bool
)

// Disabled reports whether directory config should be ignored.
func Disabled() bool {
	if IgnoreDirConfigFlag {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvDisable))) {
	case "0", "false", "no", "off":
		return true
	}
	return false
}

// Disable ignores directory config for the rest of this process.
// Used by the Kubernetes exec-credential plugin so a repo file cannot redirect it.
func Disable() {
	IgnoreDirConfigFlag = true
}

// Reset clears the process-wide cache. Tests only.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	cached = Loaded{}
	loaded = false
	IgnoreDirConfigFlag = false
}

// SetCurrentForTest installs an overlay without walking the filesystem. Tests only.
func SetCurrentForTest(cfg Config, path string) {
	mu.Lock()
	defer mu.Unlock()
	cfg.normalize()
	cached = Loaded{Config: cfg, Path: path}
	loaded = true
}

// Current returns the cached overlay. Missing or disabled config yields a zero Loaded value.
func Current() Loaded {
	mu.Lock()
	defer mu.Unlock()
	if Disabled() {
		return Loaded{}
	}
	if !loaded {
		_ = loadLocked()
	}
	return cached
}

// Load discovers directory config from the process working directory and caches it.
// A missing file is not an error. Invalid YAML or an oversized file is an error.
func Load() error {
	mu.Lock()
	defer mu.Unlock()
	return loadLocked()
}

func loadLocked() error {
	if loaded {
		return nil
	}
	loaded = true
	cached = Loaded{}
	if Disabled() {
		return nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	cfg, path, err := Discover(cwd)
	if err != nil {
		loaded = false
		return err
	}
	cached = Loaded{Config: cfg, Path: path}
	return nil
}

// Discover walks from cwd toward the filesystem root and returns the nearest overlay.
// It does not merge ancestor files. A missing file returns a zero Config and empty path.
func Discover(cwd string) (Config, string, error) {
	cwd = strings.TrimSpace(cwd)
	if cwd == "" {
		return Config{}, "", nil
	}
	dir, err := filepath.Abs(cwd)
	if err != nil {
		return Config{}, "", fmt.Errorf("resolve directory config path: %w", err)
	}

	for {
		cfg, path, err := configInDir(dir)
		if err != nil {
			return Config{}, "", err
		}
		if path != "" {
			return cfg, path, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return Config{}, "", nil
		}
		dir = parent
	}
}

func configInDir(dir string) (Config, string, error) {
	candidate := filepath.Join(dir, FileName)
	info, err := os.Stat(candidate)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Config{}, "", nil
		}
		return Config{}, "", fmt.Errorf("stat directory config %s: %w", candidate, err)
	}

	switch {
	case info.Mode().IsRegular():
		cfg, err := readConfig(candidate)
		if err != nil {
			return Config{}, "", err
		}
		return cfg, candidate, nil
	case info.IsDir():
		nested := filepath.Join(candidate, NestedConfigFile)
		nestedInfo, err := os.Stat(nested)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return Config{}, "", nil
			}
			return Config{}, "", fmt.Errorf("stat directory config %s: %w", nested, err)
		}
		if !nestedInfo.Mode().IsRegular() {
			return Config{}, "", nil
		}
		cfg, err := readConfig(nested)
		if err != nil {
			return Config{}, "", err
		}
		return cfg, nested, nil
	default:
		return Config{}, "", nil
	}
}

func readConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("read directory config %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	info, err := f.Stat()
	if err != nil {
		return Config{}, fmt.Errorf("stat directory config %s: %w", path, err)
	}
	if info.Size() > MaxSize {
		return Config{}, fmt.Errorf("directory config %s is larger than %d bytes", path, MaxSize)
	}

	data, err := io.ReadAll(io.LimitReader(f, MaxSize+1))
	if err != nil {
		return Config{}, fmt.Errorf("read directory config %s: %w", path, err)
	}
	if len(data) > MaxSize {
		return Config{}, fmt.Errorf("directory config %s is larger than %d bytes", path, MaxSize)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse directory config %s: %w", path, err)
	}
	cfg.normalize()
	return cfg, nil
}
