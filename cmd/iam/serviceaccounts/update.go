package serviceaccounts

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	updateName        string
	updateDescription string
	updateLabels      []string
	updateAnnotations []string
)

var updateCmd = &cobra.Command{
	Use:               "update <identity>",
	Short:             "Update a service account (only set flags are sent)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMServiceAccountIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		req := clientiam.UpdateServiceAccountRequest{}
		if cmd.Flags().Changed("name") {
			req.Name = &updateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = &updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = shared.KeyValuePairsToMap(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = shared.KeyValuePairsToMap(updateAnnotations)
		}
		sa, err := client.IAM().UpdateServiceAccount(ctx, args[0], req)
		if err != nil {
			return fmt.Errorf("failed to update service account: %w", err)
		}
		if sa == nil {
			return nil
		}
		body := [][]string{{sa.Identity, sa.Name, sa.Slug}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug"}, body)
		}
		return nil
	},
}

func init() {
	ServiceAccountsCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	updateCmd.Flags().StringVar(&updateName, "name", "", "Name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description (empty to clear)")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Replace labels (key=value, repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Replace annotations (key=value, repeatable)")
}
