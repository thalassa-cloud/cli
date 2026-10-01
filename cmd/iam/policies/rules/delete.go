package rules

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/iamresolve"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <policy> <rule>",
	Short:             "Remove a permission rule from a policy",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completion.CompleteIAMPolicyThenRule,
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
		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to remove this permission rule from the policy?\n  Policy: %s\n  Rule: %s\n", policy.Identity, args[1]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := client.IAM().DeleteIamPolicyRule(ctx, policy.Identity, args[1]); err != nil {
			return fmt.Errorf("failed to delete rule: %w", err)
		}
		fmt.Printf("Deleted rule %s from policy %s\n", args[1], policy.Identity)
		return nil
	},
}

func init() {
	RulesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
	deleteCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
