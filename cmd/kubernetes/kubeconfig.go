package kubernetes

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	kubeconfigrender "github.com/thalassa-cloud/cli/internal/kubernetes/kubeconfig"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var (
	kubeconfigInlineToken     bool
	kubeconfigSessionLifetime time.Duration
)

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

		if !kubeconfigInlineToken {
			// expire the session token immediately, we don't use it in this call, it's fetched by the credentials exec plugin
			kubeconfigSessionLifetime = 1 * time.Second
		}

		fmt.Fprintf(os.Stderr, "Getting kubeconfig for cluster %s\n", cluster.Name)
		session, err := client.Kubernetes().GetKubernetesClusterKubeconfigWithParams(ctx, cluster.Identity, kubernetes.KubeconfigParams{
			SessionLifetime: kubeconfigSessionLifetime,
		})
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
	KubernetesKubeConfigCmd.Flags().DurationVar(&kubeconfigSessionLifetime, "session-lifetime", 4*7*24*time.Hour, "Lifetime of the kubeconfig session token. Defaults to 4 weeks.")
}
