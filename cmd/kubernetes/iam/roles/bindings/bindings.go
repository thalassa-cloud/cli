package bindings

import "github.com/spf13/cobra"

// BindingsCmd manages Kubernetes cluster role bindings.
var BindingsCmd = &cobra.Command{
	Use:   "bindings",
	Short: "Kubernetes cluster role bindings (who receives the role)",
}
