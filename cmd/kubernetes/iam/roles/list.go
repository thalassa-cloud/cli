package roles

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/kubernetes/iam/shared"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var (
	noHeader               bool
	rolesListLabelSelector string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List Kubernetes cluster IAM roles",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		req := &kubernetes.ListKubernetesClusterRolesRequest{}
		if rolesListLabelSelector != "" {
			req.Filters = []filters.Filter{
				&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(rolesListLabelSelector)},
			}
		}
		roles, err := client.Kubernetes().ListKubernetesClusterRoles(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list roles: %w", err)
		}
		body := make([][]string, 0, len(roles))
		for _, r := range roles {
			sys := "no"
			if r.System {
				sys = "yes"
			}
			body = append(body, []string{
				r.Identity,
				r.Name,
				r.Slug,
				sys,
				fmt.Sprintf("%d", len(r.Rules)),
				fmt.Sprintf("%d", len(r.Bindings)),
			})
		}
		if len(body) == 0 {
			fmt.Println("No Kubernetes cluster roles found")
			return nil
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "System", "Rules", "Bindings"}, body)
		}
		return nil
	},
}

func init() {
	RolesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().StringVar(&rolesListLabelSelector, "label-selector", "", "Filter by labels (key=value,key2=value2)")
}
