package rules

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <role> <rule>",
	Short:             "Remove a permission rule from a role",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completion.CompleteIAMRoleThenRule,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to remove this permission rule from the role?\n  Role: %s\n  Rule: %s\n", args[0], args[1]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := client.IAM().DeleteRuleFromRole(ctx, args[0], args[1]); err != nil {
			return fmt.Errorf("failed to delete rule: %w", err)
		}
		fmt.Printf("Deleted rule %s from role %s\n", args[1], args[0])
		return nil
	},
}

func init() {
	RulesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
	deleteCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
