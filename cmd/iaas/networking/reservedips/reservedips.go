package reservedips

import (
	"github.com/spf13/cobra"
)

const NoHeaderKey = "no-header"

var noHeader bool

// ReservedIPsCmd manages reserved IP addresses.
var ReservedIPsCmd = &cobra.Command{
	Use:     "reserved-ips",
	Aliases: []string{"reserved-ip", "rip"},
	Short:   "Manage reserved IP addresses",
	Long:    "Manage reserved public IP addresses that can be associated with load balancers or NAT gateways.",
	Example: "tcloud networking reserved-ips list\ntcloud networking reserved-ips create --name my-ip --region nl-ams\ntcloud networking reserved-ips associate rip-123 --nat-gateway ngw-456",
}

func init() {
}
