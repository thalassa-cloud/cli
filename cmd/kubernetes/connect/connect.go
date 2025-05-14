package connect

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/fzf"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	//go:embed templates/kubeconfig.yaml
	kubeconfigTemplate string
)

var (
	useTempKubeconfig  bool
	tempKubeConfigPath string
)

type KubeconfigTemplateInput struct {
	APIServerURL        string
	Base64CACertificate string
	SessionToken        string
	Cluster             string
	User                string
}

var KubernetesConnectCmd = &cobra.Command{
	Use:     "connect",
	Aliases: []string{"connection", "shell", "c"},
	Short:   "Connect your shell to the Kubernetes Cluster",

	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		clusterIdentity, err := getSelectedCluster(args)
		if err != nil {
			fmt.Println(err)
			return
		}

		// get the cluster
		cluster, err := client.Kubernetes().GetKubernetesCluster(ctx, clusterIdentity)
		if err != nil {
			// try and find the cluster by name or slug
			clusters, err := client.Kubernetes().ListKubernetesClusters(ctx)
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

		fmt.Printf("Getting kubeconfig for cluster %s\n", cluster.Name)
		session, err := client.Kubernetes().GetKubernetesClusterKubeconfig(ctx, cluster.Identity)
		if err != nil {
			fmt.Println(err)
			return
		}

		i := KubeconfigTemplateInput{
			APIServerURL:        cluster.APIServerURL,
			Base64CACertificate: session.CACertificate,
			SessionToken:        session.Token,
			Cluster:             cluster.Slug,
			User:                session.Username,
		}

		tmpl, err := template.New("kubeconfig").Parse(kubeconfigTemplate)
		if err != nil {
			fmt.Println(err)
			return
		}

		w := new(strings.Builder)
		if err := tmpl.Execute(w, i); err != nil {
			fmt.Println("failed to render kubeconfig template")
			return
		}

		if useTempKubeconfig {
			// with session.Kubeconfig you can now connect to the cluster. We will create a temp file and configure the shell to use it
			//
			// 1. create a temp file
			// 2. write session.Kubeconfig to the file
			// 3. export KUBECONFIG to the file
			// 4. remove the file when the shell exits
			//
			// Step 1: create a temp file with secure permissions
			tmpFile, err := os.CreateTemp("", "kubeconfig-*.yaml")
			if err != nil {
				fmt.Println("Error creating temp file:", err)
				return
			}
			defer os.Remove(tmpFile.Name()) // clean up

			// Set secure permissions on the temp file
			if err := tmpFile.Chmod(0600); err != nil {
				fmt.Println("Error setting file permissions:", err)
				return
			}

			// Step 2: write session.Kubeconfig to the file
			if _, err := tmpFile.Write([]byte(w.String())); err != nil {
				fmt.Println("Error writing to temp file:", err)
				return
			}
			if err := tmpFile.Close(); err != nil {
				fmt.Println("Error closing temp file:", err)
				return
			}

			// Step 3: export KUBECONFIG to the file
			if err := os.Setenv("KUBECONFIG", tmpFile.Name()); err != nil {
				fmt.Println("Error setting KUBECONFIG environment variable:", err)
				return
			}

			// set the TCLOUD_CLUSTER_ID for the shell
			if err := os.Setenv("TCLOUD_CLUSTER_IDENTITY", cluster.Identity); err != nil {
				fmt.Println("Error setting TCLOUD_CLUSTER_IDENTITY environment variable:", err)
				return
			}

			if err := os.Setenv("TCLOUD_CLUSTER_NAME", cluster.Name); err != nil {
				fmt.Println("Error setting TCLOUD_CLUSTER_NAME environment variable:", err)
				return
			}

			if err := os.Setenv("TCLOUD_CLUSTER_SLUG", cluster.Slug); err != nil {
				fmt.Println("Error setting TCLOUD_CLUSTER_SLUG environment variable:", err)
				return
			}

			// Inform the user
			// fmt.Printf("KUBECONFIG set to %s\n", tmpFile.Name())
			fmt.Println("Connecting..")

			// start the subshell
			shell := os.Getenv("SHELL")
			if shell == "" {
				shell = "/bin/sh"
			}

			// Start the subshell
			subcmd := exec.Command(shell)
			subcmd.Stdin = os.Stdin
			subcmd.Stdout = os.Stdout
			subcmd.Stderr = os.Stderr

			if err := subcmd.Run(); err != nil {
				// ignore the error for now
			}

			// Step 4: remove the file when the shell exits
			if err := os.Unsetenv("KUBECONFIG"); err != nil {
				fmt.Println("Error unsetting KUBECONFIG environment variable:", err)
				return
			}
			fmt.Printf("Disconnected\n")
		} else {
			fmt.Println("Not yet implemented")
			os.Exit(1)
			return
		}
	},
}

func init() {
	KubernetesConnectCmd.Flags().BoolVar(&useTempKubeconfig, "temp", true, "use a temporary kubeconfig file")
	KubernetesConnectCmd.Flags().StringVar(&tempKubeConfigPath, "kubeconfig-path", "", "path to the kubeconfig file")
}

func getSelectedCluster(args []string) (string, error) {
	if contextstate.OrganisationFlag != "" {
		return contextstate.OrganisationFlag, nil
	}

	if len(args) == 0 && fzf.IsInteractiveMode(os.Stdout) {
		command := fmt.Sprintf("%s kubernetes clusters --no-header", os.Args[0])
		return fzf.InteractiveChoice(command)
	} else if len(args) == 1 {
		return args[0], nil
	} else {
		return "", errors.New("invalid organisation")
	}
}
