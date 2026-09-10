package routetables

import (
	"github.com/spf13/cobra"
)

// RouteTablesCmd represents the route tables command
var RouteTablesCmd = &cobra.Command{
	Use:     "routetables",
	Aliases: []string{"route-tables", "rt"},
	Short:   "Manage route tables",
	Long:    "Manage VPC route tables and their routes within the Thalassa Cloud Platform.",
	Example: "tcloud networking routetables list\ntcloud networking routetables create --name custom --vpc vpc-123\ntcloud networking routetables routes list rt-123",
}

func init() {
}
