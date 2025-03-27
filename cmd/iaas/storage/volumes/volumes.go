package volumes

import (
	"github.com/spf13/cobra"
)

// VolumesCmd represents the volumes command
var VolumesCmd = &cobra.Command{
	Use:     "volumes",
	Aliases: []string{"volume", "vol"},
	Short:   "Manage storage volumes",
}

func init() {
}
