package quotas

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientquotas "github.com/thalassa-cloud/client-go/quotas"
)

var (
	requestIncreaseNewMax int64
	requestIncreaseReason string
)

var requestIncreaseCmd = &cobra.Command{
	Use:               "request-increase <name>",
	Aliases:           []string{"increase", "request"},
	Short:             "Request a quota limit increase",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteOrganisationQuotaName,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if requestIncreaseNewMax <= 0 {
			return fmt.Errorf("--new-max-usage must be greater than 0")
		}
		if requestIncreaseReason == "" {
			return fmt.Errorf("--reason is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if err := client.Quotas().RequestQuotaIncrease(ctx, clientquotas.RequestQuotaIncreaseRequest{
			Name:        args[0],
			NewMaxUsage: requestIncreaseNewMax,
			Reason:      requestIncreaseReason,
		}); err != nil {
			return fmt.Errorf("failed to request quota increase: %w", err)
		}

		fmt.Printf("Quota increase requested for %s (new limit: %d)\n", args[0], requestIncreaseNewMax)
		return nil
	},
}

func init() {
	QuotasCmd.AddCommand(requestIncreaseCmd)
	requestIncreaseCmd.Flags().Int64Var(&requestIncreaseNewMax, "new-max-usage", 0, "Requested maximum usage limit")
	requestIncreaseCmd.Flags().StringVar(&requestIncreaseReason, "reason", "", "Reason for the increase request")
	_ = requestIncreaseCmd.MarkFlagRequired("new-max-usage")
	_ = requestIncreaseCmd.MarkFlagRequired("reason")
}
