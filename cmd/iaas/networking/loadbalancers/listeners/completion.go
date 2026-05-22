package listeners

import "github.com/thalassa-cloud/cli/internal/completion"

var (
	completeLoadbalancerID         = completion.CompleteLoadbalancerID
	completeLoadbalancerListenerID = completion.CompleteLoadbalancerListenerID
	completeTargetGroupID          = completion.CompleteTargetGroupID
	completeLoadbalancerProtocol   = completion.CompleteLoadbalancerProtocol
	completeOutputFormat           = completion.CompleteOutputFormat
)
