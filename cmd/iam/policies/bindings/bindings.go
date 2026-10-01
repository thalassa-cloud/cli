package bindings

import "github.com/spf13/cobra"

// BindingsCmd manages IAM policy bindings.
var BindingsCmd = &cobra.Command{
	Use:   "bindings",
	Short: "Policy bindings (who receives the policy)",
}
