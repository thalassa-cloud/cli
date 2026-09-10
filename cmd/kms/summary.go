package kms

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show KMS availability and regional key counts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		summary, err := client.KMS().GetSummary(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get KMS summary: %w", err)
		}

		fmt.Printf("Feature enabled: %v\n", summary.FeatureEnabled)
		body := make([][]string, 0, len(summary.Regions))
		for _, r := range summary.Regions {
			body = append(body, []string{
				r.Identity,
				r.Name,
				r.Slug,
				fmt.Sprintf("%v", r.KmsAvailable),
				fmt.Sprintf("%d", r.TotalKeys),
				fmt.Sprintf("%d", r.ActiveKeys),
				fmt.Sprintf("%d", r.DisabledKeys),
				fmt.Sprintf("%d", r.PendingDeletionKeys),
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "Available", "Total", "Active", "Disabled", "Pending Deletion"}, body)
		}
		return nil
	},
}

func init() {
	KmsCmd.AddCommand(summaryCmd)
	summaryCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
