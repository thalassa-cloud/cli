package repositories

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/client-go/thalassa"
)

var deleteArtifactsForce bool

var deleteArtifactsCmd = &cobra.Command{
	Use:     "delete-artifacts REPOSITORY [REPOSITORY...]",
	Short:   "Delete artifacts from repositories",
	Long:    "Request deletion of artifacts from repositories without deleting the repository itself.",
	Example: "tcloud registry repositories delete-artifacts --namespace crns-123 repo-456 --force",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRepositoryActions(
			cmd, args, deleteArtifactsForce,
			fmt.Sprintf("Delete artifacts from %d repository(ies) in namespace %s?\n", len(args), namespace),
			func(repoID string) { fmt.Printf("Deleting artifacts from repository: %s\n", repoID) },
			func(repoID string) { fmt.Printf("Repository %s not found\n", repoID) },
			func(repoID string) { fmt.Printf("Artifact deletion requested for repository %s\n", repoID) },
			"failed to delete artifacts",
			func(ctx context.Context, client thalassa.Client, namespace, repoID string) error {
				return client.ContainerRegistry().DeleteContainerRegistryRepositoryArtifact(ctx, namespace, repoID)
			},
		)
	},
}

func init() {
	RepositoriesCmd.AddCommand(deleteArtifactsCmd)
	deleteArtifactsCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	deleteArtifactsCmd.Flags().BoolVar(&deleteArtifactsForce, "force", false, "Skip confirmation")
	_ = deleteArtifactsCmd.MarkFlagRequired(NamespaceFlag)
	deleteArtifactsCmd.ValidArgsFunction = completeRepositoryID
	_ = deleteArtifactsCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
}
