package routetables

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	routeUpdateDestination string
	routeUpdateGateway     string
	routeUpdateNatGateway  string
	routeUpdatePeering     string
	routeUpdateAddress     string
)

var routesUpdateCmd = &cobra.Command{
	Use:               "update ROUTE_TABLE ROUTE",
	Short:             "Update a route in a route table",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeRouteTableID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.IaaS().GetRouteTableRoute(cmd.Context(), args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed to get route: %w", err)
		}

		req := iaas.UpdateRouteTableRoute{
			DestinationCidrBlock: current.DestinationCidrBlock,
		}
		if current.TargetGatewayIdentity != nil {
			req.TargetGatewayIdentity = *current.TargetGatewayIdentity
		}
		if current.TargetNatGatewayIdentity != nil {
			req.TargetNatGatewayIdentity = *current.TargetNatGatewayIdentity
		}
		if current.TargetVpcPeeringConnectionId != nil {
			req.TargetVpcPeeringConnectionId = current.TargetVpcPeeringConnectionId
		}
		if current.GatewayAddress != nil {
			req.GatewayAddress = *current.GatewayAddress
		}

		if cmd.Flags().Changed("destination") {
			req.DestinationCidrBlock = routeUpdateDestination
		}
		if cmd.Flags().Changed("gateway") {
			req.TargetGatewayIdentity = routeUpdateGateway
			req.TargetNatGatewayIdentity = ""
			req.TargetVpcPeeringConnectionId = nil
			req.GatewayAddress = ""
		}
		if cmd.Flags().Changed("nat-gateway") {
			req.TargetNatGatewayIdentity = routeUpdateNatGateway
			req.TargetGatewayIdentity = ""
			req.TargetVpcPeeringConnectionId = nil
			req.GatewayAddress = ""
		}
		if cmd.Flags().Changed("vpc-peering") {
			req.TargetVpcPeeringConnectionId = &routeUpdatePeering
			req.TargetGatewayIdentity = ""
			req.TargetNatGatewayIdentity = ""
			req.GatewayAddress = ""
		}
		if cmd.Flags().Changed("gateway-address") {
			req.GatewayAddress = routeUpdateAddress
			req.TargetGatewayIdentity = ""
			req.TargetNatGatewayIdentity = ""
			req.TargetVpcPeeringConnectionId = nil
		}

		route, err := client.IaaS().UpdateRouteTableRoute(cmd.Context(), args[0], args[1], req)
		if err != nil {
			return err
		}

		body := [][]string{{
			route.Identity,
			route.DestinationCidrBlock,
			routeTarget(*route),
			route.Type,
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Destination", "Target", "Type"}, body)
		}
		return nil
	},
}

func init() {
	RoutesCmd.AddCommand(routesUpdateCmd)

	routesUpdateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	routesUpdateCmd.Flags().StringVar(&routeUpdateDestination, "destination", "", "Destination CIDR block")
	routesUpdateCmd.Flags().StringVar(&routeUpdateGateway, "gateway", "", "Target gateway identity")
	routesUpdateCmd.Flags().StringVar(&routeUpdateNatGateway, "nat-gateway", "", "Target NAT gateway identity")
	routesUpdateCmd.Flags().StringVar(&routeUpdatePeering, "vpc-peering", "", "Target VPC peering connection identity")
	routesUpdateCmd.Flags().StringVar(&routeUpdateAddress, "gateway-address", "", "Gateway address")

	_ = routesUpdateCmd.RegisterFlagCompletionFunc("nat-gateway", completeNatGatewayID)
}
