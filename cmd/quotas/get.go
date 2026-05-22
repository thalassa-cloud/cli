package quotas

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var showExactTime bool

var getCmd = &cobra.Command{
	Use:               "get <name>",
	Short:             "Show details for a quota",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteOrganisationQuotaName,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		q, err := client.Quotas().GetOrganisationQuota(ctx, args[0])
		if err != nil {
			return fmt.Errorf("failed to get quota: %w", err)
		}

		service := "-"
		if q.Service != nil && *q.Service != "" {
			service = *q.Service
		}

		fmt.Printf("Name:         %s\n", q.Name)
		fmt.Printf("Description:  %s\n", q.Description)
		fmt.Printf("Service:      %s\n", service)
		fmt.Printf("Type:         %s\n", q.QuotaType)
		fmt.Printf("Used:         %d\n", q.CurrentUsage)
		fmt.Printf("Limit:        %d\n", q.MaxUsage)
		fmt.Printf("Usage:        %s\n", formatUsagePercent(q.CurrentUsage, q.MaxUsage))
		fmt.Printf("Created:      %s\n", formattime.FormatTime(q.CreatedAt.Local(), showExactTime))

		if len(q.IncreaseRequests) == 0 {
			return nil
		}

		fmt.Println("\nIncrease requests:")
		body := make([][]string, 0, len(q.IncreaseRequests))
		for _, r := range q.IncreaseRequests {
			reason := r.RequestedReasonMessage
			if r.DecisionReason != nil && *r.DecisionReason != "" {
				reason = *r.DecisionReason
			}
			body = append(body, []string{
				r.Identity,
				fmt.Sprintf("%d", r.NewMaxUsageRequested),
				r.Decision,
				formattime.FormatTime(r.CreatedAt.Local(), showExactTime),
				reason,
			})
		}
		table.Print([]string{"ID", "Requested limit", "Decision", "Created", "Note"}, body)
		return nil
	},
}

func init() {
	QuotasCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
