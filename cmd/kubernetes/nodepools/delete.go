package nodepools

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var (
	deleteNodePoolCluster     string
	deleteNodePoolId          string
	deleteNodePoolWait        bool
	deleteNodePoolWaitTimeout time.Duration
	deleteNodePoolForce       bool
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d", "del", "remove", "rm"},
	Short:   "Delete a Kubernetes node pool",
	Long: `Delete a node pool from a Kubernetes cluster.

This command will delete the node pool and all nodes associated with it.
The cluster must be in a ready state to delete node pools.

Examples:
  # Delete a node pool
  tcloud kubernetes nodepools delete --cluster my-cluster --nodepool worker-pool

  # Delete a node pool and wait for completion
  tcloud kubernetes nodepools delete --cluster my-cluster --nodepool worker-pool --wait

  # Delete a node pool and wait up to 45 minutes
  tcloud kubernetes nodepools delete --cluster my-cluster --nodepool worker-pool --wait --wait-timeout 45m

  # Delete a node pool without confirmation
  tcloud kubernetes nodepools delete --cluster my-cluster --nodepool worker-pool --force`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if deleteNodePoolCluster == "" {
			return fmt.Errorf("--cluster is required")
		}
		if deleteNodePoolId == "" {
			return fmt.Errorf("--nodepool is required")
		}

		clusterIdentifier := deleteNodePoolCluster
		nodePoolIdentifier := deleteNodePoolId

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

		// Get node pools for the cluster
		nodePools, err := client.Kubernetes().ListKubernetesNodePools(ctx, cluster.Identity, &kubernetes.ListKubernetesNodePoolsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list node pools: %w", err)
		}

		// Find the node pool
		var nodePool *kubernetes.KubernetesNodePool
		for _, np := range nodePools {
			if np.Identity == nodePoolIdentifier || np.Name == nodePoolIdentifier || np.Slug == nodePoolIdentifier {
				nodePool = &np
				break
			}
		}

		if nodePool == nil {
			return fmt.Errorf("node pool not found: %s", nodePoolIdentifier)
		}

		// Confirmation prompt
		proceed, err := shared.PromptYesNoUnlessForce(deleteNodePoolForce, fmt.Sprintf("Do you want to delete node pool '%s'? [y/N]: ", nodePool.Name))
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}

		err = client.Kubernetes().DeleteKubernetesNodePool(ctx, cluster.Identity, nodePool.Identity)
		if err != nil {
			return fmt.Errorf("failed to delete node pool: %w", err)
		}

		if deleteNodePoolWait {
			ctxWithTimeout, cancel, err := shared.WaitContext(ctx, deleteNodePoolWaitTimeout)
			if err != nil {
				return err
			}
			defer cancel()

			if err := client.Kubernetes().WaitUntilKubernetesNodePoolDeleted(ctxWithTimeout, cluster.Identity, nodePool.Identity); err != nil {
				return fmt.Errorf("failed to wait for node pool to be deleted: %w", err)
			}
		}

		return nil
	},
}

func init() {
	// Command is registered in kubernetesclusters.go

	deleteCmd.Flags().StringVar(&deleteNodePoolCluster, "cluster", "", "Cluster identity, name, or slug (required)")
	deleteCmd.Flags().StringVar(&deleteNodePoolId, "nodepool", "", "Node pool name, identity, or slug (required)")
	deleteCmd.Flags().BoolVar(&deleteNodePoolWait, "wait", false, "Wait for the node pool to be deleted before returning")
	deleteCmd.Flags().DurationVar(&deleteNodePoolWaitTimeout, "wait-timeout", 20*time.Minute, "Maximum time to wait for the node pool to be deleted")
	deleteCmd.Flags().BoolVar(&deleteNodePoolForce, "force", false, "Skip confirmation prompt")

	_ = deleteCmd.RegisterFlagCompletionFunc(ClusterFlag, completion.CompleteKubernetesCluster)
	_ = deleteCmd.RegisterFlagCompletionFunc("nodepool", completion.CompleteKubernetesNodePool)
}
