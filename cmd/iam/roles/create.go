package roles

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
	Short: "Create a custom organisation role",
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
		role, err := client.IAM().CreateOrganisationRole(ctx, clientiam.CreateOrganisationRoleRequest{
			Name:        createName,
			Description: createDescription,
			Labels:      shared.KeyValuePairsToMap(createLabels),
			Annotations: shared.KeyValuePairsToMap(createAnnotations),
		})
		if err != nil {
			return fmt.Errorf("failed to create role: %w", err)
		}
		body := [][]string{{role.Identity, role.Name, role.Slug}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug"}, body)
		}
		return nil
	},
}

func init() {
	RolesCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().StringVar(&createName, "name", "", "Role name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Role description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	_ = createCmd.MarkFlagRequired("name")
}
