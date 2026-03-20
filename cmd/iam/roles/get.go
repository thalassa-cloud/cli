package roles

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var getCmd = &cobra.Command{
	Use:               "get <role>",
	Short:             "Show a role including rules and bindings summary",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMOrganisationRoleIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		role, err := client.IAM().GetOrganisationRole(ctx, args[0])
		if err != nil {
			return fmt.Errorf("failed to get role: %w", err)
		}
		fmt.Printf("Identity:    %s\n", role.Identity)
		fmt.Printf("Name:        %s\n", role.Name)
		fmt.Printf("Slug:        %s\n", role.Slug)
		fmt.Printf("Description: %s\n", role.Description)
		fmt.Printf("System:      %v\n", role.System)
		fmt.Printf("Read-only:   %v\n", role.IsReadOnly)
		fmt.Printf("Created:     %s\n", formattime.FormatTime(role.CreatedAt.Local(), showExactTime))

		if len(role.Rules) > 0 {
			fmt.Println("\nRules:")
			body := make([][]string, 0, len(role.Rules))
			for _, ru := range role.Rules {
				perms := make([]string, 0, len(ru.Permissions))
				for _, p := range ru.Permissions {
					perms = append(perms, string(p))
				}
				body = append(body, []string{
					ru.Identity,
					strings.Join(ru.Resources, ","),
					strings.Join(ru.ResourceIdentities, ","),
					strings.Join(perms, ","),
					ru.Note,
				})
			}
			table.Print([]string{"ID", "Resources", "Resource IDs", "Permissions", "Note"}, body)
		}
		if len(role.Bindings) > 0 {
			fmt.Println("\nBindings:")
			body := make([][]string, 0, len(role.Bindings))
			for _, b := range role.Bindings {
				subject := ""
				switch {
				case b.AppUser != nil:
					subject = "user:" + shared.UserPtrDisplay(b.AppUser)
				case b.OrganisationTeam != nil:
					subject = "team:" + b.OrganisationTeam.Slug
				case b.ServiceAccount != nil:
					subject = "service_account:" + b.ServiceAccount.Slug
				}
				body = append(body, []string{b.Identity, b.Name, subject})
			}
			table.Print([]string{"ID", "Name", "Subject"}, body)
		}
		return nil
	},
}

func init() {
	RolesCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	getCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
