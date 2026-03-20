package invites

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	noHeader      bool
	showExactTime bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List pending organisation invites",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		invites, err := client.IAM().ListOrganisationMemberInvites(ctx, &clientiam.ListOrganisationMemberInvitesRequest{})
		if err != nil {
			return fmt.Errorf("failed to list invites: %w", err)
		}
		body := make([][]string, 0, len(invites))
		for _, inv := range invites {
			exp := ""
			if inv.ExpiresAt != nil {
				exp = formattime.FormatTime(inv.ExpiresAt.Local(), showExactTime)
			}
			body = append(body, []string{
				inv.Email,
				string(inv.Role),
				inv.InviteCode,
				formattime.FormatTime(inv.CreatedAt.Local(), showExactTime),
				exp,
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Email", "Role", "Code", "Created", "Expires"}, body)
		}
		return nil
	},
}

func init() {
	InvitesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
