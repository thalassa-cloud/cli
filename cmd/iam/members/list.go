package members

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
	Short:   "List organisation members",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		members, err := client.IAM().ListOrganisationMembers(ctx, &clientiam.ListMembersRequest{})
		if err != nil {
			return fmt.Errorf("failed to list members: %w", err)
		}
		body := make([][]string, 0, len(members))
		for _, m := range members {
			body = append(body, []string{
				m.Identity,
				string(m.MemberType),
				shared.UserPtrDisplay(m.User),
				formattime.FormatTime(m.CreatedAt.Local(), showExactTime),
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Role", "User", "Joined"}, body)
		}
		return nil
	},
}

func init() {
	MembersCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
