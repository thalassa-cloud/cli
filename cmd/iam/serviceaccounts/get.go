package serviceaccounts

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var getCmd = &cobra.Command{
	Use:               "get <identity>",
	Short:             "Show a service account",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMServiceAccountIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		sa, err := client.IAM().GetServiceAccount(ctx, args[0])
		if err != nil {
			return fmt.Errorf("failed to get service account: %w", err)
		}
		fmt.Printf("Identity:    %s\n", sa.Identity)
		fmt.Printf("Name:        %s\n", sa.Name)
		fmt.Printf("Slug:        %s\n", sa.Slug)
		if sa.Description != nil {
			fmt.Printf("Description: %s\n", *sa.Description)
		}
		fmt.Printf("Created:     %s\n", formattime.FormatTime(sa.CreatedAt.Local(), showExactTime))
		if len(sa.RoleBindings) > 0 {
			fmt.Println("\nRole bindings:")
			body := make([][]string, 0, len(sa.RoleBindings))
			for _, b := range sa.RoleBindings {
				roleName := ""
				if b.OrganisationRole != nil {
					roleName = b.OrganisationRole.Slug
				}
				body = append(body, []string{b.Identity, b.Name, roleName})
			}
			table.Print([]string{"ID", "Name", "Role"}, body)
		}
		return nil
	},
}

func init() {
	ServiceAccountsCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	getCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
