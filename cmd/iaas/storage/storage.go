package storage

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/iaas/storage/volumes"
)

// StorageCmd represents the storage command
var StorageCmd = &cobra.Command{
	Use:     "storage",
	Aliases: []string{"store", "vol", "volume", "volumes"},
	Short:   "Manage storage resources",
}

func init() {
	StorageCmd.AddCommand(volumes.VolumesCmd)
}
