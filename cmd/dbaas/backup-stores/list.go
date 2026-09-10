package backupstores

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	dbaasutil "github.com/thalassa-cloud/cli/internal/dbaas"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
)

var (
	backupStoreListNoHeader      bool
	backupStoreListShowExactTime bool
	backupStoreListShowLabels    bool
)

var backupStoreListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List database backup stores",
	Long:    "List backup stores in the organisation",
	Aliases: []string{"ls", "get"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		stores, err := client.DBaaS().ListDbObjectStores(cmd.Context(), &dbaas.ListDbObjectStoresRequest{})
		if err != nil {
			return fmt.Errorf("failed to list backup stores: %w", err)
		}
		if len(stores) == 0 {
			fmt.Println("No backup stores found")
			return nil
		}

		body := make([][]string, 0, len(stores))
		for _, store := range stores {
			regionName := ""
			if store.Region != nil {
				regionName = store.Region.Name
			}

			row := []string{
				store.Identity,
				store.Name,
				regionName,
				string(store.Status),
				store.RetentionPolicy,
				dbaasutil.FormatRetentionMode(store.RetentionMode),
				formattime.FormatTime(store.CreatedAt.Local(), backupStoreListShowExactTime),
			}

			if backupStoreListShowLabels {
				labelStrs := make([]string, 0, len(store.Labels))
				for k, v := range store.Labels {
					labelStrs = append(labelStrs, k+"="+v)
				}
				sort.Strings(labelStrs)
				if len(labelStrs) == 0 {
					labelStrs = []string{"-"}
				}
				row = append(row, strings.Join(labelStrs, ","))
			}

			body = append(body, row)
		}

		if backupStoreListNoHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Region", "Status", "Retention", "Retention Mode", "Created"}
			if backupStoreListShowLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}

		return nil
	},
}

func init() {
	BackupStoresCmd.AddCommand(backupStoreListCmd)

	backupStoreListCmd.Flags().BoolVar(&backupStoreListNoHeader, NoHeaderKey, false, "Do not print the header")
	backupStoreListCmd.Flags().BoolVar(&backupStoreListShowExactTime, "exact-time", false, "Show exact time instead of relative time")
	backupStoreListCmd.Flags().BoolVar(&backupStoreListShowLabels, "show-labels", false, "Show labels")
}
