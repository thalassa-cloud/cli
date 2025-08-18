package networking

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/iaas/networking/natgateways"
	"github.com/thalassa-cloud/cli/cmd/iaas/networking/routetables"
	"github.com/thalassa-cloud/cli/cmd/iaas/networking/securitygroups"
	"github.com/thalassa-cloud/cli/cmd/iaas/networking/subnets"
	"github.com/thalassa-cloud/cli/cmd/iaas/networking/vpcs"
)

// NetworkingCmd represents the networking command
var NetworkingCmd = &cobra.Command{
	Use:     "networking",
	Aliases: []string{"net", "network", "networks"},
	Short:   "Manage networking resources",
	Long:    "Manage networking resources in the Thalassa Cloud Platform",
}

func init() {
	NetworkingCmd.AddCommand(vpcs.VpcsCmd)
	NetworkingCmd.AddCommand(subnets.SubnetsCmd)
	NetworkingCmd.AddCommand(routetables.RouteTablesCmd)
	NetworkingCmd.AddCommand(natgateways.NatGatewaysCmd)
	NetworkingCmd.AddCommand(securitygroups.SecurityGroupsCmd)
}
