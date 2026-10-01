package policies

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/policies/bindings"
	"github.com/thalassa-cloud/cli/cmd/iam/policies/rules"
)

// PoliciesCmd manages project/organisation IAM policies.
var PoliciesCmd = &cobra.Command{
	Use:     "policies",
	Aliases: []string{"policy"},
	Short:   "IAM policies, permission rules, and bindings",
	Long: `IAM policies define permission rules and can be bound to users or service accounts
in the current organisation/project scope. System policies may be read-only; the API
enforces what you can change.`,
}

func init() {
	PoliciesCmd.AddCommand(rules.RulesCmd)
	PoliciesCmd.AddCommand(bindings.BindingsCmd)
}
