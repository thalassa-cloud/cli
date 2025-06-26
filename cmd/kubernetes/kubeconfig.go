package kubernetes

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var (
	useClusterNameInKubeconfig bool
)

func init() {
	KubernetesKubeConfigCmd.Flags().BoolVarP(&useClusterNameInKubeconfig, "cluster-name-in-context", "n", false, "use the cluster name instead of the cluster identity in the kubeconfig context")
}

var KubernetesKubeConfigCmd = &cobra.Command{
	Use:     "kubeconfig",
	Aliases: []string{},
	Short:   "Kubernetes Kubeconfig management",

	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(args) != 1 {
			fmt.Println("must provide a cluster. Missing value <cluster>")
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
			fmt.Println("cluster not found")
			return
		}

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
