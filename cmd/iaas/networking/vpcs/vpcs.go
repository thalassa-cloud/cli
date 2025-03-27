package vpcs

import (
	"github.com/spf13/cobra"
)

// VpcsCmd represents the incidents command
var VpcsCmd = &cobra.Command{
	Use:     "vpcs",
	Aliases: []string{"vpc", "virtualprivateclouds"},
	Short:   "Manage VPCs",
}

func init() {
}
