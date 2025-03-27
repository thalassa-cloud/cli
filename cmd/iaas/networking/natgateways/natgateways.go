package natgateways

import (
	"github.com/spf13/cobra"
)

// NatGatewaysCmd represents the incidents command
var NatGatewaysCmd = &cobra.Command{
	Use:     "natgateways",
	Aliases: []string{"natgateways", "ngw"},
	Short:   "Manage natgateways",
	Long:    "Manage natgateways to manage your NAT gateways within the Thalassa Cloud Platform. This command will list all the NAT gateways within your organisation.",
}

func init() {
}
