package tfs

import (
	"github.com/spf13/cobra"
)

// TfsCmd represents the tfs command
var TfsCmd = &cobra.Command{
	Use:     "tfs",
	Aliases: []string{"filesystem", "fs"},
	Short:   "Manage TFS (Thalassa File System) instances",
}

func init() {
}
