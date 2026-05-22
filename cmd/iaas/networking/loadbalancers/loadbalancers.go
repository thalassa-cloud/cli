package loadbalancers

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/iaas/networking/loadbalancers/listeners"
)

const (
	NoHeaderKey        = "no-header"
	LoadbalancerFlag   = "loadbalancer"
)

// LoadbalancersCmd represents the load balancers command.
var LoadbalancersCmd = &cobra.Command{
	Use:     "loadbalancers",
	Aliases: []string{"loadbalancer", "load-balancer", "load-balancers", "lb", "lbs"},
	Short:   "Manage load balancers",
	Long:    "Manage load balancers, listeners, and related networking resources within the Thalassa Cloud Platform.",
	Example: "tcloud networking loadbalancers list\ntcloud networking loadbalancers create --name web --subnet subnet-123\ntcloud networking loadbalancers listeners list --loadbalancer lb-123",
}

func init() {
	LoadbalancersCmd.AddCommand(listeners.ListenersCmd)
}
