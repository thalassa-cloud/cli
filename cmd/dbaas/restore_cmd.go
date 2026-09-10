package dbaas

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	dbaasutil "github.com/thalassa-cloud/cli/internal/dbaas"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/fzf"
	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

var (
	restoreName             string
	restoreDescription      string
	restoreCluster          string
	restoreFromBackup       string
	restoreTargetTime       string
	restoreTargetLSN        string
	restoreInstanceType     string
	restoreSubnet           string
	restoreVpc              string
	restoreVolumeType       string
	restoreStorage          int
	restoreReplicas         int
	restoreLabels           []string
	restoreAnnotations      []string
	restoreDeleteProtection bool
	restoreWait             bool
	restoreWaitTimeout      time.Duration
	restoreProvisionStore   bool
	restoreBackupStore      string
	restoreForce            bool
	restoreNoHeader         bool
)

var restoreCmd = &cobra.Command{
	Use:     "restore",
	Aliases: []string{"pitr-restore"},
	Short:   "Restore a database cluster from a backup",
	Long: `Create a new database cluster from an existing backup.

Restore always provisions a new cluster. Optionally specify a point-in-time
recovery target (timestamp or LSN). In a terminal with fzf installed, missing
cluster, backup, and subnet values can be selected interactively.`,
	Example: "tcloud dbaas restore --from-backup bak-123 --name restored-db --subnet subnet-123\ntcloud dbaas restore --cluster dbc-123 --target-time 2023-12-25T10:00:00Z --name restored-db",
	Args:    cobra.NoArgs,
	RunE:    runRestore,
}

func runRestore(cmd *cobra.Command, args []string) error {
	if restoreTargetTime != "" && restoreTargetLSN != "" {
		return fmt.Errorf("--target-time and --target-lsn are mutually exclusive")
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	interactive := fzf.IsInteractiveMode(os.Stdout)
	sourceCluster, err := resolveRestoreSourceCluster(cmd, client, interactive)
	if err != nil {
		return err
	}

	overview, err := fetchRestoreOverview(cmd, client, sourceCluster)
	if err != nil {
		return err
	}
	if overview != nil {
		printRestoreWindow(overview)
	}

	backup, sourceCluster, err := resolveRestoreBackup(cmd, client, sourceCluster, overview, interactive)
	if err != nil {
		return err
	}

	coverage := findBackupCoverage(overview, backup.Identity)
	if coverage != nil && !restoreForce && !dbaasutil.BackupCoverageIsRestorable(*coverage) {
		return fmt.Errorf("backup %s is not restorable (status %s); use --force to override", backup.Identity, coverage.Status)
	}

	if err := promptRestoreMode(interactive); err != nil {
		return err
	}

	recoveryTarget, err := buildRestoreRecoveryTarget(overview, coverage)
	if err != nil {
		return err
	}
	if aborted, err := confirmIncompleteWAL(coverage, recoveryTarget); err != nil {
		return err
	} else if aborted {
		return nil
	}

	if err := fillRestoreClusterDefaults(sourceCluster, interactive); err != nil {
		return err
	}

	req, err := buildRestoreClusterRequest(cmd, client, backup, recoveryTarget)
	if err != nil {
		return err
	}

	if !confirmRestore(backup, sourceCluster, recoveryTarget) {
		return nil
	}

	cluster, err := client.DBaaS().CreateDbCluster(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("failed to restore database cluster: %w", err)
	}

	if restoreWait {
		cluster, err = waitForRestoredCluster(cmd, client, cluster.Identity)
		if err != nil {
			return err
		}
	}

	printRestoredCluster(cluster)
	return nil
}

func resolveRestoreSourceCluster(cmd *cobra.Command, client thalassa.Client, interactive bool) (*dbaas.DbCluster, error) {
	if restoreCluster == "" && restoreFromBackup == "" {
		if !interactive {
			return nil, fmt.Errorf("cluster or --from-backup is required")
		}
		selected, err := fzf.InteractiveChoice(fmt.Sprintf("%s dbaas list --no-header", os.Args[0]))
		if err != nil {
			return nil, fmt.Errorf("select cluster: %w", err)
		}
		restoreCluster = selected
	}
	if restoreCluster == "" {
		return nil, nil
	}

	sourceCluster, err := client.DBaaS().GetDbCluster(cmd.Context(), restoreCluster)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return nil, fmt.Errorf("database cluster not found: %s", restoreCluster)
		}
		return nil, fmt.Errorf("failed to get database cluster: %w", err)
	}
	return sourceCluster, nil
}

func resolveRestoreBackup(cmd *cobra.Command, client thalassa.Client, sourceCluster *dbaas.DbCluster, overview *dbaas.DbObjectStoreRecoveryOverview, interactive bool) (*dbaas.DbClusterBackup, *dbaas.DbCluster, error) {
	if restoreFromBackup == "" {
		if !interactive {
			return nil, sourceCluster, fmt.Errorf("--from-backup is required")
		}
		selected, err := pickRestoreBackup(overview, restoreForce)
		if err != nil {
			return nil, sourceCluster, err
		}
		restoreFromBackup = selected
	}

	backup, err := client.DBaaS().GetDbBackup(cmd.Context(), restoreFromBackup)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return nil, sourceCluster, fmt.Errorf("backup not found: %s", restoreFromBackup)
		}
		return nil, sourceCluster, fmt.Errorf("failed to get backup: %w", err)
	}

	if sourceCluster == nil && backup.DbCluster != nil {
		sourceCluster, err = client.DBaaS().GetDbCluster(cmd.Context(), backup.DbCluster.Identity)
		if err != nil && !tcclient.IsNotFound(err) {
			return nil, nil, fmt.Errorf("failed to get source database cluster: %w", err)
		}
	}

	return backup, sourceCluster, nil
}

func promptRestoreMode(interactive bool) error {
	if restoreTargetTime != "" || restoreTargetLSN != "" || !interactive {
		return nil
	}

	mode, err := shared.PromptString("Restore mode (end, time, lsn)", "end")
	if err != nil {
		return err
	}
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "", "end", "backup":
		return nil
	case "time", "pitr":
		restoreTargetTime, err = shared.PromptString("Target time", "")
		if err != nil {
			return err
		}
		if restoreTargetTime == "" {
			return fmt.Errorf("target time is required for PITR")
		}
		return nil
	case "lsn":
		restoreTargetLSN, err = shared.PromptString("Target LSN", "")
		if err != nil {
			return err
		}
		if restoreTargetLSN == "" {
			return fmt.Errorf("target LSN is required")
		}
		return nil
	default:
		return fmt.Errorf("unknown restore mode %q: use end, time, or lsn", mode)
	}
}

func buildRestoreRecoveryTarget(overview *dbaas.DbObjectStoreRecoveryOverview, coverage *dbaas.DbObjectStoreRecoveryBackupCoverage) (*dbaas.RestoreRecoveryTarget, error) {
	if restoreTargetTime == "" && restoreTargetLSN == "" {
		return nil, nil
	}

	recoveryTarget := &dbaas.RestoreRecoveryTarget{}
	if restoreTargetTime != "" {
		parsed, err := dbaasutil.ParseRestoreTargetTime(restoreTargetTime)
		if err != nil {
			return nil, fmt.Errorf("invalid restore target time: %w", err)
		}
		if err := validateRestoreTargetTime(parsed, overview, coverage); err != nil {
			return nil, err
		}
		formatted := dbaasutil.FormatRestoreTargetTimeRFC3339(parsed)
		recoveryTarget.TargetTime = &formatted
		restoreTargetTime = formatted
	}
	if restoreTargetLSN != "" {
		lsn := strings.TrimSpace(restoreTargetLSN)
		recoveryTarget.TargetLSN = &lsn
	}
	return recoveryTarget, nil
}

func confirmIncompleteWAL(coverage *dbaas.DbObjectStoreRecoveryBackupCoverage, recoveryTarget *dbaas.RestoreRecoveryTarget) (aborted bool, err error) {
	if coverage == nil || recoveryTarget == nil {
		return false, nil
	}
	warning := dbaasutil.ChainStatusWarning(coverage.Coverage.ChainStatus, coverage.Coverage.FirstMissingWal)
	if warning == "" {
		return false, nil
	}

	fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
	proceed, err := shared.PromptYesNoUnlessForce(restoreForce, "Continue with restore anyway? [y/N]: ")
	if err != nil {
		return false, err
	}
	if !proceed {
		fmt.Println("Aborted")
		return true, nil
	}
	return false, nil
}

func buildRestoreClusterRequest(cmd *cobra.Command, client thalassa.Client, backup *dbaas.DbClusterBackup, recoveryTarget *dbaas.RestoreRecoveryTarget) (dbaas.CreateDbClusterRequest, error) {
	volumeTypes, err := client.IaaS().ListVolumeTypes(cmd.Context(), &iaas.ListVolumeTypesRequest{})
	if err != nil {
		return dbaas.CreateDbClusterRequest{}, fmt.Errorf("failed to get volume type: %w", err)
	}
	volumeType, err := iaasutil.FindVolumeTypeByIdentitySlugOrNameWithError(volumeTypes, restoreVolumeType)
	if err != nil {
		return dbaas.CreateDbClusterRequest{}, err
	}

	subnet, err := iaasutil.GetSubnetByIdentitySlugOrName(cmd.Context(), client.IaaS(), restoreSubnet)
	if err != nil {
		return dbaas.CreateDbClusterRequest{}, fmt.Errorf("failed to get subnet: %w", err)
	}
	if restoreVpc != "" {
		vpc, err := iaasutil.GetVPCByIdentitySlugOrName(cmd.Context(), client.IaaS(), restoreVpc)
		if err != nil {
			return dbaas.CreateDbClusterRequest{}, fmt.Errorf("failed to get vpc: %w", err)
		}
		if subnet.Vpc != nil && subnet.Vpc.Identity != vpc.Identity {
			return dbaas.CreateDbClusterRequest{}, fmt.Errorf("subnet %s does not belong to VPC %s", subnet.Identity, vpc.Identity)
		}
	}

	req := dbaas.CreateDbClusterRequest{
		Name:                         restoreName,
		Description:                  restoreDescription,
		Engine:                       backup.EngineType,
		EngineVersion:                backup.EngineVersion,
		DatabaseInstanceTypeIdentity: restoreInstanceType,
		VolumeTypeClassIdentity:      volumeType.Identity,
		SubnetIdentity:               subnet.Identity,
		AllocatedStorage:             uint64(restoreStorage),
		Labels:                       shared.KeyValuePairsToMap(restoreLabels),
		Annotations:                  shared.KeyValuePairsToMap(restoreAnnotations),
		DeleteProtection:             restoreDeleteProtection,
		ProvisionDbObjectStore:       restoreProvisionStore,
		RestoreFromBackupIdentity:    &restoreFromBackup,
		RestoreRecoveryTarget:        recoveryTarget,
	}
	if restoreReplicas >= 0 {
		req.Replicas = restoreReplicas
	}
	if restoreBackupStore != "" {
		req.DbObjectStoreIdentity = &restoreBackupStore
	}
	return req, nil
}

func waitForRestoredCluster(cmd *cobra.Command, client thalassa.Client, identity string) (*dbaas.DbCluster, error) {
	waitCtx, cancel, err := shared.WaitContext(cmd.Context(), restoreWaitTimeout)
	if err != nil {
		return nil, err
	}
	defer cancel()

	for {
		cluster, err := client.DBaaS().GetDbCluster(waitCtx, identity)
		if err != nil {
			return nil, fmt.Errorf("failed to get database cluster: %w", err)
		}
		if cluster.Status == dbaas.DbClusterStatusReady {
			return cluster, nil
		}
		if cluster.Status == dbaas.DbClusterStatusFailed {
			return nil, fmt.Errorf("database cluster restore failed")
		}
		select {
		case <-waitCtx.Done():
			return nil, waitCtx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

func fetchRestoreOverview(cmd *cobra.Command, client thalassa.Client, sourceCluster *dbaas.DbCluster) (*dbaas.DbObjectStoreRecoveryOverview, error) {
	if sourceCluster == nil {
		return nil, nil
	}
	ref, err := dbaasutil.ResolveBackupStore(cmd.Context(), client.DBaaS(), sourceCluster.Identity, "")
	if err != nil {
		// No backup store means we can still restore from a snapshot-style backup without PITR insight.
		fmt.Fprintf(os.Stderr, "Note: %s\n", err.Error())
		return nil, nil
	}

	overview, err := client.DBaaS().GetDbObjectStoreRecoveryOverview(cmd.Context(), ref.StoreIdentity, &dbaas.GetDbObjectStoreRecoveryOverviewRequest{
		ClusterID: sourceCluster.Identity,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get recovery overview: %w", err)
	}
	return overview, nil
}

func printRestoreWindow(overview *dbaas.DbObjectStoreRecoveryOverview) {
	fmt.Println("PITR window")
	body := [][]string{
		{"Earliest PITR From", dbaasutil.FormatOptionalTime(overview.Summary.EarliestPitrFrom, true)},
		{"Latest PITR Through", dbaasutil.FormatOptionalTime(overview.Summary.LatestPitrThroughApprox, true)},
		{"Gap Count", fmt.Sprintf("%d", overview.Summary.GapCount)},
	}
	table.Print([]string{"Field", "Value"}, body)
	if !overview.Complete || overview.Truncation.Truncated {
		fmt.Fprintln(os.Stderr, "Warning: recovery analysis is incomplete. Do not treat the PITR window as complete.")
	}
}

func pickRestoreBackup(overview *dbaas.DbObjectStoreRecoveryOverview, force bool) (string, error) {
	if overview == nil {
		return fzf.InteractiveChoice(fmt.Sprintf("%s dbaas backup list --no-header", os.Args[0]))
	}

	lines := make([]string, 0)
	for _, cluster := range overview.Clusters {
		for _, backup := range cluster.Backups {
			if !force && !dbaasutil.BackupCoverageIsRestorable(backup) {
				continue
			}
			lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
				backup.BackupIdentity,
				backup.Status,
				dbaasutil.FormatBytes(backup.SizeBytes),
				backup.Coverage.ChainStatus,
				dbaasutil.FormatOptionalTime(backup.Coverage.PitrThroughApproxAt, true),
			))
		}
	}
	if len(lines) == 0 {
		return "", fmt.Errorf("no restorable backups found in the recovery overview")
	}
	return fzf.InteractiveChoiceFromLines(lines)
}

func findBackupCoverage(overview *dbaas.DbObjectStoreRecoveryOverview, backupIdentity string) *dbaas.DbObjectStoreRecoveryBackupCoverage {
	if overview == nil {
		return nil
	}
	for i := range overview.Clusters {
		for j := range overview.Clusters[i].Backups {
			if overview.Clusters[i].Backups[j].BackupIdentity == backupIdentity {
				return &overview.Clusters[i].Backups[j]
			}
		}
	}
	return nil
}

func validateRestoreTargetTime(target time.Time, overview *dbaas.DbObjectStoreRecoveryOverview, coverage *dbaas.DbObjectStoreRecoveryBackupCoverage) error {
	if coverage != nil {
		if err := dbaasutil.ValidateTargetTime(target, dbaasutil.WindowFromBackupCoverage(*coverage)); err != nil {
			return err
		}
	}
	if overview != nil {
		if err := dbaasutil.ValidateTargetTime(target, dbaasutil.WindowFromStoreSummary(overview.Summary)); err != nil {
			return err
		}
	}
	return nil
}

func fillRestoreClusterDefaults(sourceCluster *dbaas.DbCluster, interactive bool) error {
	if restoreName == "" {
		if !interactive {
			return fmt.Errorf("name is required")
		}
		name, err := shared.PromptString("New cluster name", "")
		if err != nil {
			return err
		}
		if name == "" {
			return fmt.Errorf("name is required")
		}
		restoreName = name
	}

	if restoreInstanceType == "" && sourceCluster != nil && sourceCluster.DatabaseInstanceType != nil {
		restoreInstanceType = sourceCluster.DatabaseInstanceType.Identity
	}
	if restoreInstanceType == "" {
		if !interactive {
			return fmt.Errorf("instance-type is required")
		}
		selected, err := fzf.InteractiveChoice(fmt.Sprintf("%s dbaas instance-types --no-header", os.Args[0]))
		if err != nil {
			return fmt.Errorf("select instance type: %w", err)
		}
		restoreInstanceType = selected
	}

	if restoreStorage <= 0 && sourceCluster != nil && sourceCluster.AllocatedStorage > 0 {
		restoreStorage = int(sourceCluster.AllocatedStorage)
	}
	if restoreStorage <= 0 {
		if !interactive {
			return fmt.Errorf("storage must be greater than 0")
		}
		defaultStorage := ""
		if sourceCluster != nil && sourceCluster.AllocatedStorage > 0 {
			defaultStorage = strconv.FormatUint(sourceCluster.AllocatedStorage, 10)
		}
		value, err := shared.PromptString("Storage size in GB", defaultStorage)
		if err != nil {
			return err
		}
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || parsed <= 0 {
			return fmt.Errorf("storage must be greater than 0")
		}
		restoreStorage = parsed
	}

	if restoreSubnet == "" && sourceCluster != nil && sourceCluster.Subnet != nil {
		restoreSubnet = sourceCluster.Subnet.Identity
	}
	if restoreSubnet == "" {
		if !interactive {
			return fmt.Errorf("subnet is required")
		}
		selected, err := fzf.InteractiveChoice(fmt.Sprintf("%s networking subnets list --no-header", os.Args[0]))
		if err != nil {
			return fmt.Errorf("select subnet: %w", err)
		}
		restoreSubnet = selected
	}

	if restoreReplicas < 0 && sourceCluster != nil {
		restoreReplicas = sourceCluster.Replicas
	}

	return nil
}

func confirmRestore(backup *dbaas.DbClusterBackup, sourceCluster *dbaas.DbCluster, recoveryTarget *dbaas.RestoreRecoveryTarget) bool {
	sourceName := "-"
	if sourceCluster != nil {
		sourceName = fmt.Sprintf("%s (%s)", sourceCluster.Name, sourceCluster.Identity)
	}

	target := "backup end"
	if recoveryTarget != nil {
		if recoveryTarget.TargetTime != nil {
			target = "time " + *recoveryTarget.TargetTime
		}
		if recoveryTarget.TargetLSN != nil {
			target = "LSN " + *recoveryTarget.TargetLSN
		}
	}

	var summary strings.Builder
	fmt.Fprintf(&summary, "Restore will create a new cluster:\n")
	fmt.Fprintf(&summary, "  Name: %s\n", restoreName)
	fmt.Fprintf(&summary, "  Source cluster: %s\n", sourceName)
	fmt.Fprintf(&summary, "  Backup: %s\n", backup.Identity)
	fmt.Fprintf(&summary, "  Engine: %s %s\n", backup.EngineType, backup.EngineVersion)
	fmt.Fprintf(&summary, "  Recovery target: %s\n", target)
	fmt.Fprintf(&summary, "  Subnet: %s\n", restoreSubnet)
	fmt.Fprintf(&summary, "  Instance type: %s\n", restoreInstanceType)
	fmt.Fprintf(&summary, "  Storage: %d GB\n", restoreStorage)
	fmt.Fprintf(&summary, "Continue? [y/N]: ")

	proceed, err := shared.PromptYesNoUnlessForce(restoreForce, summary.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read confirmation: %v\n", err)
		return false
	}
	if !proceed {
		fmt.Println("Aborted")
		return false
	}
	return true
}

func printRestoredCluster(cluster *dbaas.DbCluster) {
	vpcName := ""
	if cluster.Vpc != nil {
		vpcName = cluster.Vpc.Name
	}
	engineVersion := cluster.EngineVersion
	if cluster.DatabaseEngineVersion != nil {
		engineVersion = cluster.DatabaseEngineVersion.EngineVersion
	}
	instanceType := ""
	if cluster.DatabaseInstanceType != nil {
		instanceType = cluster.DatabaseInstanceType.Name
	}

	body := [][]string{{
		cluster.Identity,
		cluster.Name,
		vpcName,
		string(cluster.Engine),
		engineVersion,
		instanceType,
		fmt.Sprintf("%d", cluster.Replicas),
		fmt.Sprintf("%d GB", cluster.AllocatedStorage),
		string(cluster.Status),
		formattime.FormatTime(cluster.CreatedAt.Local(), false),
	}}
	if restoreNoHeader {
		table.Print(nil, body)
	} else {
		table.Print([]string{"ID", "Name", "VPC", "Engine", "Version", "Instance Type", "Replicas", "Storage", "Status", "Age"}, body)
	}
}

func init() {
	DbaasCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().BoolVar(&restoreNoHeader, NoHeaderKey, false, "Do not print the header")
	restoreCmd.Flags().StringVar(&restoreName, "name", "", "Name of the restored database cluster (required)")
	restoreCmd.Flags().StringVar(&restoreDescription, "description", "", "Description of the restored database cluster")
	restoreCmd.Flags().StringVar(&restoreCluster, "cluster", "", "Source database cluster identity")
	restoreCmd.Flags().StringVar(&restoreFromBackup, "from-backup", "", "Backup identity to restore from")
	restoreCmd.Flags().StringVar(&restoreTargetTime, "target-time", "", restoreTargetTimeDescription)
	restoreCmd.Flags().StringVar(&restoreTargetLSN, "target-lsn", "", "Point-in-time recovery target LSN")
	restoreCmd.Flags().StringVar(&restoreInstanceType, "instance-type", "", "Instance type (defaults to the source cluster)")
	restoreCmd.Flags().StringVar(&restoreSubnet, "subnet", "", "Subnet identity, slug, or name (defaults to the source cluster)")
	restoreCmd.Flags().StringVar(&restoreVpc, "vpc", "", "VPC identity, slug, or name")
	restoreCmd.Flags().StringVar(&restoreVolumeType, "volume-type", "block", "Volume type")
	restoreCmd.Flags().IntVar(&restoreStorage, "storage", 0, "Storage size in GB (defaults to the source cluster)")
	restoreCmd.Flags().IntVar(&restoreReplicas, "replicas", -1, "Number of replicas (defaults to the source cluster)")
	restoreCmd.Flags().StringSliceVar(&restoreLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	restoreCmd.Flags().StringSliceVar(&restoreAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	restoreCmd.Flags().BoolVar(&restoreDeleteProtection, "delete-protection", false, "Enable delete protection")
	restoreCmd.Flags().BoolVar(&restoreWait, "wait", false, "Wait for the restored cluster to be available before returning")
	restoreCmd.Flags().DurationVar(&restoreWaitTimeout, "wait-timeout", 20*time.Minute, "Maximum time to wait for the restored cluster to be available")
	restoreCmd.Flags().BoolVar(&restoreProvisionStore, "with-backup-bucket", false, "Provision a backup store for the restored cluster")
	restoreCmd.Flags().StringVar(&restoreBackupStore, "backup-store", "", "Backup store identity to attach to the restored cluster")
	restoreCmd.Flags().BoolVar(&restoreForce, "force", false, "Skip confirmation and allow restore from backups with incomplete WAL coverage")

	_ = restoreCmd.RegisterFlagCompletionFunc("cluster", completion.CompleteDbClusterID)
	_ = restoreCmd.RegisterFlagCompletionFunc("from-backup", completion.CompleteDbBackupID)
	_ = restoreCmd.RegisterFlagCompletionFunc("vpc", completion.CompleteVPCID)
	_ = restoreCmd.RegisterFlagCompletionFunc("subnet", completion.CompleteSubnetEnhanced)
	_ = restoreCmd.RegisterFlagCompletionFunc("instance-type", completion.CompleteDbInstanceType)
	_ = restoreCmd.RegisterFlagCompletionFunc("backup-store", completion.CompleteDbBackupStoreID)
}
