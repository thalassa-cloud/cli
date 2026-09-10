package recovery

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	dbaasutil "github.com/thalassa-cloud/cli/internal/dbaas"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
)

var (
	walBackupStore   string
	walCluster       string
	walNoHeader      bool
	walShowExactTime bool
	walTimeline      string
	walKinds         []string
	walFlat          bool
	walLog           string
	walLimit         int
	walContinuation  string
	walQuery         string
)

var walCmd = &cobra.Command{
	Use:   "wal",
	Short: "Inspect WAL archives for recovery",
	Long:  "Browse WAL hierarchy, list WAL segments, and search WAL objects in a backup store",
}

var walHierarchyCmd = &cobra.Command{
	Use:   "hierarchy [cluster]",
	Short: "Show WAL archive hierarchy",
	Long:  "Show WAL timelines and log folders for a cluster in a backup store",
	Args:  cobra.MaximumNArgs(1),
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

		clusterIdentity, storeIdentity, err := resolveWalScope(cmd, args)
		if err != nil {
			return err
		}

		hierarchy, err := client.DBaaS().GetDbObjectStoreRecoveryWalHierarchy(cmd.Context(), storeIdentity, clusterIdentity, &dbaas.GetDbObjectStoreRecoveryWalHierarchyRequest{
			Timeline: walTimeline,
			Kinds:    walKinds,
		})
		if err != nil {
			return fmt.Errorf("failed to get WAL hierarchy: %w", err)
		}

		printTruncationWarning(hierarchy.Complete, hierarchy.Truncation)

		meta := [][]string{
			{"Cluster", hierarchy.ClusterIdentity},
			{"WAL Directory", emptyDash(hierarchy.WalDirectory)},
			{"Analysed At", formattime.FormatTime(hierarchy.AnalysedAt.Local(), walShowExactTime)},
			{"Complete", dbaasutil.FormatGoodYesNo(hierarchy.Complete)},
		}
		if walNoHeader {
			table.Print(nil, meta)
		} else {
			table.Print([]string{"Field", "Value"}, meta)
		}

		if len(hierarchy.Timelines) == 0 {
			fmt.Println("No WAL timelines found")
			return nil
		}

		body := make([][]string, 0)
		for _, timeline := range hierarchy.Timelines {
			for _, log := range timeline.Logs {
				body = append(body, []string{
					fmt.Sprintf("%d", timeline.Timeline),
					emptyDash(timeline.TimelineHex),
					fmt.Sprintf("%d", log.Log),
					emptyDash(log.LogHex),
					fmt.Sprintf("%d", log.ObjectCount),
					dbaasutil.FormatByteCount(log.SizeBytes),
					emptyDash(log.FirstWal),
					emptyDash(log.LastWal),
					fmt.Sprintf("%d", log.GapCount),
					dbaasutil.FormatOptionalTime(log.OldestArchivedAt, walShowExactTime),
					dbaasutil.FormatOptionalTime(log.NewestArchivedAt, walShowExactTime),
				})
			}
		}

		fmt.Println("WAL log folders")
		if walNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Timeline", "Timeline Hex", "Log", "Log Hex", "Objects", "Size", "First WAL", "Last WAL", "Gaps", "Oldest", "Newest"}, body)
		}

		return nil
	},
}

var walListCmd = &cobra.Command{
	Use:   "list [cluster]",
	Short: "List WAL segment objects",
	Long:  "List WAL segment objects for a timeline/log folder, or as a flat listing with --flat",
	Args:  cobra.MaximumNArgs(1),
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

		if !walFlat && (walTimeline == "" || walLog == "") {
			return fmt.Errorf("--timeline and --log are required unless --flat is set")
		}

		clusterIdentity, storeIdentity, err := resolveWalScope(cmd, args)
		if err != nil {
			return err
		}

		listReq := &dbaas.ListDbObjectStoreRecoveryWalSegmentsRequest{
			Timeline:          walTimeline,
			Log:               walLog,
			ContinuationToken: walContinuation,
			Limit:             walLimit,
			Kinds:             walKinds,
		}
		if walFlat {
			listReq.View = "flat"
		}

		segments, err := client.DBaaS().ListDbObjectStoreRecoveryWalSegments(cmd.Context(), storeIdentity, clusterIdentity, listReq)
		if err != nil {
			return fmt.Errorf("failed to list WAL segments: %w", err)
		}

		if len(segments.Objects) == 0 {
			fmt.Println("No WAL objects found")
			return nil
		}

		body := walObjectsBody(segments.Objects, walShowExactTime)
		if walNoHeader {
			table.Print(nil, body)
		} else {
			table.Print(walObjectHeaders(), body)
		}

		if segments.NextContinuationToken != nil && *segments.NextContinuationToken != "" {
			fmt.Printf("Next continuation token: %s\n", *segments.NextContinuationToken)
		}

		return nil
	},
}

var walSearchCmd = &cobra.Command{
	Use:   "search [cluster]",
	Short: "Search WAL archive objects",
	Long:  "Search WAL archive objects by WAL name or hex substring",
	Args:  cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return completion.CompleteDbClusterID(cmd, args, toComplete)
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if walQuery == "" {
			return fmt.Errorf("query is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusterIdentity, storeIdentity, err := resolveWalScope(cmd, args)
		if err != nil {
			return err
		}

		result, err := client.DBaaS().SearchDbObjectStoreRecoveryWal(cmd.Context(), storeIdentity, clusterIdentity, &dbaas.SearchDbObjectStoreRecoveryWalRequest{
			Query: walQuery,
			Limit: walLimit,
		})
		if err != nil {
			return fmt.Errorf("failed to search WAL archive: %w", err)
		}

		if len(result.Objects) == 0 {
			fmt.Println("No WAL objects matched the query")
			return nil
		}

		body := walObjectsBody(result.Objects, walShowExactTime)
		if walNoHeader {
			table.Print(nil, body)
		} else {
			table.Print(walObjectHeaders(), body)
		}

		return nil
	},
}

func resolveWalScope(cmd *cobra.Command, args []string) (clusterIdentity, storeIdentity string, err error) {
	clusterIdentity = walCluster
	if len(args) > 0 {
		clusterIdentity = args[0]
	}
	if clusterIdentity == "" {
		return "", "", fmt.Errorf("cluster identity is required")
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return "", "", fmt.Errorf("failed to create client: %w", err)
	}

	ref, err := dbaasutil.ResolveBackupStore(cmd.Context(), client.DBaaS(), clusterIdentity, walBackupStore)
	if err != nil {
		return "", "", err
	}
	return clusterIdentity, ref.StoreIdentity, nil
}

func walObjectHeaders() []string {
	return []string{"File", "Kind", "WAL", "Timeline", "Log", "Segment", "Size", "Archived"}
}

func walObjectsBody(objects []dbaas.DbObjectStoreRecoveryWalObject, showExactTime bool) [][]string {
	body := make([][]string, 0, len(objects))
	for _, object := range objects {
		name := object.FileName
		if name == "" {
			name = object.Key
		}
		body = append(body, []string{
			emptyDash(name),
			emptyDash(string(object.Kind)),
			emptyDash(object.WalName),
			emptyDash(object.TimelineHex),
			emptyDash(object.LogHex),
			emptyDash(object.SegmentHex),
			dbaasutil.FormatByteCount(object.SizeBytes),
			dbaasutil.FormatOptionalTime(object.ArchivedAt, showExactTime),
		})
	}
	return body
}

func registerWalFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&walBackupStore, "backup-store", "", "Backup store identity (optional when the cluster has a backup store)")
	cmd.Flags().StringVar(&walCluster, "cluster", "", "Cluster identity (alternative to the positional argument)")
	cmd.Flags().BoolVar(&walNoHeader, NoHeaderKey, false, "Do not print the header")
	cmd.Flags().BoolVar(&walShowExactTime, "exact-time", false, "Show exact time instead of relative time")
	cmd.Flags().StringSliceVar(&walKinds, "kinds", []string{}, "Filter by WAL object kinds (wal, backup_history, timeline_history)")
	_ = cmd.RegisterFlagCompletionFunc("backup-store", completion.CompleteDbBackupStoreID)
	_ = cmd.RegisterFlagCompletionFunc("cluster", completion.CompleteDbClusterID)
}

func init() {
	RecoveryCmd.AddCommand(walCmd)
	walCmd.AddCommand(walHierarchyCmd)
	walCmd.AddCommand(walListCmd)
	walCmd.AddCommand(walSearchCmd)

	registerWalFlags(walHierarchyCmd)
	walHierarchyCmd.Flags().StringVar(&walTimeline, "timeline", "", "Filter to a single timeline")

	registerWalFlags(walListCmd)
	walListCmd.Flags().BoolVar(&walFlat, "flat", false, "List WAL objects without requiring timeline and log")
	walListCmd.Flags().StringVar(&walTimeline, "timeline", "", "Timeline hex (required unless --flat)")
	walListCmd.Flags().StringVar(&walLog, "log", "", "Log hex (required unless --flat)")
	walListCmd.Flags().IntVar(&walLimit, "limit", 0, "Maximum number of objects to return")
	walListCmd.Flags().StringVar(&walContinuation, "continuation-token", "", "Continuation token from a previous listing")

	registerWalFlags(walSearchCmd)
	walSearchCmd.Flags().StringVar(&walQuery, "query", "", "WAL name or hex substring to search for (required)")
	walSearchCmd.Flags().IntVar(&walLimit, "limit", 0, "Maximum number of objects to return")
	_ = walSearchCmd.MarkFlagRequired("query")
}
