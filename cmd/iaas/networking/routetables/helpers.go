package routetables

import (
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	completeRouteTableID = completion.CompleteRouteTableID
	completeVPCID        = completion.CompleteVPCID
	completeOutputFormat = completion.CompleteOutputFormat
	completeNatGatewayID = completion.CompleteNatGatewayID
)

func descriptionString(desc *string) string {
	if desc == nil || *desc == "" {
		return "-"
	}
	return *desc
}

func vpcName(rt iaas.RouteTable) string {
	if rt.Vpc == nil {
		return ""
	}
	return rt.Vpc.Name
}

func routeTarget(entry iaas.RouteEntry) string {
	switch {
	case entry.TargetNatGatewayIdentity != nil && *entry.TargetNatGatewayIdentity != "":
		return "nat:" + *entry.TargetNatGatewayIdentity
	case entry.TargetGatewayIdentity != nil && *entry.TargetGatewayIdentity != "":
		return "gateway:" + *entry.TargetGatewayIdentity
	case entry.TargetVpcPeeringConnectionId != nil && *entry.TargetVpcPeeringConnectionId != "":
		return "peering:" + *entry.TargetVpcPeeringConnectionId
	case entry.GatewayAddress != nil && *entry.GatewayAddress != "":
		return "address:" + *entry.GatewayAddress
	default:
		return "-"
	}
}
