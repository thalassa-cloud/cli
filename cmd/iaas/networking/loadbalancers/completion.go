package loadbalancers

import "github.com/thalassa-cloud/cli/internal/completion"

var (
	completeLoadbalancerID   = completion.CompleteLoadbalancerID
	completeVPCID            = completion.CompleteVPCID
	completeRegion           = completion.CompleteRegion
	completeSubnetID         = completion.CompleteSubnetID
	completeSecurityGroupID  = completion.CompleteSecurityGroupID
	completeOutputFormat     = completion.CompleteOutputFormat
	completeTargetGroupID    = completion.CompleteTargetGroupID
)
