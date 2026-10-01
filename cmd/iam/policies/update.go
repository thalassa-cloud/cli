package policies

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/iamresolve"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	updateDescription string
	updateLabels      []string
	updateAnnotations []string
)

var updateCmd = &cobra.Command{
	Use:               "update <policy>",
	Short:             "Update an IAM policy (only set flags are changed)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMPolicyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		current, err := iamresolve.ResolveIamPolicyRef(ctx, client.IAM(), args[0])
		if err != nil {
			return err
		}
		current, err = client.IAM().GetIamPolicy(ctx, current.Identity)
		if err != nil {
			return fmt.Errorf("failed to get IAM policy: %w", err)
		}

		req := clientiam.UpdateIamPolicyRequest{
			Description:  current.Description,
			Labels:       current.Labels,
			Annotations:  current.Annotations,
			Conditionals: current.Conditionals,
		}
		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = shared.KeyValuePairsToMap(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = shared.KeyValuePairsToMap(updateAnnotations)
		}
		if req.Labels == nil {
			req.Labels = map[string]string{}
		}
		if req.Annotations == nil {
			req.Annotations = map[string]string{}
		}

		policy, err := client.IAM().UpdateIamPolicy(ctx, current.Identity, req)
		if err != nil {
			return fmt.Errorf("failed to update IAM policy: %w", err)
		}
		body := [][]string{{policy.Identity, policy.Name, policy.Slug}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug"}, body)
		}
		return nil
	},
}

func init() {
	PoliciesCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Policy description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Replace labels (key=value, repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Replace annotations (key=value, repeatable)")
}
