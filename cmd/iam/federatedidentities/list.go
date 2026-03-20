package federatedidentities

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
	Short:   "List federated identities",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		req := &clientiam.ListFederatedIdentitiesRequest{}
		if listSelector != "" {
			req.Filters = []filters.Filter{
				&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(listSelector)},
			}
		}
		list, err := client.IAM().ListFederatedIdentities(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list federated identities: %w", err)
		}
		body := make([][]string, 0, len(list))
		for _, fi := range list {
			prov := ""
			if fi.Provider != nil {
				prov = fi.Provider.Name
			}
			body = append(body, []string{
				fi.Identity,
				fi.Name,
				fi.ProviderSubject,
				prov,
				string(fi.Status),
				formattime.FormatTime(fi.CreatedAt.Local(), showExactTime),
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Subject", "Provider", "Status", "Created"}, body)
		}
		return nil
	},
}

func init() {
	FederatedIdentitiesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	listCmd.Flags().StringVar(&listSelector, "label-selector", "", "Filter by labels (key=value,key2=value2)")
}
