package connect

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/config/securefile"
	"github.com/thalassa-cloud/cli/internal/fzf"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	kubeconfigrender "github.com/thalassa-cloud/cli/internal/kubernetes/kubeconfig"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	useTempKubeconfig bool
	kubeconfigPath    string
	inlineToken       bool
)

var KubernetesConnectCmd = &cobra.Command{
	Use:               "connect",
	Aliases:           []string{"connection", "shell", "c"},
	Short:             "Connect your shell to the Kubernetes Cluster",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completion.CompleteKubernetesCluster,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return err
		}

		clusterRef, err := getSelectedCluster(args)
		if err != nil {
			return err
		}

		cluster, err := kuberesolve.ResolveKubernetesClusterRef(ctx, client.Kubernetes(), clusterRef)
		if err != nil {
			return err
		}

		fmt.Printf("Getting kubeconfig for cluster %s\n", cluster.Name)
		session, err := client.Kubernetes().GetKubernetesClusterKubeconfig(ctx, cluster.Identity)
		if err != nil {
			return err
		}

		authMode := kubeconfigrender.AuthModeExec
		if inlineToken {
			authMode = kubeconfigrender.AuthModeToken
		}

		config, err := kubeconfigrender.Render(kubeconfigrender.NewRenderInput(cluster, session), authMode)
		if err != nil {
			return err
		}

		path, cleanup, err := prepareKubeconfigFile(cluster.Slug, config)
		if err != nil {
			return err
		}
		if cleanup != nil {
			defer cleanup()
		}

		if err := os.Setenv("KUBECONFIG", path); err != nil {
			return fmt.Errorf("set KUBECONFIG: %w", err)
		}
		defer func() { _ = os.Unsetenv("KUBECONFIG") }()

		if err := os.Setenv("TCLOUD_CLUSTER_IDENTITY", cluster.Identity); err != nil {
			return fmt.Errorf("set TCLOUD_CLUSTER_IDENTITY: %w", err)
		}
		if err := os.Setenv("TCLOUD_CLUSTER_NAME", cluster.Name); err != nil {
			return fmt.Errorf("set TCLOUD_CLUSTER_NAME: %w", err)
		}
		if err := os.Setenv("TCLOUD_CLUSTER_SLUG", cluster.Slug); err != nil {
			return fmt.Errorf("set TCLOUD_CLUSTER_SLUG: %w", err)
		}

		fmt.Println("Connecting..")
		if !useTempKubeconfig {
			fmt.Printf("Using kubeconfig at %s\n", path)
		}

		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}

		subcmd := exec.Command(shell)
		subcmd.Stdin = os.Stdin
		subcmd.Stdout = os.Stdout
		subcmd.Stderr = os.Stderr
		if err := subcmd.Run(); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				os.Exit(exitErr.ExitCode())
			}
			return err
		}

		fmt.Printf("Disconnected\n")
		return nil
	},
}

func init() {
	KubernetesConnectCmd.Flags().BoolVar(&useTempKubeconfig, "temp", true, "use a temporary kubeconfig file that is removed when the shell exits")
	KubernetesConnectCmd.Flags().StringVar(&kubeconfigPath, "kubeconfig-path", "", "path to write the kubeconfig when --temp=false")
	KubernetesConnectCmd.Flags().BoolVar(&inlineToken, "inline-token", false, "embed the session token in the kubeconfig instead of using a kubectl exec credential plugin")
}

func prepareKubeconfigFile(clusterSlug, config string) (path string, cleanup func(), err error) {
	if useTempKubeconfig {
		tmpFile, err := os.CreateTemp("", "kubeconfig-*.yaml")
		if err != nil {
			return "", nil, fmt.Errorf("create temp kubeconfig: %w", err)
		}
		if err := tmpFile.Close(); err != nil {
			_ = os.Remove(tmpFile.Name())
			return "", nil, fmt.Errorf("close temp kubeconfig: %w", err)
		}
		if err := securefile.Write(tmpFile.Name(), []byte(config)); err != nil {
			_ = os.Remove(tmpFile.Name())
			return "", nil, err
		}
		return tmpFile.Name(), func() { _ = os.Remove(tmpFile.Name()) }, nil
	}

	path = kubeconfigPath
	if path == "" {
		home, err := homedir.Dir()
		if err != nil {
			return "", nil, fmt.Errorf("resolve home directory: %w", err)
		}
		path = filepath.Join(home, ".kube", fmt.Sprintf("tcloud-%s.yaml", clusterSlug))
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", nil, fmt.Errorf("create kubeconfig directory: %w", err)
	}
	if err := securefile.Write(path, []byte(config)); err != nil {
		return "", nil, err
	}
	return path, nil, nil
}

func getSelectedCluster(args []string) (string, error) {
	if len(args) == 0 && fzf.IsInteractiveMode(os.Stdout) {
		command := fmt.Sprintf("%s kubernetes clusters --no-header", os.Args[0])
		return fzf.InteractiveChoice(command)
	}
	if len(args) == 1 {
		return args[0], nil
	}
	if contextstate.OrganisationFlag != "" {
		return "", errors.New("must provide a cluster when organisation is set via flag")
	}
	return "", errors.New("must provide a cluster")
}
