package serviceaccounts

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	createName        string
	createDescription string
	createLabels      []string
	createAnnotations []string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a service account",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		if createName == "" {
			return fmt.Errorf("--name is required")
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		req := clientiam.CreateServiceAccountRequest{
			Name:        createName,
			Labels:      shared.KeyValuePairsToMap(createLabels),
			Annotations: shared.KeyValuePairsToMap(createAnnotations),
		}
		if createDescription != "" {
			req.Description = &createDescription
		}
		sa, err := client.IAM().CreateServiceAccount(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to create service account: %w", err)
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
	ServiceAccountsCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().StringVar(&createName, "name", "", "Service account name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	_ = createCmd.MarkFlagRequired("name")
}
