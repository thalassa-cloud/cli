package dirconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const fileMode = 0o644

// WriteFile writes cfg to path. It refuses to overwrite an existing path unless force is true.
// The file form is intended to be committed, so credentials must never be included in cfg.
func WriteFile(path string, cfg Config, force bool) error {
	cfg.normalize()
	if cfg.Empty() {
		return fmt.Errorf("directory config has no values to write")
	}

	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("%s is a directory; write %s or remove the directory", path, filepath.Join(path, NestedConfigFile))
		}
		if !force {
			return fmt.Errorf("%s already exists (use --force to overwrite)", path)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal directory config: %w", err)
	}
	if err := os.WriteFile(path, data, fileMode); err != nil {
		return fmt.Errorf("write directory config %s: %w", path, err)
	}
	return nil
}
