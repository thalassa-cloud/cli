package repositories

import "github.com/spf13/cobra"

const (
	NoHeaderKey   = "no-header"
	NamespaceFlag = "namespace"
)

// RepositoriesCmd manages container registry repositories.
var RepositoriesCmd = &cobra.Command{
	Use:     "repositories",
	Aliases: []string{"repository", "repos", "repo"},
	Short:   "Manage container registry repositories",
	Long:    "List and manage repositories within a container registry namespace. All commands require --namespace.",
	Example: "tcloud registry repositories list --namespace crns-123\ntcloud registry repositories delete --namespace crns-123 repo-456",
}
