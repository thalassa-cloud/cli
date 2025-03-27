package subnets

import (
	"github.com/spf13/cobra"
)

// SubnetsCmd represents the incidents command
var SubnetsCmd = &cobra.Command{
	Use:     "subnets",
	Aliases: []string{"subnets"},
	Short:   "Manage subnets",
}

func init() {
}
