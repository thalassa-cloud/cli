package repositories

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
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
		if namespace == "" {
			return fmt.Errorf("--namespace is required")
		}

		if !deleteForce {
			fmt.Printf("Delete %d repository(ies) and all artifacts from namespace %s?\n", len(args), namespace)
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
			fmt.Printf("Deleting repository: %s\n", repoID)
			if err := client.ContainerRegistry().DeleteContainerRegistryRepositoryWithAllArtifacts(cmd.Context(), namespace, repoID); err != nil {
				if tcclient.IsNotFound(err) {
					fmt.Printf("Repository %s not found\n", repoID)
					continue
				}
				return fmt.Errorf("failed to delete repository: %w", err)
			}
			fmt.Printf("Repository %s deleted successfully\n", repoID)
		}
		return nil
	},
}

func init() {
	RepositoriesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation")
	deleteCmd.MarkFlagRequired(NamespaceFlag)
	deleteCmd.ValidArgsFunction = completeRepositoryID
	deleteCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
}
