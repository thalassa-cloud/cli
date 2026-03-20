package federatedidentityproviders

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	noHeader      bool
	showExactTime bool
	listSelector  string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List federated identity providers",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		req := &clientiam.ListFederatedIdentityProvidersRequest{}
		if listSelector != "" {
			req.Filters = []filters.Filter{
				&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(listSelector)},
			}
		}
		list, err := client.IAM().ListFederatedIdentityProviders(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list providers: %w", err)
		}
		body := make([][]string, 0, len(list))
		for _, p := range list {
			body = append(body, []string{
				p.Identity,
				p.Name,
				p.ProviderIssuer,
				string(p.Status),
				formattime.FormatTime(p.CreatedAt.Local(), showExactTime),
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Issuer", "Status", "Created"}, body)
		}
		return nil
	},
}

func init() {
	FederatedIdentityProvidersCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	listCmd.Flags().StringVar(&listSelector, "label-selector", "", "Filter by labels (key=value,key2=value2)")
}
