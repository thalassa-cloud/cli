package roles

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/kubernetes/iam/roles/bindings"
	"github.com/thalassa-cloud/cli/cmd/kubernetes/iam/roles/rules"
)

// RolesCmd represents Kubernetes cluster IAM roles.
var RolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Kubernetes cluster IAM roles, permission rules, and bindings",
	Long: `Custom Kubernetes cluster roles define RBAC-style permission rules and can be bound to users,
teams, or service accounts. System roles may be read-only; the API enforces what you can change.`,
}

func init() {
	RolesCmd.AddCommand(rules.RulesCmd)
	RolesCmd.AddCommand(bindings.BindingsCmd)
}
