package zones

import (
	"github.com/spf13/cobra"
)

// ZonesCmd manages DNS zones.
var ZonesCmd = &cobra.Command{
	Use:     "zones",
	Aliases: []string{"zone", "z"},
	Short:   "Manage DNS zones",
}

var (
	noHeader      bool
	showExactTime bool
)
