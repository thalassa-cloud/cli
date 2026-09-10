package routetables

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	routeCreateDestination string
	routeCreateGateway     string
	routeCreateNatGateway  string
	routeCreatePeering     string
	routeCreateAddress     string
)

var routesCreateCmd = &cobra.Command{
	Use:               "create ROUTE_TABLE",
	Short:             "Create a route in a route table",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeRouteTableID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if routeCreateDestination == "" {
			return fmt.Errorf("destination is required")
		}

		targets := 0
		if routeCreateGateway != "" {
			targets++
		}
		if routeCreateNatGateway != "" {
			targets++
		}
		if routeCreatePeering != "" {
			targets++
		}
		if routeCreateAddress != "" {
			targets++
		}
		if targets != 1 {
			return fmt.Errorf("exactly one of --gateway, --nat-gateway, --vpc-peering, or --gateway-address must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req := iaas.CreateRouteTableRoute{
			DestinationCidrBlock:     routeCreateDestination,
			TargetGatewayIdentity:    routeCreateGateway,
			TargetNatGatewayIdentity: routeCreateNatGateway,
			GatewayAddress:           routeCreateAddress,
		}
		if routeCreatePeering != "" {
			req.TargetVpcPeeringConnectionId = &routeCreatePeering
		}

		route, err := client.IaaS().CreateRouteTableRoute(cmd.Context(), args[0], req)
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
	RoutesCmd.AddCommand(routesCreateCmd)

	routesCreateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	routesCreateCmd.Flags().StringVar(&routeCreateDestination, "destination", "", "Destination CIDR block")
	routesCreateCmd.Flags().StringVar(&routeCreateGateway, "gateway", "", "Target gateway identity")
	routesCreateCmd.Flags().StringVar(&routeCreateNatGateway, "nat-gateway", "", "Target NAT gateway identity")
	routesCreateCmd.Flags().StringVar(&routeCreatePeering, "vpc-peering", "", "Target VPC peering connection identity")
	routesCreateCmd.Flags().StringVar(&routeCreateAddress, "gateway-address", "", "Gateway address")

	_ = routesCreateCmd.MarkFlagRequired("destination")
	_ = routesCreateCmd.RegisterFlagCompletionFunc("nat-gateway", completeNatGatewayID)
}
