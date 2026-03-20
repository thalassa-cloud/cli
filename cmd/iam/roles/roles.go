package roles

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/roles/bindings"
	"github.com/thalassa-cloud/cli/cmd/iam/roles/rules"
)

// RolesCmd represents organisation IAM roles.
var RolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Organisation IAM roles, permission rules, and bindings",
	Long: `Custom organisation roles define permission rules and can be bound to users, teams,
or service accounts. System roles may be read-only; the API enforces what you can change.`,
}

func init() {
	RolesCmd.AddCommand(rules.RulesCmd)
	RolesCmd.AddCommand(bindings.BindingsCmd)
}
