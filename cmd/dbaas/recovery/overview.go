package recovery

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	dbaasutil "github.com/thalassa-cloud/cli/internal/dbaas"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
)

var (
	overviewBackupStore        string
	overviewBackupID           string
	overviewIncludeUnavailable bool
	overviewNoHeader           bool
	overviewShowExactTime      bool
	overviewExcludeUnavailable bool
)

var overviewCmd = &cobra.Command{
	Use:   "overview [cluster]",
	Short: "Show backup recovery and PITR overview",
	Long: `Show aggregated recovery information for a database cluster or backup store.

The overview includes backup counts, the approximate PITR window, WAL chain coverage
per backup, and any gaps detected in the archive.`,
	Example: "tcloud dbaas recovery overview dbc-123\ntcloud dbaas recovery overview --backup-store bst-123",
	Args:    cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return completion.CompleteDbClusterID(cmd, args, toComplete)
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusterIdentity := ""
		if len(args) > 0 {
			clusterIdentity = args[0]
		}

		ref, err := dbaasutil.ResolveBackupStore(cmd.Context(), client.DBaaS(), clusterIdentity, overviewBackupStore)
		if err != nil {
			return err
		}

		req := &dbaas.GetDbObjectStoreRecoveryOverviewRequest{
			ClusterID: clusterIdentity,
			BackupID:  overviewBackupID,
		}
		if cmd.Flags().Changed("include-unavailable") {
			include := overviewIncludeUnavailable
			req.IncludeUnavailable = &include
		}
		if overviewExcludeUnavailable {
			include := false
			req.IncludeUnavailable = &include
		}

		overview, err := client.DBaaS().GetDbObjectStoreRecoveryOverview(cmd.Context(), ref.StoreIdentity, req)
		if err != nil {
			return fmt.Errorf("failed to get recovery overview: %w", err)
		}

		printTruncationWarning(overview.Complete, overview.Truncation)
		printOverviewSummary(overview, overviewShowExactTime, overviewNoHeader)

		for _, cluster := range overview.Clusters {
			printClusterOverview(cluster, overviewShowExactTime, overviewNoHeader)
		}

		return nil
	},
}

func printTruncationWarning(complete bool, truncation dbaas.DbObjectStoreRecoveryTruncation) {
	if complete && !truncation.Truncated {
		return
	}
	reason := "recovery analysis is incomplete"
	if truncation.Reason != nil && strings.TrimSpace(*truncation.Reason) != "" {
		reason = strings.TrimSpace(*truncation.Reason)
	}
	if truncation.MaxWalObjectsConsidered != nil {
		fmt.Fprintf(os.Stderr, "Warning: %s (max WAL objects considered: %d). Do not treat the PITR window as complete.\n", reason, *truncation.MaxWalObjectsConsidered)
		return
	}
	fmt.Fprintf(os.Stderr, "Warning: %s. Do not treat the PITR window as complete.\n", reason)
}

func printOverviewSummary(overview *dbaas.DbObjectStoreRecoveryOverview, showExactTime, noHeader bool) {
	summary := overview.Summary
	body := [][]string{
		{"Backup Store", overview.ObjectStoreIdentity},
		{"Analysed At", formattime.FormatTime(overview.AnalysedAt.Local(), showExactTime)},
		{"Complete", dbaasutil.FormatGoodYesNo(overview.Complete)},
		{"Completed Backups", fmt.Sprintf("%d", summary.CompletedBackupCount)},
		{"Unavailable Backups", fmt.Sprintf("%d", summary.UnavailableBackupCount)},
		{"Failed Backups", fmt.Sprintf("%d", summary.FailedBackupCount)},
		{"Oldest Backup", dbaasutil.FormatOptionalTime(summary.OldestBackupAt, showExactTime)},
		{"Newest Backup", dbaasutil.FormatOptionalTime(summary.NewestBackupAt, showExactTime)},
		{"Earliest PITR From", dbaasutil.FormatOptionalTime(summary.EarliestPitrFrom, showExactTime)},
		{"Latest PITR Through", dbaasutil.FormatOptionalTime(summary.LatestPitrThroughApprox, showExactTime)},
		{"Latest Archived WAL", emptyDash(summary.LatestArchivedWal)},
		{"Latest WAL Archived At", dbaasutil.FormatOptionalTime(summary.LatestWalArchivedAt, showExactTime)},
		{"Timelines", emptyDash(strings.Join(summary.Timelines, ", "))},
		{"Gap Count", fmt.Sprintf("%d", summary.GapCount)},
	}

	fmt.Println("Recovery summary")
	if noHeader {
		table.Print(nil, body)
	} else {
		table.Print([]string{"Field", "Value"}, body)
	}
}

func printClusterOverview(cluster dbaas.DbObjectStoreRecoveryClusterOverview, showExactTime, noHeader bool) {
	exists := dbaasutil.FormatGoodYesNo(cluster.ClusterExists)
	clusterName := cluster.ClusterName
	if clusterName == "" {
		clusterName = "-"
	}

	body := [][]string{
		{"Cluster ID", emptyDash(cluster.ClusterIdentity)},
		{"Cluster Name", clusterName},
		{"Cluster Exists", exists},
		{"Latest Archived WAL", emptyDash(cluster.LatestArchivedWal)},
		{"Latest WAL Archived At", dbaasutil.FormatOptionalTime(cluster.LatestWalArchivedAt, showExactTime)},
		{"Timelines", emptyDash(strings.Join(cluster.Timelines, ", "))},
		{"Gap Count", fmt.Sprintf("%d", cluster.GapCount)},
	}
	if cluster.RecoveryWindow != nil {
		body = append(body,
			[]string{"First Recoverability", dbaasutil.FormatOptionalTime(cluster.RecoveryWindow.FirstRecoverabilityPoint, showExactTime)},
			[]string{"Last Successful Backup", dbaasutil.FormatOptionalTime(cluster.RecoveryWindow.LastSuccessfulBackupTime, showExactTime)},
			[]string{"Last Failed Backup", dbaasutil.FormatOptionalTime(cluster.RecoveryWindow.LastFailedBackupTime, showExactTime)},
		)
	}

	fmt.Printf("\nCluster %s\n", emptyDash(cluster.ClusterIdentity))
	if noHeader {
		table.Print(nil, body)
	} else {
		table.Print([]string{"Field", "Value"}, body)
	}

	if len(cluster.Backups) == 0 {
		fmt.Println("No backups in recovery overview")
		return
	}

	backupBody := make([][]string, 0, len(cluster.Backups))
	for _, backup := range cluster.Backups {
		firstMissing := "-"
		if backup.Coverage.FirstMissingWal != nil && *backup.Coverage.FirstMissingWal != "" {
			firstMissing = *backup.Coverage.FirstMissingWal
		}
		backupBody = append(backupBody, []string{
			backup.BackupIdentity,
			dbaasutil.FormatStatus(emptyDash(backup.Status)),
			dbaasutil.FormatBytes(backup.SizeBytes),
			emptyDash(backup.BeginWal),
			emptyDash(backup.EndWal),
			emptyDash(backup.BeginLsn),
			emptyDash(backup.EndLsn),
			dbaasutil.FormatStatus(emptyDash(string(backup.Coverage.ChainStatus))),
			dbaasutil.FormatOptionalTime(backup.Coverage.PitrThroughApproxAt, showExactTime),
			firstMissing,
		})
	}

	fmt.Println("Backup coverage")
	if noHeader {
		table.Print(nil, backupBody)
	} else {
		table.Print([]string{"Backup", "Status", "Size", "Begin WAL", "End WAL", "Begin LSN", "End LSN", "Chain", "PITR Through", "First Missing WAL"}, backupBody)
	}
}

func emptyDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func init() {
	RecoveryCmd.AddCommand(overviewCmd)

	overviewCmd.Flags().StringVar(&overviewBackupStore, "backup-store", "", "Backup store identity (optional when a cluster is provided)")
	overviewCmd.Flags().StringVar(&overviewBackupID, "backup", "", "Limit the overview to a single backup identity")
	overviewCmd.Flags().BoolVar(&overviewIncludeUnavailable, "include-unavailable", true, "Include unavailable backups in the overview")
	overviewCmd.Flags().BoolVar(&overviewExcludeUnavailable, "exclude-unavailable", false, "Exclude unavailable backups from the overview")
	overviewCmd.Flags().BoolVar(&overviewNoHeader, NoHeaderKey, false, "Do not print the header")
	overviewCmd.Flags().BoolVar(&overviewShowExactTime, "exact-time", false, "Show exact time instead of relative time")

	_ = overviewCmd.RegisterFlagCompletionFunc("backup-store", completion.CompleteDbBackupStoreID)
	_ = overviewCmd.RegisterFlagCompletionFunc("backup", completion.CompleteDbBackupID)
}
