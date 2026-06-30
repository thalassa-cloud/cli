package repositories

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var deleteArtifactsForce bool

var deleteArtifactsCmd = &cobra.Command{
	Use:     "delete-artifacts REPOSITORY [REPOSITORY...]",
	Short:   "Delete artifacts from repositories",
	Long:    "Request deletion of artifacts from repositories without deleting the repository itself.",
	Example: "tcloud registry repositories delete-artifacts --namespace crns-123 repo-456 --force",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			return fmt.Errorf("--namespace is required")
		}

		if !deleteArtifactsForce {
			fmt.Printf("Delete artifacts from %d repository(ies) in namespace %s?\n", len(args), namespace)
			var confirm string
			fmt.Print("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		for _, repoID := range args {
			fmt.Printf("Deleting artifacts from repository: %s\n", repoID)
			if err := client.ContainerRegistry().DeleteContainerRegistryRepositoryArtifact(cmd.Context(), namespace, repoID); err != nil {
				if tcclient.IsNotFound(err) {
					fmt.Printf("Repository %s not found\n", repoID)
					continue
				}
				return fmt.Errorf("failed to delete artifacts: %w", err)
			}
			fmt.Printf("Artifact deletion requested for repository %s\n", repoID)
		}
		return nil
	},
}

func init() {
	RepositoriesCmd.AddCommand(deleteArtifactsCmd)
	deleteArtifactsCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	deleteArtifactsCmd.Flags().BoolVar(&deleteArtifactsForce, "force", false, "Skip confirmation")
	deleteArtifactsCmd.MarkFlagRequired(NamespaceFlag)
	deleteArtifactsCmd.ValidArgsFunction = completeRepositoryID
	deleteArtifactsCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
}
