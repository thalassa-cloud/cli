package repositories

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/client-go/thalassa"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:     "delete REPOSITORY [REPOSITORY...]",
	Short:   "Delete repositories and all artifacts",
	Long:    "Permanently delete repositories and all contained artifacts.",
	Example: "tcloud registry repositories delete --namespace crns-123 repo-456 --force",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRepositoryActions(
			cmd, args, deleteForce,
			fmt.Sprintf("Delete %d repository(ies) and all artifacts from namespace %s?\n", len(args), namespace),
			func(repoID string) { fmt.Printf("Deleting repository: %s\n", repoID) },
			func(repoID string) { fmt.Printf("Repository %s not found\n", repoID) },
			func(repoID string) { fmt.Printf("Repository %s deleted successfully\n", repoID) },
			"failed to delete repository",
			func(ctx context.Context, client thalassa.Client, namespace, repoID string) error {
				return client.ContainerRegistry().DeleteContainerRegistryRepositoryWithAllArtifacts(ctx, namespace, repoID)
			},
		)
	},
}

func init() {
	RepositoriesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation")
	_ = deleteCmd.MarkFlagRequired(NamespaceFlag)
	deleteCmd.ValidArgsFunction = completeRepositoryID
	_ = deleteCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
}
