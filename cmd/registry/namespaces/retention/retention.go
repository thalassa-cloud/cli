package retention

import "github.com/spf13/cobra"

const NamespaceFlag = "namespace"

// RetentionCmd manages retention policy operations.
var RetentionCmd = &cobra.Command{
	Use:     "retention",
	Aliases: []string{"retention-policy"},
	Short:   "Run retention policies",
	Example: "tcloud registry namespaces retention run --namespace crns-123",
}
