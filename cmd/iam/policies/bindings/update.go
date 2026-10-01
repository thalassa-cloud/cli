package bindings

import (
	"fmt"
	"strings"

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
	Use:               "update <policy> <binding>",
	Short:             "Update a policy binding (only set flags are changed)",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completion.CompleteIAMPolicyThenBinding,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		policy, err := iamresolve.ResolveIamPolicyRef(ctx, client.IAM(), args[0])
		if err != nil {
			return err
		}

		bindings, err := client.IAM().ListIamPolicyBindings(ctx, policy.Identity, &clientiam.ListIamPolicyBindingsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list bindings: %w", err)
		}
		var current *clientiam.IamPolicyBinding
		for i := range bindings {
			b := &bindings[i]
			if b.Identity == args[1] || strings.EqualFold(b.Slug, args[1]) || strings.EqualFold(b.Name, args[1]) {
				current = b
				break
			}
		}
		if current == nil {
			return fmt.Errorf("policy binding not found: %s", args[1])
		}

		req := clientiam.UpdateIamPolicyBindingRequest{
			Description: current.Description,
			Labels:      current.Labels,
			Annotations: current.Annotations,
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

		binding, err := client.IAM().UpdateIamPolicyBinding(ctx, policy.Identity, current.Identity, req)
		if err != nil {
			return fmt.Errorf("failed to update binding: %w", err)
		}
		body := [][]string{{binding.Identity, binding.Name, binding.Slug}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug"}, body)
		}
		return nil
	},
}

func init() {
	BindingsCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Binding description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Replace labels (key=value, repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Replace annotations (key=value, repeatable)")
}
