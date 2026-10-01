package observability

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/observability/workspaces"
)

// ObservabilityCmd manages observability resources.
var ObservabilityCmd = &cobra.Command{
	Use:     "observability",
	Aliases: []string{"obs"},
	Short:   "Manage observability workspaces (metrics and logs)",
	Long: `Manage Thalassa observability workspaces that provide Prometheus (Cortex) metrics
and Loki log ingestion/query endpoints for the current organisation/project scope.`,
}

func init() {
	ObservabilityCmd.AddCommand(workspaces.WorkspacesCmd)
}
