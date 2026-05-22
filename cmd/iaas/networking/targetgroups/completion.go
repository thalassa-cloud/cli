package targetgroups

import "github.com/thalassa-cloud/cli/internal/completion"

var (
	completeTargetGroupID        = completion.CompleteTargetGroupID
	completeVPCID                = completion.CompleteVPCID
	completeLoadbalancerProtocol   = completion.CompleteLoadbalancerProtocol
	completeLoadbalancingPolicy    = completion.CompleteLoadbalancingPolicy
	completeMachineID            = completion.CompleteMachineID
	completeOutputFormat         = completion.CompleteOutputFormat
)
