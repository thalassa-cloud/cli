package roles

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	noHeader               bool
	showExactTime          bool
	rolesListLabelSelector string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List organisation roles",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		req := &clientiam.ListOrganisationRolesRequest{}
		if rolesListLabelSelector != "" {
			req.Filters = []filters.Filter{
				&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(rolesListLabelSelector)},
			}
		}
		roles, err := client.IAM().ListOrganisationRoles(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list roles: %w", err)
		}
		body := make([][]string, 0, len(roles))
		for _, r := range roles {
			sys := "no"
			if r.System {
				sys = "yes"
			}
			ro := "no"
			if r.IsReadOnly {
				ro = "yes"
			}
			body = append(body, []string{
				r.Identity,
				r.Name,
				r.Slug,
				sys,
				ro,
				fmt.Sprintf("%d", len(r.Rules)),
				fmt.Sprintf("%d", len(r.Bindings)),
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "System", "Read-only", "Rules", "Bindings"}, body)
		}
		return nil
	},
}

func init() {
	RolesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	listCmd.Flags().StringVar(&rolesListLabelSelector, "label-selector", "", "Filter by labels (key=value,key2=value2)")
}
