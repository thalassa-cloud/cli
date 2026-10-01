package accesselevations

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

type reviewAction struct {
	use       string
	short     string
	confirm   string
	errPrefix string
	call      func(ctx context.Context, api *clientiam.Client, identity string, review clientiam.ReviewAccessElevationRequest) (*clientiam.IamAccessElevationRequest, error)
}

func registerReviewCommand(a reviewAction) {
	var force bool
	var reviewNote string

	cmd := &cobra.Command{
		Use:               a.use + " <identity>",
		Short:             a.short,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.CompleteIAMAccessElevationIdentity,
		RunE: func(cmd *cobra.Command, args []string) error {
			ok, err := shared.PromptDestructiveUnlessForce(force, fmt.Sprintf(
				"%s\n  Request: %s\n", a.confirm, args[0]))
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}

			client, err := thalassaclient.GetThalassaClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			req, err := a.call(cmd.Context(), client.IAM(), args[0], clientiam.ReviewAccessElevationRequest{
				ReviewNote: reviewNote,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", a.errPrefix, err)
			}

			printAccessElevationDetails(req, false)
			return nil
		},
	}

	AccessElevationsCmd.AddCommand(cmd)
	cmd.Flags().BoolVar(&force, shared.ForceKey, false, "Skip the confirmation prompt")
	cmd.Flags().StringVar(&reviewNote, "review-note", "", "Optional note attached to the review")
}

func init() {
	registerReviewCommand(reviewAction{
		use:       "approve",
		short:     "Approve a pending access elevation request",
		confirm:   "Are you sure you want to approve this access elevation request?",
		errPrefix: "failed to approve access elevation",
		call: func(ctx context.Context, api *clientiam.Client, identity string, review clientiam.ReviewAccessElevationRequest) (*clientiam.IamAccessElevationRequest, error) {
			return api.ApproveAccessElevation(ctx, identity, review)
		},
	})
	registerReviewCommand(reviewAction{
		use:       "reject",
		short:     "Reject a pending access elevation request",
		confirm:   "Are you sure you want to reject this access elevation request?",
		errPrefix: "failed to reject access elevation",
		call: func(ctx context.Context, api *clientiam.Client, identity string, review clientiam.ReviewAccessElevationRequest) (*clientiam.IamAccessElevationRequest, error) {
			return api.RejectAccessElevation(ctx, identity, review)
		},
	})
	registerReviewCommand(reviewAction{
		use:       "revoke",
		short:     "Revoke an approved access elevation",
		confirm:   "Are you sure you want to revoke this access elevation?",
		errPrefix: "failed to revoke access elevation",
		call: func(ctx context.Context, api *clientiam.Client, identity string, review clientiam.ReviewAccessElevationRequest) (*clientiam.IamAccessElevationRequest, error) {
			return api.RevokeAccessElevation(ctx, identity, review)
		},
	})
}
