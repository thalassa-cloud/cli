package listeners

import "github.com/spf13/cobra"

const (
	NoHeaderKey      = "no-header"
	LoadbalancerFlag = "loadbalancer"
)

// ListenersCmd manages load balancer listeners.
var ListenersCmd = &cobra.Command{
	Use:     "listeners",
	Aliases: []string{"listener"},
	Short:   "Manage load balancer listeners",
	Long:    "Manage listeners on a load balancer. All commands require --loadbalancer.",
	Example: "tcloud networking loadbalancers listeners list --loadbalancer lb-123\ntcloud networking loadbalancers listeners create --loadbalancer lb-123 --name http --port 80 --protocol http --target-group tg-123",
}
