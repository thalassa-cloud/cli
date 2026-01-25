package nodepools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var (
	updateNodePoolCluster        *string
	updateNodePoolName           *string
	updateNodePoolMachineType    *string
	updateNodePoolReplicas       *int
	updateNodePoolEnableAS       *bool
	updateNodePoolMinNodes       *int
	updateNodePoolMaxNodes       *int
	updateNodePoolEnableAH       *bool
	updateNodePoolUpgradeStrat   *string
	updateNodePoolTaints         []string
	updateNodePoolSecurityGroups []string
	updateNodePoolWait           bool
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "edit", "modify"},
	Short:   "Update a Kubernetes node pool",
	Long: `Update an existing node pool in a Kubernetes cluster.

Only fields specified with flags will be updated. All other fields will remain unchanged.

Examples:
  # Update node pool replicas
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --num-nodes 5

  # Enable autoscaling
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --enable-autoscaling --min-nodes 2 --max-nodes 10

  # Update machine type
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --machine-type pgp-large

  # Update multiple fields
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --num-nodes 5 --enable-autohealing

Note: Use 'tcloud kubernetes nodepools label' and 'tcloud kubernetes nodepools annotate' commands
to manage labels and annotations separately.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if updateNodePoolCluster == nil || *updateNodePoolCluster == "" {
			return fmt.Errorf("--cluster is required")
		}
		if updateNodePoolName == nil || *updateNodePoolName == "" {
			return fmt.Errorf("--name is required")
		}

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
			if strings.EqualFold(c.Identity, *updateNodePoolCluster) ||
				strings.EqualFold(c.Name, *updateNodePoolCluster) ||
				strings.EqualFold(c.Slug, *updateNodePoolCluster) {
				cluster = &c
				break
			}
		}

		if cluster == nil {
			return fmt.Errorf("cluster not found: %s", *updateNodePoolCluster)
		}

		// Get node pools for the cluster
		nodePools, err := client.Kubernetes().ListKubernetesNodePools(ctx, cluster.Identity, &kubernetes.ListKubernetesNodePoolsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list node pools: %w", err)
		}

		// Find the node pool
		var nodePool *kubernetes.KubernetesNodePool
		for _, np := range nodePools {
			if strings.EqualFold(np.Identity, *updateNodePoolName) ||
				strings.EqualFold(np.Name, *updateNodePoolName) ||
				strings.EqualFold(np.Slug, *updateNodePoolName) {
				nodePool = &np
				break
			}
		}

		if nodePool == nil {
			return fmt.Errorf("node pool not found: %s", *updateNodePoolName)
		}

		// Get fresh node pool data
		nodePool, err = client.Kubernetes().GetKubernetesNodePool(ctx, cluster.Identity, nodePool.Identity)
		if err != nil {
			return fmt.Errorf("failed to get node pool: %w", err)
		}

		// Build update request, starting with existing values
		updateReq := kubernetes.UpdateKubernetesNodePool{}

		// Update machine type if specified
		if updateNodePoolMachineType != nil && *updateNodePoolMachineType != "" {
			resolvedMachineType, err := ResolveMachineType(ctx, client, *updateNodePoolMachineType)
			if err != nil {
				return err
			}
			updateReq.MachineType = resolvedMachineType
		}

		// Update replicas/autoscaling
		if updateNodePoolEnableAS != nil {
			updateReq.EnableAutoscaling = updateNodePoolEnableAS

			// Validate autoscaling configuration if enabling
			if *updateNodePoolEnableAS {
				minNodes := nodePool.MinReplicas
				maxNodes := nodePool.MaxReplicas
				if updateNodePoolMinNodes != nil {
					minNodes = *updateNodePoolMinNodes
				}
				if updateNodePoolMaxNodes != nil {
					maxNodes = *updateNodePoolMaxNodes
				}

				if err := ValidateAutoscalingConfig(true, minNodes, maxNodes, 0); err != nil {
					return err
				}

				if updateNodePoolMinNodes != nil {
					updateReq.MinReplicas = updateNodePoolMinNodes
				}
				if updateNodePoolMaxNodes != nil {
					updateReq.MaxReplicas = updateNodePoolMaxNodes
				}
			} else {
				// Disabling autoscaling, need replicas
				replicas := nodePool.Replicas
				if updateNodePoolReplicas != nil {
					replicas = *updateNodePoolReplicas
				}
				if replicas < 1 {
					return fmt.Errorf("num-nodes must be at least 1 when autoscaling is disabled")
				}
				updateReq.Replicas = &replicas
			}
		} else {
			// Autoscaling not being changed, but replicas might be
			if updateNodePoolReplicas != nil {
				if nodePool.EnableAutoscaling {
					return fmt.Errorf("cannot set --num-nodes when autoscaling is enabled. Use --min-nodes and --max-nodes instead, or disable autoscaling first")
				}
				if *updateNodePoolReplicas < 1 {
					return fmt.Errorf("num-nodes must be at least 1")
				}
				updateReq.Replicas = updateNodePoolReplicas
			}
			if updateNodePoolMinNodes != nil {
				if !nodePool.EnableAutoscaling {
					return fmt.Errorf("cannot set --min-nodes when autoscaling is disabled. Enable autoscaling first")
				}
				updateReq.MinReplicas = updateNodePoolMinNodes
			}
			if updateNodePoolMaxNodes != nil {
				if !nodePool.EnableAutoscaling {
					return fmt.Errorf("cannot set --max-nodes when autoscaling is disabled. Enable autoscaling first")
				}
				updateReq.MaxReplicas = updateNodePoolMaxNodes
			}
		}

		// Update autohealing
		if updateNodePoolEnableAH != nil {
			updateReq.EnableAutoHealing = updateNodePoolEnableAH
		}

		// Update upgrade strategy
		if updateNodePoolUpgradeStrat != nil && *updateNodePoolUpgradeStrat != "" {
			strategy := ParseUpgradeStrategy(*updateNodePoolUpgradeStrat)
			updateReq.UpgradeStrategy = &strategy
		}

		// Update taints
		if len(updateNodePoolTaints) > 0 {
			taints, err := ParseTaints(updateNodePoolTaints)
			if err != nil {
				return err
			}
			nodeSettings := kubernetes.KubernetesNodeSettings{}

			// Preserve existing labels and annotations
			if nodePool.NodeSettings.Labels != nil {
				nodeSettings.Labels = nodePool.NodeSettings.Labels
			}
			if nodePool.NodeSettings.Annotations != nil {
				nodeSettings.Annotations = nodePool.NodeSettings.Annotations
			}

			nodeSettings.Taints = taints
			updateReq.NodeSettings = &nodeSettings
		}

		// Update security groups
		if len(updateNodePoolSecurityGroups) > 0 {
			updateReq.SecurityGroupAttachments = updateNodePoolSecurityGroups
		}

		// Check if there are any updates
		if updateReq.MachineType == "" &&
			updateReq.EnableAutoscaling == nil &&
			updateReq.Replicas == nil &&
			updateReq.MinReplicas == nil &&
			updateReq.MaxReplicas == nil &&
			updateReq.EnableAutoHealing == nil &&
			updateReq.UpgradeStrategy == nil &&
			updateReq.NodeSettings == nil &&
			len(updateReq.SecurityGroupAttachments) == 0 {
			return fmt.Errorf("no updates specified")
		}

		updatedNodePool, err := client.Kubernetes().UpdateKubernetesNodePool(ctx, cluster.Identity, nodePool.Identity, updateReq)
		if err != nil {
			return fmt.Errorf("failed to update node pool: %w", err)
		}

		if updateNodePoolWait {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
			defer cancel()

			_, err = client.Kubernetes().WaitUntilKubernetesNodePoolReady(ctxWithTimeout, cluster.Identity, updatedNodePool.Identity)
			if err != nil {
				return fmt.Errorf("failed to wait for node pool to be ready: %w", err)
			}
			// Get fresh node pool data
			updatedNodePool, err = client.Kubernetes().GetKubernetesNodePool(ctx, cluster.Identity, updatedNodePool.Identity)
			if err != nil {
				return fmt.Errorf("failed to get node pool: %w", err)
			}
		}

		// Output in table format like IaaS commands
		replicas := formatReplicasForUpdate(updatedNodePool)
		body := [][]string{
			{
				updatedNodePool.Identity,
				updatedNodePool.Name,
				cluster.Name,
				replicas,
				updatedNodePool.MachineType.Name,
				string(updatedNodePool.Status),
				formattime.FormatTime(updatedNodePool.CreatedAt.Local(), false),
			},
		}
		table.Print([]string{"ID", "Name", "Cluster", "Replicas", "Type", "Status", "Age"}, body)

		return nil
	},
}

func init() {
	// Command is registered in nodepools.go

	updateNodePoolCluster = updateCmd.Flags().String("cluster", "", "Cluster identity, name, or slug (required)")
	updateNodePoolName = updateCmd.Flags().String("name", "", "Node pool name, identity, or slug (required)")
	updateNodePoolMachineType = updateCmd.Flags().String("machine-type", "", "Machine type for the node pool")
	updateNodePoolReplicas = updateCmd.Flags().Int("num-nodes", 0, "Number of nodes in the node pool (only when autoscaling is disabled)")
	updateNodePoolEnableAS = updateCmd.Flags().Bool("enable-autoscaling", false, "Enable autoscaling for the node pool")
	updateNodePoolMinNodes = updateCmd.Flags().Int("min-nodes", 0, "Minimum number of nodes (when autoscaling is enabled)")
	updateNodePoolMaxNodes = updateCmd.Flags().Int("max-nodes", 0, "Maximum number of nodes (when autoscaling is enabled)")
	updateNodePoolEnableAH = updateCmd.Flags().Bool("enable-autohealing", false, "Enable autohealing for the node pool")
	updateNodePoolUpgradeStrat = updateCmd.Flags().String("upgrade-strategy", "", "Upgrade strategy: manual, auto, always, on-delete, inplace, or never")
	updateCmd.Flags().StringSliceVar(&updateNodePoolTaints, "node-taints", []string{}, "Node taints in key=value:effect or key:effect format (e.g., 'dedicated=gpu:NoSchedule'). Replaces existing taints.")
	updateCmd.Flags().StringSliceVar(&updateNodePoolSecurityGroups, "security-groups", []string{}, "Security group identities to attach to node pool machines")
	updateCmd.Flags().BoolVar(&updateNodePoolWait, "wait", false, "Wait for the node pool update to complete")

	updateCmd.RegisterFlagCompletionFunc(ClusterFlag, completion.CompleteKubernetesCluster)
	updateCmd.RegisterFlagCompletionFunc("machine-type", completion.CompleteMachineType)
	updateCmd.RegisterFlagCompletionFunc("upgrade-strategy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"manual", "auto", "always", "on-delete", "inplace", "never"}, cobra.ShellCompDirectiveNoFileComp
	})
}

// formatReplicasForUpdate formats the replicas string for update output
func formatReplicasForUpdate(np *kubernetes.KubernetesNodePool) string {
	if np.EnableAutoscaling {
		return fmt.Sprintf("%d-%d", np.MinReplicas, np.MaxReplicas)
	}
	return fmt.Sprintf("%d", np.Replicas)
}
