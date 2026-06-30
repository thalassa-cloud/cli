package rules

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/kubernetes/iam/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <role> <rule>",
	Short:             "Remove a permission rule from a Kubernetes cluster role",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completion.CompleteKubernetesClusterRoleThenRule,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		role, err := kuberesolve.ResolveKubernetesClusterRoleRef(ctx, client.Kubernetes(), args[0])
		if err != nil {
			return err
		}
		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to remove this permission rule from the role?\n  Role: %s\n  Rule: %s\n", role.Name, args[1]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := client.Kubernetes().DeleteClusterRoleRule(ctx, role.Identity, args[1]); err != nil {
			return fmt.Errorf("failed to delete rule: %w", err)
		}
		fmt.Printf("Deleted rule %s from role %s\n", args[1], role.Name)
		return nil
	},
}

func init() {
	RulesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
}
