package bindings

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/iamresolve"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var noHeader bool

var listCmd = &cobra.Command{
	Use:               "list <role>",
	Short:             "List bindings for a role",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMOrganisationRoleIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		role, err := iamresolve.ResolveOrganisationRoleRef(ctx, client.IAM(), args[0])
		if err != nil {
			return err
		}

		bindings, err := client.IAM().ListRoleBindings(ctx, role.Identity, &clientiam.ListRoleBindingsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list bindings: %w", err)
		}
		body := make([][]string, 0, len(bindings))
		for _, b := range bindings {
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
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Subject"}, body)
		}
		return nil
	},
}

func init() {
	BindingsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
