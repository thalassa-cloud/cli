package backupstores

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	dbaasutil "github.com/thalassa-cloud/cli/internal/dbaas"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupStoreViewShowExactTime bool
	backupStoreViewNoHeader      bool
)

var backupStoreViewCmd = &cobra.Command{
	Use:               "view",
	Short:             "View backup store details",
	Long:              "View detailed information about a database backup store, including recovery windows",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDbBackupStoreID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		storeIdentity := args[0]
		store, err := client.DBaaS().GetDbObjectStore(cmd.Context(), storeIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("backup store not found: %s", storeIdentity)
			}
			return fmt.Errorf("failed to get backup store: %w", err)
		}

		regionName := ""
		if store.Region != nil {
			regionName = store.Region.Name
		}

		body := [][]string{
			{"ID", store.Identity},
			{"Name", store.Name},
			{"Status", dbaasutil.FormatStatus(string(store.Status))},
			{"Region", regionName},
			{"Retention Policy", store.RetentionPolicy},
			{"Retention Mode", dbaasutil.FormatRetentionMode(store.RetentionMode)},
			{"Delete Protection", dbaasutil.FormatBoolYesNo(store.DeleteProtection)},
			{"Created", formattime.FormatTime(store.CreatedAt.Local(), backupStoreViewShowExactTime)},
		}

		if store.Description != "" {
			body = append(body, []string{"Description", store.Description})
		}
		if store.StatusMessage != "" {
			body = append(body, []string{"Status Message", store.StatusMessage})
		}
		if store.UpdatedAt != nil {
			body = append(body, []string{"Updated", formattime.FormatTime(store.UpdatedAt.Local(), backupStoreViewShowExactTime)})
		}
		if store.ObjectStorageBucket != nil {
			body = append(body, []string{"Bucket", store.ObjectStorageBucket.Name})
			body = append(body, []string{"Bucket ID", store.ObjectStorageBucket.Identity})
		}

		if len(store.Labels) > 0 {
			labelStrs := make([]string, 0, len(store.Labels))
			for k, v := range store.Labels {
				labelStrs = append(labelStrs, k+"="+v)
			}
			sort.Strings(labelStrs)
			body = append(body, []string{"Labels", strings.Join(labelStrs, ", ")})
		}
		if len(store.Annotations) > 0 {
			annotationStrs := make([]string, 0, len(store.Annotations))
			for k, v := range store.Annotations {
				annotationStrs = append(annotationStrs, k+"="+v)
			}
			sort.Strings(annotationStrs)
			body = append(body, []string{"Annotations", strings.Join(annotationStrs, ", ")})
		}

		if backupStoreViewNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Field", "Value"}, body)
		}

		printRecoveryWindows(store, backupStoreViewShowExactTime, backupStoreViewNoHeader)
		return nil
	},
}

func init() {
	BackupStoresCmd.AddCommand(backupStoreViewCmd)

	backupStoreViewCmd.Flags().BoolVar(&backupStoreViewNoHeader, NoHeaderKey, false, "Do not print the header")
	backupStoreViewCmd.Flags().BoolVar(&backupStoreViewShowExactTime, "exact-time", false, "Show exact time instead of relative time")
}

func printRecoveryWindows(store *dbaas.DbObjectStore, showExactTime, noHeader bool) {
	if len(store.ServerRecoveryWindow) == 0 {
		return
	}

	servers := make([]string, 0, len(store.ServerRecoveryWindow))
	for name := range store.ServerRecoveryWindow {
		servers = append(servers, name)
	}
	sort.Strings(servers)

	storeReady := store.Status == dbaas.ObjectStatusReady
	body := make([][]string, 0, len(servers))
	for _, name := range servers {
		window := store.ServerRecoveryWindow[name]
		server := name
		if server == "" {
			server = "-"
		}
		body = append(body, []string{
			server,
			dbaasutil.FormatOptionalTime(window.FirstRecoverabilityPoint, showExactTime),
			dbaasutil.FormatOptionalTime(window.LastSuccessfulBackupTime, showExactTime),
			dbaasutil.FormatOptionalTime(window.LastFailedBackupTime, showExactTime),
			dbaasutil.FormatGoodYesNo(window.PitrAvailable(storeReady, storeReady)),
		})
	}

	fmt.Println()
	if !noHeader {
		fmt.Println("Recovery windows")
		table.Print([]string{"Server", "First Recoverability", "Last Successful Backup", "Last Failed Backup", "PITR"}, body)
		return
	}
	table.Print(nil, body)
}
