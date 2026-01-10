package kubernetes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	deleteClusterWait  bool
	deleteClusterForce bool
)

var deleteCmd = &cobra.Command{
	Use:     "delete <cluster>",
	Aliases: []string{"d", "del", "remove", "rm"},
	Short:   "Delete a Kubernetes cluster",
	Long: `Delete a Kubernetes cluster and all associated resources.

This command will delete the cluster and all node pools, nodes, and other resources
associated with it. This operation cannot be undone.

Examples:
  # Delete a cluster
  tcloud kubernetes delete my-cluster

  # Delete a cluster and wait for completion
  tcloud kubernetes delete my-cluster --wait

  # Delete a cluster without confirmation
  tcloud kubernetes delete my-cluster --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		clusterIdentifier := args[0]

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
			if c.Identity == clusterIdentifier || c.Name == clusterIdentifier || c.Slug == clusterIdentifier {
				cluster = &c
				break
			}
		}

		if cluster == nil {
			return fmt.Errorf("cluster not found: %s", clusterIdentifier)
		}

		// Get fresh cluster data
		cluster, err = client.Kubernetes().GetKubernetesCluster(ctx, cluster.Identity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("cluster not found: %s", clusterIdentifier)
			}
			return fmt.Errorf("failed to get cluster: %w", err)
		}

		// Confirmation prompt
		if !deleteClusterForce {
			fmt.Printf("Do you want to delete cluster '%s'? [y/N]: ", cluster.Name)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				return nil
			}
		}

		err = client.Kubernetes().DeleteKubernetesCluster(ctx, cluster.Identity)
		if err != nil {
			return fmt.Errorf("failed to delete cluster: %w", err)
		}

		if deleteClusterWait {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Minute)
			defer cancel()

			for {
				select {
				case <-ctxWithTimeout.Done():
					return fmt.Errorf("timeout waiting for cluster to be deleted")
				default:
				}

				// Try to get the cluster - if it doesn't exist, it's deleted
				_, err := client.Kubernetes().GetKubernetesCluster(ctxWithTimeout, cluster.Identity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						// Cluster not found means it's deleted
						return nil
					}
					return fmt.Errorf("failed to get cluster status: %w", err)
				}

				time.Sleep(5 * time.Second)
			}
		}

		return nil
	},
}

func init() {
	KubernetesCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteClusterWait, "wait", false, "Wait for the cluster to be deleted before returning")
	deleteCmd.Flags().BoolVar(&deleteClusterForce, "force", false, "Skip confirmation prompt")
}
