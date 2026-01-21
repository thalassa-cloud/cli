package dbaas

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateClusterName             string
	updateClusterDescription      string
	updateClusterInstanceType     string
	updateClusterStorage          int
	updateClusterReplicas         int
	updateClusterLabels           []string
	updateClusterAnnotations      []string
	updateClusterDeleteProtection bool
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:               "update",
	Short:             "Update a database cluster",
	Long:              "Update properties of an existing database cluster.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDbClusterID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusterIdentity := args[0]

		// Get current cluster
		current, err := client.DBaaS().GetDbCluster(cmd.Context(), clusterIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("database cluster not found: %s", clusterIdentity)
			}
			return fmt.Errorf("failed to get cluster: %w", err)
		}

		req := dbaas.UpdateDbClusterRequest{
			Name:             current.Name,
			Description:      current.Description,
			Labels:           current.Labels,
			Annotations:      current.Annotations,
			DeleteProtection: current.DeleteProtection,
		}

		if current.DatabaseInstanceType != nil {
			req.DatabaseInstanceTypeIdentity = &current.DatabaseInstanceType.Identity
		}
		req.AllocatedStorage = current.AllocatedStorage
		if current.Replicas > 0 {
			req.Replicas = current.Replicas
		}

		// Update name if provided
		if cmd.Flags().Changed("name") {
			req.Name = updateClusterName
		}

		// Update description if provided
		if cmd.Flags().Changed("description") {
			req.Description = updateClusterDescription
		}

		// Update instance type if provided
		if cmd.Flags().Changed("instance-type") {
			req.DatabaseInstanceTypeIdentity = &updateClusterInstanceType
		}

		// Update storage if provided
		if cmd.Flags().Changed("storage") {
			req.AllocatedStorage = uint64(updateClusterStorage)
		}

		// Update replicas if provided
		if cmd.Flags().Changed("replicas") {
			req.Replicas = updateClusterReplicas
		}

		// Parse labels from key=value format
		if cmd.Flags().Changed("labels") {
			labels := make(map[string]string)
			for _, label := range updateClusterLabels {
				parts := strings.SplitN(label, "=", 2)
				if len(parts) == 2 {
					labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Labels = labels
		}

		// Parse annotations from key=value format
		if cmd.Flags().Changed("annotations") {
			annotations := make(map[string]string)
			for _, annotation := range updateClusterAnnotations {
				parts := strings.SplitN(annotation, "=", 2)
				if len(parts) == 2 {
					annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Annotations = annotations
		}

		// Update delete protection if provided
		if cmd.Flags().Changed("delete-protection") {
			req.DeleteProtection = updateClusterDeleteProtection
		}

		cluster, err := client.DBaaS().UpdateDbCluster(cmd.Context(), clusterIdentity, req)
		if err != nil {
			return fmt.Errorf("failed to update database cluster: %w", err)
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
			},
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "VPC", "Engine", "Version", "Instance Type", "Replicas", "Storage", "Status"}, body)
		}

		return nil
	},
}

func init() {
	DbaasCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	updateCmd.Flags().StringVar(&updateClusterName, "name", "", "Name of the database cluster")
	updateCmd.Flags().StringVar(&updateClusterDescription, "description", "", "Description of the database cluster")
	updateCmd.Flags().StringVar(&updateClusterInstanceType, "instance-type", "", "Instance type")
	updateCmd.Flags().IntVar(&updateClusterStorage, "storage", 0, "Storage size in GB")
	updateCmd.Flags().IntVar(&updateClusterReplicas, "replicas", -1, "Number of replicas")
	updateCmd.Flags().StringSliceVar(&updateClusterLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	updateCmd.Flags().StringSliceVar(&updateClusterAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	updateCmd.Flags().BoolVar(&updateClusterDeleteProtection, "delete-protection", false, "Enable or disable delete protection")
}
