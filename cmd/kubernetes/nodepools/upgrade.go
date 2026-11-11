package nodepools

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/kubernetes/version"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/thalassa"
)

var (
	upgradeNodePoolCluster     string
	upgradeNodePoolName        string
	upgradeNodePoolVersion     string
	upgradeAllNodePools        bool
	upgradeNodePoolDryRun      bool
	upgradeNodePoolWait        bool
	upgradeNodePoolWaitTimeout time.Duration
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade a Kubernetes node pool version",
	Long: `Upgrade the Kubernetes version of a node pool.

The node pool update API replaces the whole node pool. This command changes only the Kubernetes version and resends the current machine type, replicas, autoscaling, node settings, and security groups.

The target version follows the Kubernetes version policy: the newest patch of the current minor, otherwise the newest patch of the next minor. A node pool is not upgraded past the cluster control plane version.

Examples:
  # Show the version a node pool would move to
  tcloud kubernetes nodepools upgrade --cluster my-cluster --name worker

  # Upgrade a node pool to a specific version
  tcloud kubernetes nodepools upgrade --cluster my-cluster --name worker --version 1.32.2 --dry-run=false

  # Upgrade every node pool in a cluster
  tcloud kubernetes nodepools upgrade --cluster my-cluster --all --dry-run=false`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if upgradeAllNodePools && upgradeNodePoolName != "" {
			return fmt.Errorf("--name cannot be used with --all")
		}
		if !upgradeAllNodePools && upgradeNodePoolName == "" {
			return fmt.Errorf("--name is required")
		}

		clusterIdentifier, err := kuberesolve.RequireClusterRef(upgradeNodePoolCluster)
		if err != nil {
			return err
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusters, err := client.Kubernetes().ListKubernetesClusters(ctx, &kubernetes.ListKubernetesClustersRequest{})
		if err != nil {
			return fmt.Errorf("failed to list clusters: %w", err)
		}
		cluster, err := FindCluster(clusters, clusterIdentifier)
		if err != nil {
			return err
		}
		cluster, err = client.Kubernetes().GetKubernetesCluster(ctx, cluster.Identity)
		if err != nil {
			return fmt.Errorf("failed to get cluster: %w", err)
		}

		clusterVersion, err := semver.Parse(cluster.ClusterVersion.KubernetesVersion)
		if err != nil {
			return fmt.Errorf("invalid cluster version %q: %w", cluster.ClusterVersion.KubernetesVersion, err)
		}

		versions, err := client.Kubernetes().ListKubernetesVersions(ctx)
		if err != nil {
			return fmt.Errorf("failed to list kubernetes versions: %w", err)
		}

		nodePools, err := client.Kubernetes().ListKubernetesNodePools(ctx, cluster.Identity, &kubernetes.ListKubernetesNodePoolsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list node pools: %w", err)
		}

		selected, err := selectNodePoolsForUpgrade(nodePools, cluster.Name)
		if err != nil {
			return err
		}

		failed := 0
		for _, candidate := range selected {
			if err := upgradeNodePool(ctx, client, cluster, clusterVersion, versions, candidate.Identity); err != nil {
				if upgradeAllNodePools {
					fmt.Printf("Node pool %s: %v\n", candidate.Name, err)
					failed++
					continue
				}
				return err
			}
		}
		if failed > 0 {
			return fmt.Errorf("failed to upgrade %d node pools", failed)
		}
		return nil
	},
}

func init() {
	upgradeCmd.Flags().StringVar(&upgradeNodePoolCluster, "cluster", "", "Cluster identity, name, or slug (defaults to kubernetes.cluster in .thalassa)")
	upgradeCmd.Flags().StringVar(&upgradeNodePoolName, "name", "", "Node pool name, identity, or slug (required unless --all is set)")
	upgradeCmd.Flags().BoolVar(&upgradeAllNodePools, "all", false, "Upgrade all node pools in the cluster")
	upgradeCmd.Flags().StringVarP(&upgradeNodePoolVersion, "version", "v", "", "Version to upgrade to (identity, slug, name, or Kubernetes version). Defaults to the newest patch, or the next minor when no patch is available")
	upgradeCmd.Flags().BoolVar(&upgradeNodePoolDryRun, "dry-run", true, "Print the upgrade without changing the node pool")
	upgradeCmd.Flags().BoolVar(&upgradeNodePoolWait, "wait", false, "Wait for the node pool upgrade to complete")
	upgradeCmd.Flags().DurationVar(&upgradeNodePoolWaitTimeout, "wait-timeout", 20*time.Minute, "Maximum time to wait for the node pool upgrade to complete")

	_ = upgradeCmd.RegisterFlagCompletionFunc(ClusterFlag, completion.CompleteKubernetesCluster)
	_ = upgradeCmd.RegisterFlagCompletionFunc("name", completion.CompleteKubernetesNodePool)
	_ = upgradeCmd.RegisterFlagCompletionFunc("version", completion.CompleteKubernetesVersion)
}

func selectNodePoolsForUpgrade(nodePools []kubernetes.KubernetesNodePool, clusterName string) ([]kubernetes.KubernetesNodePool, error) {
	if upgradeAllNodePools {
		if len(nodePools) == 0 {
			return nil, fmt.Errorf("no node pools found in cluster %s", clusterName)
		}
		return nodePools, nil
	}

	nodePool, err := FindNodePool(nodePools, upgradeNodePoolName)
	if err != nil {
		return nil, err
	}
	return []kubernetes.KubernetesNodePool{*nodePool}, nil
}

func upgradeNodePool(ctx context.Context, client thalassa.Client, cluster *kubernetes.KubernetesCluster, clusterVersion semver.Version, versions []kubernetes.KubernetesVersion, nodePoolIdentity string) error {
	nodePool, err := client.Kubernetes().GetKubernetesNodePool(ctx, cluster.Identity, nodePoolIdentity)
	if err != nil {
		return fmt.Errorf("failed to get node pool: %w", err)
	}
	if nodePool.KubernetesVersion == nil || nodePool.KubernetesVersion.KubernetesVersion == "" {
		return fmt.Errorf("node pool %s has no kubernetes version", nodePool.Name)
	}

	currentVersion, err := semver.Parse(nodePool.KubernetesVersion.KubernetesVersion)
	if err != nil {
		return fmt.Errorf("invalid node pool version %q: %w", nodePool.KubernetesVersion.KubernetesVersion, err)
	}

	target, err := version.SelectNodePoolUpgradeVersion(currentVersion, clusterVersion, versions, upgradeNodePoolVersion)
	if err != nil {
		if errors.Is(err, version.ErrAlreadyAtLatestVersion) || errors.Is(err, version.ErrAlreadyAtClusterVersion) {
			fmt.Printf("Node pool %s is %s\n", nodePool.Name, err.Error())
			return nil
		}
		return err
	}
	if target.Identity == "" {
		return fmt.Errorf("selected kubernetes version %s has no identity", target.KubernetesVersion)
	}
	if target.KubernetesVersion == nodePool.KubernetesVersion.KubernetesVersion || target.Identity == nodePool.KubernetesVersion.Identity {
		fmt.Printf("Node pool %s is already at version %s\n", nodePool.Name, nodePool.KubernetesVersion.KubernetesVersion)
		return nil
	}

	if upgradeNodePoolDryRun {
		fmt.Printf("Dry run: would upgrade node pool %s from %s to %s\n", nodePool.Name, nodePool.KubernetesVersion.KubernetesVersion, target.KubernetesVersion)
		return nil
	}

	fmt.Printf("Upgrading node pool %s from %s to %s\n", nodePool.Name, nodePool.KubernetesVersion.KubernetesVersion, target.KubernetesVersion)
	versionIdentity := target.Identity
	updateReq, err := BuildNodePoolUpdateRequest(nodePool, NodePoolUpdateOverrides{
		KubernetesVersionIdentity: &versionIdentity,
	})
	if err != nil {
		return err
	}

	updatedNodePool, err := client.Kubernetes().UpdateKubernetesNodePool(ctx, cluster.Identity, nodePool.Identity, updateReq)
	if err != nil {
		return fmt.Errorf("failed to upgrade node pool: %w", err)
	}

	if upgradeNodePoolWait {
		ctxWithTimeout, cancel, err := shared.WaitContext(ctx, upgradeNodePoolWaitTimeout)
		if err != nil {
			return err
		}
		defer cancel()

		_, err = client.Kubernetes().WaitUntilKubernetesNodePoolReady(ctxWithTimeout, cluster.Identity, updatedNodePool.Identity)
		if err != nil {
			return fmt.Errorf("failed to wait for node pool to be ready: %w", err)
		}
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
			target.KubernetesVersion,
			string(updatedNodePool.Status),
			formattime.FormatTime(updatedNodePool.CreatedAt.Local(), false),
		},
	}
	table.Print([]string{"ID", "Name", "Cluster", "Replicas", "Type", "Version", "Status", "Age"}, body)
	return nil
}
