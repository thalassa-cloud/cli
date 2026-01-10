package vpcpeering

import (
	"github.com/spf13/cobra"
)

// VpcPeeringCmd represents the vpc peering command
var VpcPeeringCmd = &cobra.Command{
	Use:     "vpc-peering",
	Aliases: []string{"vpcpeering"},
	Short:   "Manage VPC peering connections",
}

func init() {
}
