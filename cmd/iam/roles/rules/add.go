package rules

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	noHeader               bool
	ruleResources          []string
	ruleResourceIdentities []string
	rulePermissions        []string
	ruleNote               string
)

var addCmd = &cobra.Command{
	Use:               "add <role>",
	Short:             "Add a permission rule to a role",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMOrganisationRoleIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if len(rulePermissions) == 0 {
			return fmt.Errorf("at least one --permission is required (create, read, update, delete, list, *)")
		}
		perms := make([]clientiam.PermissionType, 0, len(rulePermissions))
		for _, p := range rulePermissions {
			perms = append(perms, clientiam.PermissionType(p))
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		rule, err := client.IAM().AddRoleRule(ctx, args[0], clientiam.OrganisationRolePermissionRule{
			Resources:          ruleResources,
			ResourceIdentities: ruleResourceIdentities,
			Permissions:        perms,
			Note:               ruleNote,
		})
		if err != nil {
			return fmt.Errorf("failed to add rule: %w", err)
		}
		if rule == nil {
			return nil
		}
		body := [][]string{{rule.Identity, strings.Join(ruleResources, ","), strings.Join(rulePermissions, ",")}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Rule ID", "Resources", "Permissions"}, body)
		}
		return nil
	},
}

func init() {
	RulesCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	addCmd.Flags().StringSliceVar(&ruleResources, "resource", nil, "Resource type (repeatable)")
	addCmd.Flags().StringSliceVar(&ruleResourceIdentities, "resource-identity", nil, "Concrete resource identity (repeatable)")
	addCmd.Flags().StringSliceVar(&rulePermissions, "permission", nil, "Permission: create, read, update, delete, list, or * (repeatable)")
	addCmd.Flags().StringVar(&ruleNote, "note", "", "Human-readable note for the rule")
	_ = addCmd.RegisterFlagCompletionFunc("permission", completion.CompleteIAMPermissionType)
}
