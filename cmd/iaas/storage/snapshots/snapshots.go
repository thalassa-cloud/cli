package snapshots

import (
	"github.com/spf13/cobra"
)

// SnapshotsCmd represents the snapshots command
var SnapshotsCmd = &cobra.Command{
	Use:     "snapshots",
	Aliases: []string{"snapshot", "snap"},
	Short:   "Manage volume snapshots",
}

func init() {
}
