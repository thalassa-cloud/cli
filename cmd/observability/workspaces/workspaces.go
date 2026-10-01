package workspaces

import (
	"fmt"

	"github.com/spf13/cobra"

	clientobs "github.com/thalassa-cloud/client-go/observability"
)

// WorkspacesCmd manages observability workspaces.
var WorkspacesCmd = &cobra.Command{
	Use:     "workspaces",
	Aliases: []string{"workspace", "ws"},
	Short:   "Manage observability workspaces",
	Long: `Create and manage unified Prometheus + Loki workspaces.

Each workspace exposes remote-write / push URLs for ingest and query URLs for
Prometheus and Loki in the current organisation/project context.`,
}

func regionName(ws clientobs.ObservabilityWorkspace) string {
	if ws.Region == nil {
		return "-"
	}
	if ws.Region.Name != "" {
		return ws.Region.Name
	}
	if ws.Region.Slug != "" {
		return ws.Region.Slug
	}
	return ws.Region.Identity
}

func retentionString(days *int) string {
	if days == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *days)
}

func boolYesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func printWorkspaceDetails(ws *clientobs.ObservabilityWorkspace) {
	fmt.Printf("ID:          %s\n", ws.Identity)
	fmt.Printf("Name:        %s\n", ws.Name)
	fmt.Printf("Description: %s\n", ws.Description)
	fmt.Printf("Status:      %s\n", ws.Status)
	if ws.StatusMessage != "" {
		fmt.Printf("Status msg:  %s\n", ws.StatusMessage)
	}
	fmt.Printf("Region:      %s\n", regionName(*ws))
	fmt.Printf("Retention:   %s days\n", retentionString(ws.RetentionDays))
	fmt.Printf("Prometheus:  %s\n", boolYesNo(ws.PrometheusEnabled))
	fmt.Printf("Loki:        %s\n", boolYesNo(ws.LokiEnabled))

	if ws.RemoteWriteURL != "" {
		fmt.Printf("Remote write:      %s\n", ws.RemoteWriteURL)
	}
	if ws.RemoteWriteOTLPURL != "" {
		fmt.Printf("Remote write OTLP: %s\n", ws.RemoteWriteOTLPURL)
	}
	if ws.PrometheusQueryURL != "" {
		fmt.Printf("Prometheus query:  %s\n", ws.PrometheusQueryURL)
	}
	if ws.AlertingURL != "" {
		fmt.Printf("Alerting:          %s\n", ws.AlertingURL)
	}
	if ws.PushURL != "" {
		fmt.Printf("Loki push:         %s\n", ws.PushURL)
	}
	if ws.PushOTLPURL != "" {
		fmt.Printf("Loki push OTLP:    %s\n", ws.PushOTLPURL)
	}
	if ws.LokiQueryURL != "" {
		fmt.Printf("Loki query:        %s\n", ws.LokiQueryURL)
	}
}
