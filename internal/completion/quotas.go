package completion

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

// CompleteOrganisationQuotaName completes organisation quota names.
func CompleteOrganisationQuotaName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	quotas, err := client.Quotas().ListOrganisationQuotas(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	out := make([]string, 0, len(quotas))
	for _, q := range quotas {
		desc := q.Description
		if desc == "" {
			desc = fmt.Sprintf("%d / %d", q.CurrentUsage, q.MaxUsage)
		}
		out = append(out, q.Name+"\t"+desc)
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}
