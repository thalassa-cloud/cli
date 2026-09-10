package natgateways

import (
	"github.com/thalassa-cloud/cli/internal/completion"
)

// Re-export completion functions for convenience
var (
	completeNatGatewayID    = completion.CompleteNatGatewayID
	completeVPCID           = completion.CompleteVPCID
	completeRegion          = completion.CompleteRegion
	completeOutputFormat    = completion.CompleteOutputFormat
	completeSubnetID        = completion.CompleteSubnetID
	completeSecurityGroupID = completion.CompleteSecurityGroupID
	completeReservedIPID    = completion.CompleteReservedIPID
)
