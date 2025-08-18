package natgateways

import (
	"github.com/spf13/cobra"
)

// NatGatewaysCmd represents the NAT gateways command
var NatGatewaysCmd = &cobra.Command{
	Use:     "natgateways",
	Aliases: []string{"natgateways", "ngw"},
	Short:   "Manage NAT gateways",
	Long:    "Manage NAT gateways within the Thalassa Cloud Platform. This command will list all the NAT gateways within your organisation.",
	Example: "tcloud networking natgateways list\ntcloud networking natgateways list --region us-west-1\ntcloud networking natgateways view ngw-123",
}

func init() {
}
