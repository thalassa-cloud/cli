package dbaas

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
)

var (
	createClusterName             string
	createClusterDescription      string
	createClusterEngine           string
	createClusterEngineVersion    string
	createClusterInstanceType     string
	createClusterSubnet           string
	createClusterStorage          int
	createClusterReplicas         int
	createClusterLabels           []string
	createClusterAnnotations      []string
	createClusterDeleteProtection bool
	createClusterWait             bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a database cluster",
	Long:  "Create a new database cluster in the Thalassa Cloud Platform.",
	Args:  cobra.NoArgs,
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

		// Resolve subnet if provided
		var subnetIdentity string
		if createClusterSubnet != "" {
			subnet, err := client.IaaS().GetSubnet(cmd.Context(), createClusterSubnet)
			if err != nil {
				return fmt.Errorf("failed to get subnet: %w", err)
			}
			subnetIdentity = subnet.Identity
		}

		// Parse labels from key=value format
		labels := make(map[string]string)
		for _, label := range createClusterLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Parse annotations from key=value format
		annotations := make(map[string]string)
		for _, annotation := range createClusterAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		req := dbaas.CreateDbClusterRequest{
			Name:                         createClusterName,
			Description:                  createClusterDescription,
			Engine:                       dbaas.DbClusterDatabaseEngine(createClusterEngine),
			EngineVersion:                createClusterEngineVersion,
			DatabaseInstanceTypeIdentity: createClusterInstanceType,
			SubnetIdentity:               subnetIdentity,
			AllocatedStorage:             uint64(createClusterStorage),
			Labels:                       labels,
			Annotations:                  annotations,
			DeleteProtection:             createClusterDeleteProtection,
		}

		if createClusterReplicas > 0 {
			req.Replicas = createClusterReplicas
		}

		cluster, err := client.DBaaS().CreateDbCluster(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create database cluster: %w", err)
		}

		if createClusterWait {
			// Poll until cluster is available
			// Note: WaitUntilDbClusterIsAvailable may not be available, so we poll manually
			for {
				cluster, err = client.DBaaS().GetDbCluster(cmd.Context(), cluster.Identity)
				if err != nil {
					return fmt.Errorf("failed to get database cluster: %w", err)
				}
				if cluster.Status == dbaas.DbClusterStatusReady {
					break
				}
				if cluster.Status == dbaas.DbClusterStatusFailed {
					return fmt.Errorf("database cluster creation failed")
				}
				// Simple polling with sleep
				select {
				case <-cmd.Context().Done():
					return cmd.Context().Err()
				case <-time.After(5 * time.Second):
					// Continue polling
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
	createCmd.Flags().StringVar(&createClusterSubnet, "subnet", "", "Subnet identity")
	createCmd.Flags().IntVar(&createClusterStorage, "storage", 0, "Storage size in GB (required)")
	createCmd.Flags().IntVar(&createClusterReplicas, "replicas", 0, "Number of replicas (default: 0)")
	createCmd.Flags().StringSliceVar(&createClusterLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	createCmd.Flags().StringSliceVar(&createClusterAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	createCmd.Flags().BoolVar(&createClusterDeleteProtection, "delete-protection", false, "Enable delete protection")
	createCmd.Flags().BoolVar(&createClusterWait, "wait", false, "Wait for the database cluster to be available before returning")

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("engine")
	_ = createCmd.MarkFlagRequired("engine-version")
	_ = createCmd.MarkFlagRequired("instance-type")
	_ = createCmd.MarkFlagRequired("vpc")
	_ = createCmd.MarkFlagRequired("storage")
}
