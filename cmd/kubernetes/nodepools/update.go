package nodepools

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/thalassa"
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
	updateNodePoolVersion        *string
	updateNodePoolWait           bool
	updateNodePoolWaitTimeout    time.Duration
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "edit", "modify"},
	Short:   "Update a Kubernetes node pool",
	Long: `Update an existing node pool in a Kubernetes cluster.

Only fields specified with flags will be updated. All other fields will remain unchanged.
The update API replaces the node pool, so omitted settings are filled from the current node pool.

Examples:
  # Update node pool replicas
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --num-nodes 5

  # Enable autoscaling
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --enable-autoscaling --min-nodes 2 --max-nodes 10

  # Update machine type
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --machine-type pgp-large

  # Update the Kubernetes version
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --kubernetes-version 1.32.2

  # Update multiple fields
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --num-nodes 5 --enable-autohealing

To choose the next version within the Kubernetes version policy, use 'tcloud kubernetes nodepools upgrade'.

Note: Use 'tcloud kubernetes nodepools label' and 'tcloud kubernetes nodepools annotate' commands
to manage labels and annotations separately. Those values are resent with the update so they are not cleared.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		clusterFlag := ""
		if updateNodePoolCluster != nil {
			clusterFlag = *updateNodePoolCluster
		}
		clusterIdentifier, err := kuberesolve.RequireClusterRef(clusterFlag)
		if err != nil {
			return err
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

		cluster, err := FindCluster(clusters, clusterIdentifier)
		if err != nil {
			return err
		}

		// Get node pools for the cluster
		nodePools, err := client.Kubernetes().ListKubernetesNodePools(ctx, cluster.Identity, &kubernetes.ListKubernetesNodePoolsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list node pools: %w", err)
		}

		nodePool, err := FindNodePool(nodePools, *updateNodePoolName)
		if err != nil {
			return err
		}

		// Get fresh node pool data
		nodePool, err = client.Kubernetes().GetKubernetesNodePool(ctx, cluster.Identity, nodePool.Identity)
		if err != nil {
			return fmt.Errorf("failed to get node pool: %w", err)
		}

		if err := validateNodePoolUpdateFlags(cmd, nodePool); err != nil {
			return err
		}

		overrides, err := nodePoolUpdateOverrides(ctx, cmd, client)
		if err != nil {
			return err
		}

		updateReq, err := BuildNodePoolUpdateRequest(nodePool, overrides)
		if err != nil {
			return err
		}

		updatedNodePool, err := client.Kubernetes().UpdateKubernetesNodePool(ctx, cluster.Identity, nodePool.Identity, updateReq)
		if err != nil {
			return fmt.Errorf("failed to update node pool: %w", err)
		}

		if updateNodePoolWait {
			ctxWithTimeout, cancel, err := shared.WaitContext(ctx, updateNodePoolWaitTimeout)
			if err != nil {
				return err
			}
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

	updateNodePoolCluster = updateCmd.Flags().String("cluster", "", "Cluster identity, name, or slug (defaults to kubernetes.cluster in .thalassa)")
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
	updateNodePoolVersion = updateCmd.Flags().String("kubernetes-version", "", "Kubernetes version for the node pool (identity, slug, name, or version)")
	updateCmd.Flags().BoolVar(&updateNodePoolWait, "wait", false, "Wait for the node pool update to complete")
	updateCmd.Flags().DurationVar(&updateNodePoolWaitTimeout, "wait-timeout", 20*time.Minute, "Maximum time to wait for the node pool update to complete")

	_ = updateCmd.RegisterFlagCompletionFunc(ClusterFlag, completion.CompleteKubernetesCluster)
	_ = updateCmd.RegisterFlagCompletionFunc("machine-type", completion.CompleteMachineType)
	_ = updateCmd.RegisterFlagCompletionFunc("kubernetes-version", completion.CompleteKubernetesVersion)
	_ = updateCmd.RegisterFlagCompletionFunc("upgrade-strategy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"manual", "auto", "always", "on-delete", "inplace", "never"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func validateNodePoolUpdateFlags(cmd *cobra.Command, nodePool *kubernetes.KubernetesNodePool) error {
	if !nodePoolUpdateRequested(cmd) {
		return fmt.Errorf("no updates specified")
	}
	return validateNodePoolSize(cmd, nodePool, effectiveNodePoolSize(cmd, nodePool))
}

func nodePoolUpdateRequested(cmd *cobra.Command) bool {
	names := []string{
		"machine-type",
		"num-nodes",
		"enable-autoscaling",
		"min-nodes",
		"max-nodes",
		"enable-autohealing",
		"upgrade-strategy",
		"node-taints",
		"security-groups",
		"kubernetes-version",
	}
	for _, name := range names {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}

type nodePoolSize struct {
	enableAutoscaling bool
	minNodes          int
	maxNodes          int
	replicas          int
}

func effectiveNodePoolSize(cmd *cobra.Command, nodePool *kubernetes.KubernetesNodePool) nodePoolSize {
	size := nodePoolSize{
		enableAutoscaling: nodePool.EnableAutoscaling,
		minNodes:          nodePool.MinReplicas,
		maxNodes:          nodePool.MaxReplicas,
		replicas:          nodePool.Replicas,
	}
	changed := cmd.Flags().Changed
	if changed("enable-autoscaling") && updateNodePoolEnableAS != nil {
		size.enableAutoscaling = *updateNodePoolEnableAS
	}
	if changed("min-nodes") && updateNodePoolMinNodes != nil {
		size.minNodes = *updateNodePoolMinNodes
	}
	if changed("max-nodes") && updateNodePoolMaxNodes != nil {
		size.maxNodes = *updateNodePoolMaxNodes
	}
	if changed("num-nodes") && updateNodePoolReplicas != nil {
		size.replicas = *updateNodePoolReplicas
	}
	return size
}

func validateNodePoolSize(cmd *cobra.Command, nodePool *kubernetes.KubernetesNodePool, size nodePoolSize) error {
	changed := cmd.Flags().Changed
	if changed("enable-autoscaling") {
		return validateAutoscalingChange(size)
	}
	if err := validateReplicaFlag(changed, nodePool.EnableAutoscaling, size.replicas); err != nil {
		return err
	}
	if err := validateReplicaBoundFlag(changed("min-nodes"), nodePool.EnableAutoscaling, "min-nodes"); err != nil {
		return err
	}
	if err := validateReplicaBoundFlag(changed("max-nodes"), nodePool.EnableAutoscaling, "max-nodes"); err != nil {
		return err
	}
	if nodePool.EnableAutoscaling && (changed("min-nodes") || changed("max-nodes")) {
		return ValidateAutoscalingConfig(true, size.minNodes, size.maxNodes, 0)
	}
	return nil
}

func validateAutoscalingChange(size nodePoolSize) error {
	if size.enableAutoscaling {
		return ValidateAutoscalingConfig(true, size.minNodes, size.maxNodes, 0)
	}
	if size.replicas < 1 {
		return fmt.Errorf("num-nodes must be at least 1 when autoscaling is disabled")
	}
	return nil
}

func validateReplicaFlag(changed func(string) bool, autoscalingEnabled bool, replicas int) error {
	if !changed("num-nodes") {
		return nil
	}
	if autoscalingEnabled {
		return fmt.Errorf("cannot set --num-nodes when autoscaling is enabled. Use --min-nodes and --max-nodes instead, or disable autoscaling first")
	}
	if replicas < 1 {
		return fmt.Errorf("num-nodes must be at least 1")
	}
	return nil
}

func validateReplicaBoundFlag(changed bool, autoscalingEnabled bool, flagName string) error {
	if changed && !autoscalingEnabled {
		return fmt.Errorf("cannot set --%s when autoscaling is disabled. Enable autoscaling first", flagName)
	}
	return nil
}

func nodePoolUpdateOverrides(ctx context.Context, cmd *cobra.Command, client thalassa.Client) (NodePoolUpdateOverrides, error) {
	var overrides NodePoolUpdateOverrides
	changed := cmd.Flags().Changed

	if changed("machine-type") {
		if updateNodePoolMachineType == nil || *updateNodePoolMachineType == "" {
			return overrides, fmt.Errorf("machine-type is required")
		}
		resolvedMachineType, err := ResolveMachineType(ctx, client, *updateNodePoolMachineType)
		if err != nil {
			return overrides, err
		}
		overrides.MachineType = &resolvedMachineType
	}
	if changed("num-nodes") {
		overrides.Replicas = updateNodePoolReplicas
	}
	if changed("min-nodes") {
		overrides.MinReplicas = updateNodePoolMinNodes
	}
	if changed("max-nodes") {
		overrides.MaxReplicas = updateNodePoolMaxNodes
	}
	if changed("enable-autoscaling") {
		overrides.EnableAutoscaling = updateNodePoolEnableAS
	}
	if changed("enable-autohealing") {
		overrides.EnableAutoHealing = updateNodePoolEnableAH
	}
	if changed("upgrade-strategy") {
		if updateNodePoolUpgradeStrat == nil || *updateNodePoolUpgradeStrat == "" {
			return overrides, fmt.Errorf("upgrade-strategy is required")
		}
		strategy := ParseUpgradeStrategy(*updateNodePoolUpgradeStrat)
		overrides.UpgradeStrategy = &strategy
	}
	if changed("node-taints") {
		taints, err := ParseTaints(updateNodePoolTaints)
		if err != nil {
			return overrides, err
		}
		overrides.Taints = taints
		overrides.ReplaceTaints = true
	}
	if changed("security-groups") {
		overrides.SecurityGroupAttachments = append([]string(nil), updateNodePoolSecurityGroups...)
		overrides.ReplaceSecurityGroups = true
	}
	if changed("kubernetes-version") {
		if updateNodePoolVersion == nil || *updateNodePoolVersion == "" {
			return overrides, fmt.Errorf("kubernetes-version is required")
		}
		versionIdentity, err := ResolveKubernetesVersion(ctx, client, *updateNodePoolVersion)
		if err != nil {
			return overrides, err
		}
		overrides.KubernetesVersionIdentity = &versionIdentity
	}

	return overrides, nil
}

// formatReplicasForUpdate formats the replicas string for update output
func formatReplicasForUpdate(np *kubernetes.KubernetesNodePool) string {
	if np.EnableAutoscaling {
		return fmt.Sprintf("%d-%d", np.MinReplicas, np.MaxReplicas)
	}
	return fmt.Sprintf("%d", np.Replicas)
}
