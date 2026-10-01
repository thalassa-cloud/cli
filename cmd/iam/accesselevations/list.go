package accesselevations

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List access elevation requests",
	Long: `List access elevation requests visible to approvers in the current scope.

Use --mine to list only requests you created.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		var requests []clientiam.IamAccessElevationRequest
		if mineFlag {
			requests, err = client.IAM().ListMyAccessElevations(cmd.Context())
		} else {
			requests, err = client.IAM().ListAccessElevations(cmd.Context())
		}
		if err != nil {
			return fmt.Errorf("failed to list access elevations: %w", err)
		}

		body := make([][]string, 0, len(requests))
		for _, req := range requests {
			policy := ""
			if req.IamPolicy != nil {
				policy = req.IamPolicy.Name
				if policy == "" {
					policy = req.IamPolicy.Identity
				}
			}
			requester := shared.UserPtrDisplay(req.Requester)
			body = append(body, []string{
				req.Identity,
				string(req.Status),
				policy,
				requester,
				req.Reason,
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Status", "Policy", "Requester", "Reason"}, body)
		}
		return nil
	},
}

func init() {
	AccessElevationsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&mineFlag, "mine", false, "List only requests you created")
}
