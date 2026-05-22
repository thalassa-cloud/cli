package quotas

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientquotas "github.com/thalassa-cloud/client-go/quotas"
)

const NoHeaderKey = "no-header"

var (
	noHeader              bool
	showIncreaseRequests  bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "get", "g"},
	Short:   "List organisation quotas",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		quotas, err := client.Quotas().ListOrganisationQuotas(ctx)
		if err != nil {
			return fmt.Errorf("failed to list quotas: %w", err)
		}
		if len(quotas) == 0 {
			fmt.Println("No quotas found")
			return nil
		}

		body := make([][]string, 0, len(quotas))
		for _, q := range quotas {
			service := "-"
			if q.Service != nil && *q.Service != "" {
				service = *q.Service
			}
			row := []string{
				q.Name,
				service,
				q.QuotaType,
				fmt.Sprintf("%d", q.CurrentUsage),
				fmt.Sprintf("%d", q.MaxUsage),
				formatUsagePercent(q.CurrentUsage, q.MaxUsage),
			}
			if showIncreaseRequests {
				row = append(row, formatIncreaseRequests(q.IncreaseRequests))
			}
			body = append(body, row)
		}

		headers := []string{"Name", "Service", "Type", "Used", "Limit", "Usage %"}
		if showIncreaseRequests {
			headers = append(headers, "Requested increases")
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print(headers, body)
		}
		return nil
	},
}

func formatUsagePercent(used, limit int64) string {
	if limit <= 0 {
		return "-"
	}
	return fmt.Sprintf("%.0f%%", float64(used)/float64(limit)*100)
}

func formatIncreaseRequests(requests []clientquotas.IncreaseOrganisationQuotaRequest) string {
	if len(requests) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(requests))
	for _, r := range requests {
		decision := r.Decision
		if decision == "" {
			decision = "pending"
		}
		parts = append(parts, fmt.Sprintf("%d (%s)", r.NewMaxUsageRequested, decision))
	}
	return strings.Join(parts, ", ")
}

func init() {
	QuotasCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showIncreaseRequests, "show-increase-requests", false, "Include requested increase limits and their decision status")
}
