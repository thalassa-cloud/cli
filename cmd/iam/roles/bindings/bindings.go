package bindings

import "github.com/spf13/cobra"

// BindingsCmd manages role bindings.
var BindingsCmd = &cobra.Command{
	Use:   "bindings",
	Short: "Role bindings (who receives the role)",
}
