package kubernetes

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var KubernetesKubeConfigCmd = &cobra.Command{
	Use:               "kubeconfig",
	Aliases:           []string{},
	Short:             "Kubernetes Kubeconfig management",
	ValidArgsFunction: completion.CompleteKubernetesCluster,
	Args:              cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "must provide a cluster. Missing value <cluster>")
			return
		}

		clusterIdentity := args[0]
		// get the cluster
		cluster, err := client.Kubernetes().GetKubernetesCluster(ctx, clusterIdentity)
		if err != nil {
			// try and find the cluster by name or slug
			clusters, err := client.Kubernetes().ListKubernetesClusters(ctx, &kubernetes.ListKubernetesClustersRequest{})
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, potentialCluster := range clusters {
				if potentialCluster.Name == clusterIdentity || potentialCluster.Slug == clusterIdentity {
					cluster = &potentialCluster
					break
				}
			}
		}
		if cluster == nil {
			fmt.Fprintln(os.Stderr, "cluster not found")
			return
		}

		fmt.Fprintf(os.Stderr, "Getting kubeconfig for cluster %s\n", cluster.Name)
		session, err := client.Kubernetes().GetKubernetesClusterKubeconfig(ctx, cluster.Identity)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(session.Kubeconfig)
	},
}

func init() {
	KubernetesCmd.AddCommand(KubernetesKubeConfigCmd)
}
