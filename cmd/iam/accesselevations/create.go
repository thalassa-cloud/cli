package accesselevations

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/iamresolve"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	createPolicy      string
	createReason      string
	createDuration    string
	createExpiresAt   string
	createLabels      []string
	createAnnotations []string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Request temporary elevation to an IAM policy",
	Long: `Create an access elevation request for yourself against a policy.

Exactly one of --duration or --expires-at is required.`,
	Example: `  tcloud iam access-elevations create --policy admin --reason "incident response" --duration 2h
  tcloud iam access-elevations create --policy pol-xxx --reason "break-glass" --expires-at 2026-10-01T20:00:00Z`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if createExpiresAt != "" && createDuration != "" {
			return fmt.Errorf("specify at most one of --expires-at or --duration")
		}
		if createExpiresAt == "" && createDuration == "" {
			return fmt.Errorf("specify one of --duration or --expires-at")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		policy, err := iamresolve.ResolveIamPolicyRef(cmd.Context(), client.IAM(), createPolicy)
		if err != nil {
			return err
		}

		create := clientiam.CreateAccessElevationRequest{
			PolicyIdentity: policy.Identity,
			Reason:         createReason,
			Duration:       createDuration,
			Labels:         shared.KeyValuePairsToMap(createLabels),
			Annotations:    shared.KeyValuePairsToMap(createAnnotations),
		}
		if create.Labels == nil {
			create.Labels = map[string]string{}
		}
		if create.Annotations == nil {
			create.Annotations = map[string]string{}
		}
		if createExpiresAt != "" {
			t, err := time.Parse(time.RFC3339, createExpiresAt)
			if err != nil {
				return fmt.Errorf("--expires-at must be RFC3339 (e.g. 2026-12-31T23:59:59Z): %w", err)
			}
			create.ExpiresAt = &t
		}

		req, err := client.IAM().CreateAccessElevation(cmd.Context(), create)
		if err != nil {
			return fmt.Errorf("failed to create access elevation: %w", err)
		}

		printAccessElevationDetails(req, false)
		return nil
	},
}

func init() {
	AccessElevationsCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createPolicy, "policy", "", "IAM policy identity, slug, or name (required)")
	createCmd.Flags().StringVar(&createReason, "reason", "", "Reason for the elevation request (required)")
	createCmd.Flags().StringVar(&createDuration, "duration", "", "Relative TTL (e.g. 2h, 24h)")
	createCmd.Flags().StringVar(&createExpiresAt, "expires-at", "", "Absolute expiry time (RFC3339)")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	_ = createCmd.MarkFlagRequired("policy")
	_ = createCmd.MarkFlagRequired("reason")
	_ = createCmd.RegisterFlagCompletionFunc("policy", completion.CompleteIAMPolicyIdentityFlag)
}
