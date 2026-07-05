package kubernetes

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	kubeconfigrender "github.com/thalassa-cloud/cli/internal/kubernetes/kubeconfig"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var kubeconfigInlineToken bool

var KubernetesKubeConfigCmd = &cobra.Command{
	Use:               "kubeconfig",
	Short:             "Print a kubeconfig for a Kubernetes cluster",
	ValidArgsFunction: completion.CompleteKubernetesCluster,
	Args:              cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return err
		}

		cluster, err := kuberesolve.ResolveKubernetesClusterRef(ctx, client.Kubernetes(), args[0])
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Getting kubeconfig for cluster %s\n", cluster.Name)
		session, err := client.Kubernetes().GetKubernetesClusterKubeconfig(ctx, cluster.Identity)
		if err != nil {
			return err
		}

		authMode := kubeconfigrender.AuthModeExec
		if kubeconfigInlineToken {
			authMode = kubeconfigrender.AuthModeToken
		}

		config, err := kubeconfigrender.Render(kubeconfigrender.NewRenderInput(cluster, session), authMode)
		if err != nil {
			return err
		}

		fmt.Print(config)
		return nil
	},
}

func init() {
	KubernetesCmd.AddCommand(KubernetesKubeConfigCmd)
	KubernetesKubeConfigCmd.Flags().BoolVar(&kubeconfigInlineToken, "inline-token", false, "embed the session token in the kubeconfig instead of using a kubectl exec credential plugin")
}
