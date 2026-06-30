package registry

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/registry/namespaces"
	"github.com/thalassa-cloud/cli/cmd/registry/repositories"
)

const (
	NoHeaderKey   = "no-header"
	NamespaceFlag = "namespace"
)

// RegistryCmd represents the container registry command.
var RegistryCmd = &cobra.Command{
	Use:     "registry",
	Aliases: []string{"container-registry", "containerregistry", "cr"},
	Short:   "Manage the Thalassa container registry",
	Long:    "Manage container registry namespaces, repositories, configuration, and retention policies.",
	Example: "tcloud registry namespaces list\ntcloud registry namespaces create --namespace my-app --region eu-west-1\ntcloud registry repositories list --namespace crns-123",
}

func init() {
	RegistryCmd.AddCommand(namespaces.NamespacesCmd)
	RegistryCmd.AddCommand(repositories.RepositoriesCmd)
}
