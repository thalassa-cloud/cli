package dbaas

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	createClusterName                                 string
	createClusterDescription                          string
	createClusterEngine                               string
	createClusterEngineVersion                        string
	createClusterInstanceType                         string
	createClusterVpc                                  string
	createClusterSubnet                               string
	createClusterStorage                              int
	createclusterVolumeType                           string
	createClusterReplicas                             int
	createClusterLabels                               []string
	createClusterAnnotations                          []string
	createClusterDeleteProtection                     bool
	createClusterWait                                 bool
	createClusterWaitTimeout                          time.Duration
	createClusterProvisionDbBackupObjectStorageBucket bool
	createClusterDbBackupObjectStorageId              string
	createClusterRestoreFromBackup                    string
	createClusterRestoreTargetTime                    string
	createClusterRestoreTargetLSN                     string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"create-cluster", "cluster-create"},
	Short:   "Create a database cluster",
	Long:    "Create a new database cluster in the Thalassa Cloud Platform. Use --restore-from-backup to create a cluster from an existing backup.",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if createClusterName == "" {
			return fmt.Errorf("name is required")
		}
		if createClusterEngine == "" {
			return fmt.Errorf("engine is required")
		}
		if createClusterEngineVersion == "" {
			return fmt.Errorf("engine-version is required")
		}
		if createClusterInstanceType == "" {
			return fmt.Errorf("instance-type is required")
		}
		if createClusterStorage <= 0 {
			return fmt.Errorf("storage must be greater than 0")
		}
		if createClusterReplicas < 0 {
			return fmt.Errorf("replicas must be 0 or greater")
		}
		if createClusterRestoreTargetTime != "" && createClusterRestoreTargetLSN != "" {
			return fmt.Errorf("--restore-target-time and --restore-target-lsn are mutually exclusive")
		}
		if (createClusterRestoreTargetTime != "" || createClusterRestoreTargetLSN != "") && createClusterRestoreFromBackup == "" {
			return fmt.Errorf("--restore-from-backup is required when using a restore recovery target")
		}

		// Resolve volume type
		volumeTypes, err := client.IaaS().ListVolumeTypes(cmd.Context(), &iaas.ListVolumeTypesRequest{})
		if err != nil {
			return fmt.Errorf("failed to get volume type: %w", err)
		}
		volumeType, err := iaasutil.FindVolumeTypeByIdentitySlugOrNameWithError(volumeTypes, createclusterVolumeType)
		if err != nil {
			return err
		}

		createclusterVolumeType = volumeType.Identity

		// Resolve subnet (required)
		if createClusterSubnet == "" {
			return fmt.Errorf("subnet is required")
		}
		subnet, err := iaasutil.GetSubnetByIdentitySlugOrName(cmd.Context(), client.IaaS(), createClusterSubnet)
		if err != nil {
			return fmt.Errorf("failed to get subnet: %w", err)
		}
		subnetIdentity := subnet.Identity

		// Resolve and validate VPC if provided
		if createClusterVpc != "" {
			vpc, err := iaasutil.GetVPCByIdentitySlugOrName(cmd.Context(), client.IaaS(), createClusterVpc)
			if err != nil {
				return fmt.Errorf("failed to get vpc: %w", err)
			}
			// Validate that the subnet belongs to the specified VPC
			if subnet.Vpc != nil && subnet.Vpc.Identity != vpc.Identity {
				return fmt.Errorf("subnet %s does not belong to VPC %s", subnetIdentity, vpc.Identity)
			}
		}

		// Parse labels from key=value format
		labels := make(map[string]string)
		for _, label := range createClusterLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			} else {
				labels[label] = ""
			}
		}

		// Parse annotations from key=value format
		annotations := make(map[string]string)
		for _, annotation := range createClusterAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			} else {
				annotations[annotation] = ""
			}
		}

		req := dbaas.CreateDbClusterRequest{
			Name:                         createClusterName,
			Description:                  createClusterDescription,
			Engine:                       dbaas.DbClusterDatabaseEngine(createClusterEngine),
			EngineVersion:                createClusterEngineVersion,
			DatabaseInstanceTypeIdentity: createClusterInstanceType,
			VolumeTypeClassIdentity:      createclusterVolumeType,
			SubnetIdentity:               subnetIdentity,
			AllocatedStorage:             uint64(createClusterStorage),
			Labels:                       labels,
			Annotations:                  annotations,
			DeleteProtection:             createClusterDeleteProtection,
			ProvisionDbObjectStore:       createClusterProvisionDbBackupObjectStorageBucket,
		}
		if createClusterDbBackupObjectStorageId != "" { // todo; ensure this exists
			req.DbObjectStoreIdentity = &createClusterDbBackupObjectStorageId
		}

		if createClusterReplicas > 0 {
			req.Replicas = createClusterReplicas
		}

		if createClusterRestoreFromBackup != "" {
			_, err := client.DBaaS().GetDbBackup(cmd.Context(), createClusterRestoreFromBackup)
			if err != nil {
				if tcclient.IsNotFound(err) {
					return fmt.Errorf("backup not found: %s", createClusterRestoreFromBackup)
				}
				return fmt.Errorf("failed to get backup: %w", err)
			}
			req.RestoreFromBackupIdentity = &createClusterRestoreFromBackup

			if createClusterRestoreTargetTime != "" || createClusterRestoreTargetLSN != "" {
				recoveryTarget := &dbaas.RestoreRecoveryTarget{}
				if createClusterRestoreTargetTime != "" {
					targetTime, err := parseBarmanRestoreTargetTime(createClusterRestoreTargetTime)
					if err != nil {
						return fmt.Errorf("invalid restore target time: %w", err)
					}
					recoveryTarget.TargetTime = &targetTime
				}
				if createClusterRestoreTargetLSN != "" {
					lsn := strings.TrimSpace(createClusterRestoreTargetLSN)
					recoveryTarget.TargetLSN = &lsn
				}
				req.RestoreRecoveryTarget = recoveryTarget
			}
		}

		cluster, err := client.DBaaS().CreateDbCluster(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create database cluster: %w", err)
		}

		if createClusterWait {
			waitCtx, cancel, err := shared.WaitContext(cmd.Context(), createClusterWaitTimeout)
			if err != nil {
				return err
			}
			defer cancel()
			for {
				cluster, err = client.DBaaS().GetDbCluster(waitCtx, cluster.Identity)
				if err != nil {
					return fmt.Errorf("failed to get database cluster: %w", err)
				}
				if cluster.Status == dbaas.DbClusterStatusReady {
					break
				}
				if cluster.Status == dbaas.DbClusterStatusFailed {
					return fmt.Errorf("database cluster creation failed")
				}
				select {
				case <-waitCtx.Done():
					return waitCtx.Err()
				case <-time.After(5 * time.Second):
				}
			}
		}

		// Output in table format
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

		body := [][]string{
			{
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
			},
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "VPC", "Engine", "Version", "Instance Type", "Replicas", "Storage", "Status", "Age"}, body)
		}

		return nil
	},
}

func init() {
	DbaasCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createClusterName, "name", "", "Name of the database cluster (required)")
	createCmd.Flags().StringVar(&createClusterDescription, "description", "", "Description of the database cluster")
	createCmd.Flags().StringVar(&createClusterEngine, "engine", "", "Database engine (e.g., postgres) (required)")
	createCmd.Flags().StringVar(&createClusterEngineVersion, "engine-version", "", "Engine version (required)")
	createCmd.Flags().StringVar(&createClusterInstanceType, "instance-type", "", "Instance type (required)")

	createCmd.Flags().StringVar(&createClusterVpc, "vpc", "", "VPC identity, slug, or name")
	createCmd.Flags().StringVar(&createClusterSubnet, "subnet", "", "Subnet identity, slug, or name (required)")
	createCmd.Flags().StringVar(&createclusterVolumeType, "volume-type", "block", "Volume type")
	createCmd.Flags().IntVar(&createClusterStorage, "storage", 0, "Storage size in GB (required)")
	createCmd.Flags().IntVar(&createClusterReplicas, "replicas", 0, "Number of replicas (default: 0)")
	createCmd.Flags().StringSliceVar(&createClusterLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	createCmd.Flags().StringSliceVar(&createClusterAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	createCmd.Flags().BoolVar(&createClusterDeleteProtection, "delete-protection", false, "Enable delete protection")
	createCmd.Flags().BoolVar(&createClusterWait, "wait", false, "Wait for the database cluster to be available before returning")
	createCmd.Flags().DurationVar(&createClusterWaitTimeout, "wait-timeout", 20*time.Minute, "Maximum time to wait for the database cluster to be available")
	createCmd.Flags().BoolVar(&createClusterProvisionDbBackupObjectStorageBucket, "with-backup-bucket", false, "Provision a backup object storage bucket for the database cluster")
	createCmd.Flags().StringVar(&createClusterDbBackupObjectStorageId, "backup-object-storage-id", "", "Backup object storage ID (enables backup storage, requires --with-backup-bucket=false)")
	createCmd.Flags().StringVar(&createClusterRestoreFromBackup, "restore-from-backup", "", "Backup identity to restore the cluster from")
	createCmd.Flags().StringVar(&createClusterRestoreTargetTime, "restore-target-time", "", barmanTargetTimeDescription+" (requires --restore-from-backup)")
	createCmd.Flags().StringVar(&createClusterRestoreTargetLSN, "restore-target-lsn", "", "Point-in-time recovery target LSN (requires --restore-from-backup)")

	// Register completions
	_ = createCmd.RegisterFlagCompletionFunc("vpc", completion.CompleteVPCID)
	_ = createCmd.RegisterFlagCompletionFunc("subnet", completion.CompleteSubnetEnhanced)
	_ = createCmd.RegisterFlagCompletionFunc("engine-version", completion.CompleteDbEngineVersion)
	_ = createCmd.RegisterFlagCompletionFunc("restore-from-backup", completion.CompleteDbBackupID)

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("engine")
	_ = createCmd.MarkFlagRequired("engine-version")
	_ = createCmd.MarkFlagRequired("instance-type")
	_ = createCmd.MarkFlagRequired("subnet")
	_ = createCmd.MarkFlagRequired("volume-type")
	_ = createCmd.MarkFlagRequired("storage")
}
