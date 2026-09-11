package dir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/dirconfig"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
)

var (
	initCluster string
	initForce   bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Write a .thalassa file in the current directory",
	Long: `Write a .thalassa file in the current directory from the effective CLI context and an optional Kubernetes cluster.

The file contains names and references only. Credentials stay in ~/.tcloud or the system credential store. Existing files are not overwritten unless --force is set.`,
	Example: "tcloud dir init --cluster prod-cluster",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}

		cfg := dirconfig.Config{
			Context:      contextstate.Name(),
			Organisation: contextstate.Organisation(),
			Project:      contextstate.Project(),
		}
		if cluster, ok := kuberesolve.PreferredClusterRef(initCluster); ok {
			cfg.Kubernetes.Cluster = cluster
		}
		if cfg.Empty() {
			return fmt.Errorf("nothing to write; set a context or pass --cluster")
		}

		path := filepath.Join(cwd, dirconfig.FileName)
		if err := dirconfig.WriteFile(path, cfg, initForce); err != nil {
			return err
		}
		fmt.Printf("Wrote %s\n", path)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initCluster, "cluster", "", "Kubernetes cluster identity, name, or slug to pin in the file")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite an existing .thalassa file")
}
