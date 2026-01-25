package nodepools

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/fzf"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var (
	createNodePoolCluster        string
	createNodePoolName           string
	createNodePoolMachineType    string
	createNodePoolReplicas       int
	createNodePoolEnableAS       bool
	createNodePoolMinNodes       int
	createNodePoolMaxNodes       int
	createNodePoolSubnet         string
	createNodePoolAZs            []string
	createNodePoolEnableAH       bool
	createNodePoolUpgradeStrat   string
	createNodePoolLabels         []string
	createNodePoolAnnotations    []string
	createNodePoolTaints         []string
	createNodePoolSecurityGroups []string
	createNodePoolWait           bool
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c", "new"},
	Short:   "Create a Kubernetes node pool",
	Long: `Create a new node pool in a Kubernetes cluster.

Examples:
  # Create a node pool with minimal configuration
  tcloud kubernetes nodepools create --cluster my-cluster --name worker --machine-type pgp-medium --num-nodes 3 --availability-zone nl-01a

  # Create a node pool with autoscaling
  tcloud kubernetes nodepools create --cluster my-cluster --name worker --machine-type pgp-medium --enable-autoscaling --min-nodes 1 --max-nodes 5 --availability-zone nl-01a

  # Create node pools in multiple availability zones
  tcloud kubernetes nodepools create --cluster my-cluster --name worker --machine-type pgp-medium --num-nodes 3 --availability-zone nl-01a --availability-zone nl-01b`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if createNodePoolCluster == "" {
			return fmt.Errorf("--cluster is required")
		}
		clusterIdentifier := createNodePoolCluster

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Resolve cluster
		clusters, err := client.Kubernetes().ListKubernetesClusters(ctx, &kubernetes.ListKubernetesClustersRequest{})
		if err != nil {
			return fmt.Errorf("failed to list clusters: %w", err)
		}

		var cluster *kubernetes.KubernetesCluster
		for _, c := range clusters {
			if strings.EqualFold(c.Identity, clusterIdentifier) || strings.EqualFold(c.Name, clusterIdentifier) || strings.EqualFold(c.Slug, clusterIdentifier) {
				cluster = &c
				break
			}
		}

		if cluster == nil {
			return fmt.Errorf("cluster not found: %s", clusterIdentifier)
		}

		// Get fresh cluster data to ensure we have the latest version
		cluster, err = client.Kubernetes().GetKubernetesCluster(ctx, cluster.Identity)
		if err != nil {
			return fmt.Errorf("failed to get cluster: %w", err)
		}

		// Validate cluster is ready
		if !strings.EqualFold(cluster.Status, "ready") {
			return fmt.Errorf("cluster is not ready (status: %s). Please wait for the cluster to be ready before creating node pools", cluster.Status)
		}

		// Resolve machine type
		if createNodePoolMachineType == "" {
			if fzf.IsInteractiveMode(os.Stdout) {
				command := fmt.Sprintf("%s compute machine-types --no-header", os.Args[0])
				selected, err := fzf.InteractiveChoice(command)
				if err != nil {
					return fmt.Errorf("machine-type is required: %w", err)
				}
				createNodePoolMachineType = selected
			} else {
				return fmt.Errorf("machine-type is required")
			}
		}

		resolvedMachineType, err := ResolveMachineType(ctx, client, createNodePoolMachineType)
		if err != nil {
			return err
		}

		// Validate autoscaling configuration
		if err := ValidateAutoscalingConfig(createNodePoolEnableAS, createNodePoolMinNodes, createNodePoolMaxNodes, createNodePoolReplicas); err != nil {
			return err
		}

		// Determine node pool subnet
		nodePoolSubnetIdentity, err := ResolveNodePoolSubnet(ctx, client, createNodePoolSubnet, cluster)
		if err != nil {
			return err
		}

		// Determine availability zones
		availabilityZones, err := DetermineAvailabilityZones(ctx, client, createNodePoolAZs, cluster)
		if err != nil {
			return err
		}

		// Node pool name
		baseNodePoolName := createNodePoolName
		if baseNodePoolName == "" {
			baseNodePoolName = "worker"
		}

		// Build node pool config
		config := NodePoolConfig{
			Name:           baseNodePoolName,
			MachineType:    resolvedMachineType,
			Replicas:       createNodePoolReplicas,
			EnableAS:       createNodePoolEnableAS,
			MinNodes:       createNodePoolMinNodes,
			MaxNodes:       createNodePoolMaxNodes,
			EnableAH:       createNodePoolEnableAH,
			UpgradeStrat:   createNodePoolUpgradeStrat,
			Labels:         createNodePoolLabels,
			Annotations:    createNodePoolAnnotations,
			Taints:         createNodePoolTaints,
			SecurityGroups: createNodePoolSecurityGroups,
		}

		// Create node pool in each availability zone
		for _, az := range availabilityZones {
			// Determine node pool name
			nodePoolName := BuildNodePoolName(baseNodePoolName, az, len(availabilityZones) > 1)

			// Build create request
			createNodePoolReq, err := BuildNodePoolCreateRequest(config, cluster, az, nodePoolSubnetIdentity)
			if err != nil {
				return err
			}
			createNodePoolReq.Name = nodePoolName

			nodePool, err := client.Kubernetes().CreateKubernetesNodePool(ctx, cluster.Identity, createNodePoolReq)
			if err != nil {
				return fmt.Errorf("failed to create node pool in availability zone %s: %w", az, err)
			}

			if createNodePoolWait {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
				defer cancel()

				_, err = client.Kubernetes().WaitUntilKubernetesNodePoolReady(ctxWithTimeout, cluster.Identity, nodePool.Identity)
				if err != nil {
					return fmt.Errorf("failed to wait for node pool to be ready: %w", err)
				}
				// Get fresh node pool data
				nodePool, err = client.Kubernetes().GetKubernetesNodePool(ctx, cluster.Identity, nodePool.Identity)
				if err != nil {
					return fmt.Errorf("failed to get node pool: %w", err)
				}
			}

			// Output in table format
			replicas := formatReplicas(nodePool)
			body := [][]string{
				{
					nodePool.Identity,
					nodePool.Name,
					cluster.Name,
					replicas,
					resolvedMachineType,
					string(nodePool.Status),
					formattime.FormatTime(nodePool.CreatedAt.Local(), false),
				},
			}
			table.Print([]string{"ID", "Name", "Cluster", "Replicas", "Type", "Status", "Age"}, body)
		}

		return nil
	},
}

func init() {
	// Command is registered in kubernetesclusters.go

	createCmd.Flags().StringVar(&createNodePoolCluster, "cluster", "", "Cluster identity, name, or slug (required)")
	createCmd.Flags().StringVar(&createNodePoolName, "name", "worker", "Name of the node pool (default: worker)")
	createCmd.Flags().StringVar(&createNodePoolMachineType, "machine-type", "", "Machine type for the node pool (required)")
	createCmd.Flags().IntVar(&createNodePoolReplicas, "num-nodes", 1, "Number of nodes in the node pool (ignored if --enable-autoscaling is set)")
	createCmd.Flags().BoolVar(&createNodePoolEnableAS, "enable-autoscaling", false, "Enable autoscaling for the node pool")
	createCmd.Flags().IntVar(&createNodePoolMinNodes, "min-nodes", 0, "Minimum number of nodes (required when autoscaling is enabled)")
	createCmd.Flags().IntVar(&createNodePoolMaxNodes, "max-nodes", 3, "Maximum number of nodes (required when autoscaling is enabled)")
	createCmd.Flags().StringVar(&createNodePoolSubnet, "subnet", "", "Subnet for the node pool (defaults to cluster subnet)")
	createCmd.Flags().StringSliceVar(&createNodePoolAZs, "availability-zone", []string{}, "Availability zone for the node pool (can be specified multiple times). If not specified, a random AZ from the cluster's region will be selected.")
	createCmd.Flags().BoolVar(&createNodePoolEnableAH, "enable-autohealing", false, "Enable autohealing for the node pool")
	createCmd.Flags().StringVar(&createNodePoolUpgradeStrat, "upgrade-strategy", "auto", "Upgrade strategy: manual, auto, always, on-delete, inplace, or never")
	createCmd.Flags().StringSliceVar(&createNodePoolLabels, "node-labels", []string{}, "Node labels in key=value format (applied to Kubernetes nodes)")
	createCmd.Flags().StringSliceVar(&createNodePoolAnnotations, "node-annotations", []string{}, "Node annotations in key=value format (applied to Kubernetes nodes)")
	createCmd.Flags().StringSliceVar(&createNodePoolTaints, "node-taints", []string{}, "Node taints in key=value:effect or key:effect format (e.g., 'dedicated=gpu:NoSchedule')")
	createCmd.Flags().StringSliceVar(&createNodePoolSecurityGroups, "security-groups", []string{}, "Security group identities to attach to node pool machines")
	createCmd.Flags().BoolVar(&createNodePoolWait, "wait", false, "Wait for the node pool to be ready before returning")

	createCmd.RegisterFlagCompletionFunc("cluster", completion.CompleteKubernetesCluster)
	createCmd.RegisterFlagCompletionFunc("machine-type", completion.CompleteMachineType)
	createCmd.RegisterFlagCompletionFunc("subnet", completion.CompleteSubnetEnhanced)
	createCmd.RegisterFlagCompletionFunc("upgrade-strategy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"manual", "auto", "always", "on-delete", "inplace", "never"}, cobra.ShellCompDirectiveNoFileComp
	})
}
