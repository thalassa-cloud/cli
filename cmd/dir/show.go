package dir

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/thalassa-cloud/cli/internal/dirconfig"
)

var showCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"status", "view"},
	Short:   "Show the directory-local config in effect",
	Long:    "Show the nearest .thalassa (or .thalassa/config.yaml) found by walking up from the current directory.",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}

		cfg, path, err := dirconfig.Discover(cwd)
		if err != nil {
			return err
		}
		if path == "" {
			fmt.Println("No directory config found")
			return nil
		}

		fmt.Printf("path: %s\n", path)
		if dirconfig.Disabled() {
			fmt.Println("applied: no (ignored by --ignore-dir-config or THALASSA_DIR_CONFIG)")
		} else {
			fmt.Println("applied: yes")
		}

		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("marshal directory config: %w", err)
		}
		fmt.Print(string(data))
		return nil
	},
}
