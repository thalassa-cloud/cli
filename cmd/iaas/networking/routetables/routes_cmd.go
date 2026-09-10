package routetables

import (
	"github.com/spf13/cobra"
)

// RoutesCmd manages routes within a route table.
var RoutesCmd = &cobra.Command{
	Use:     "routes",
	Aliases: []string{"route"},
	Short:   "Manage routes in a route table",
}

func init() {
	RouteTablesCmd.AddCommand(RoutesCmd)
}
