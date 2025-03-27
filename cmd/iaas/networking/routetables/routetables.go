package routetables

import (
	"github.com/spf13/cobra"
)

// RouteTablesCmd represents the incidents command
var RouteTablesCmd = &cobra.Command{
	Use:     "routetables",
	Aliases: []string{"routetables"},
	Short:   "Manage routetables",
}

func init() {
}
