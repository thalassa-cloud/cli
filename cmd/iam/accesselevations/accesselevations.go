package accesselevations

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

// AccessElevationsCmd manages temporary IAM access elevation requests.
var AccessElevationsCmd = &cobra.Command{
	Use:     "access-elevations",
	Aliases: []string{"access-elevation", "elevations", "ae"},
	Short:   "Request and review temporary IAM access elevations",
	Long: `Request temporary elevation to an IAM policy, and approve, reject, or revoke
requests in the current organisation/project scope.

Use --mine on list/get to work with your own requests. create and cancel always
act on your requests; approve, reject, and revoke use the approver API.`,
}

var (
	noHeader      bool
	showExactTime bool
	mineFlag      bool
)

func printAccessElevationDetails(req *clientiam.IamAccessElevationRequest, exactTime bool) {
	fmt.Printf("Identity:    %s\n", req.Identity)
	fmt.Printf("Status:      %s\n", req.Status)
	fmt.Printf("Reason:      %s\n", req.Reason)
	fmt.Printf("Expires:     %s\n", formattime.FormatTime(req.RequestedExpiresAt.Local(), exactTime))
	fmt.Printf("Created:     %s\n", formattime.FormatTime(req.CreatedAt.Local(), exactTime))
	if req.UpdatedAt != nil {
		fmt.Printf("Updated:     %s\n", formattime.FormatTime(req.UpdatedAt.Local(), exactTime))
	}
	if req.Project != nil {
		fmt.Printf("Project:     %s (%s)\n", req.Project.Name, req.Project.Identity)
	} else if req.ProjectId != "" {
		fmt.Printf("Project:     %s\n", req.ProjectId)
	}
	if req.IamPolicy != nil {
		fmt.Printf("Policy:      %s (%s)\n", req.IamPolicy.Name, req.IamPolicy.Identity)
	}
	if req.Requester != nil {
		fmt.Printf("Requester:   %s\n", shared.UserPtrDisplay(req.Requester))
	}
	if req.ReviewedAt != nil {
		fmt.Printf("Reviewed:    %s\n", formattime.FormatTime(req.ReviewedAt.Local(), exactTime))
	}
	if req.ReviewedBy != nil {
		fmt.Printf("Reviewed by: %s\n", shared.UserPtrDisplay(req.ReviewedBy))
	}
	if req.ReviewNote != nil && *req.ReviewNote != "" {
		fmt.Printf("Review note: %s\n", *req.ReviewNote)
	}
	if req.IamPolicyBinding != nil {
		fmt.Printf("Binding:     %s (%s)\n", req.IamPolicyBinding.Name, req.IamPolicyBinding.Identity)
	}
}
