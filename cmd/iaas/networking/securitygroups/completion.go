package securitygroups

import (
	"github.com/thalassa-cloud/cli/internal/completion"
)

// Re-export completion functions for convenience
var (
	completeSecurityGroupID = completion.CompleteSecurityGroupID
	completeVPCID          = completion.CompleteVPCID
	completeOutputFormat   = completion.CompleteOutputFormat
)
