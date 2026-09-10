package reservedips

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	associateLoadbalancer string
	associateNatGateway   string
)

var associateCmd = &cobra.Command{
	Use:     "associate",
	Short:   "Associate a reserved IP with a resource",
	Long:    "Associate a reserved IP with exactly one of a load balancer or NAT gateway.",
	Example: "tcloud networking reserved-ips associate rip-123 --loadbalancer lb-456\ntcloud networking reserved-ips associate rip-123 --nat-gateway ngw-789",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hasLB := associateLoadbalancer != ""
		hasNGW := associateNatGateway != ""
		if hasLB == hasNGW {
			return fmt.Errorf("exactly one of --loadbalancer or --nat-gateway must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req := iaas.AssociateReservedIpRequest{}
		if hasLB {
			req.LoadbalancerIdentity = &associateLoadbalancer
		}
		if hasNGW {
			req.NatGatewayIdentity = &associateNatGateway
		}

		rip, err := client.IaaS().AssociateReservedIP(cmd.Context(), args[0], req)
		if err != nil {
			return err
		}

		body := [][]string{{
			rip.Identity,
			rip.Name,
			regionName(*rip),
			string(rip.Status),
			dashIfEmpty(rip.IPv4Address),
			dashIfEmpty(rip.IPv6Address),
			attachedTo(*rip),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Region", "Status", "IPv4", "IPv6", "AttachedTo"}, body)
		}
		return nil
	},
}

func init() {
	ReservedIPsCmd.AddCommand(associateCmd)

	associateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	associateCmd.Flags().StringVar(&associateLoadbalancer, "loadbalancer", "", "Load balancer identity to associate with")
	associateCmd.Flags().StringVar(&associateNatGateway, "nat-gateway", "", "NAT gateway identity to associate with")

	associateCmd.ValidArgsFunction = completeReservedIPID
	_ = associateCmd.RegisterFlagCompletionFunc("loadbalancer", completeLoadbalancerID)
	_ = associateCmd.RegisterFlagCompletionFunc("nat-gateway", completeNatGatewayID)
}
