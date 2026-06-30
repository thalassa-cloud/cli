package repositories

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"gopkg.in/yaml.v3"
)

var outputFormat string

var viewCmd = &cobra.Command{
	Use:     "view REPOSITORY",
	Short:   "View repository details",
	Aliases: []string{"show", "get", "describe"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			return fmt.Errorf("--namespace is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		repo, err := client.ContainerRegistry().GetContainerRegistryRepository(cmd.Context(), namespace, args[0])
		if err != nil {
			return fmt.Errorf("failed to get repository: %w", err)
		}

		if outputFormat == "yaml" {
			yamlData, err := yaml.Marshal(repo)
			if err != nil {
				return fmt.Errorf("failed to marshal to YAML: %w", err)
			}
			fmt.Print(string(yamlData))
			return nil
		}

		fmt.Printf("Repository Details:\n")
		fmt.Printf("  ID: %s\n", repo.Identity)
		fmt.Printf("  Image: %s\n", repo.Image)
		fmt.Printf("  Full Name: %s\n", repo.FullName)
		fmt.Printf("  Description: %s\n", repo.Description)
		fmt.Printf("  Tags: %d\n", repo.TagCount)
		fmt.Printf("  Artifacts: %d\n", repo.ArtifactCount)
		fmt.Printf("  Size: %s\n", formatSizeMB(repo.TotalSizeBytes))
		if repo.LastPushedAt != nil {
			fmt.Printf("  Last Pushed: %s\n", formattime.FormatTime(repo.LastPushedAt.Local(), false))
		}
		if repo.LastPulledAt != nil {
			fmt.Printf("  Last Pulled: %s\n", formattime.FormatTime(repo.LastPulledAt.Local(), false))
		}
		if len(repo.Tags) > 0 {
			fmt.Printf("  Tags:\n")
			for _, tag := range repo.Tags {
				fmt.Printf("    - %s (%s, %.2f MB)\n", tag.Tag, tag.Sha256, tag.SizeMb)
			}
		}
		return nil
	},
}

func init() {
	RepositoriesCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	viewCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml)")
	viewCmd.MarkFlagRequired(NamespaceFlag)
	viewCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
	viewCmd.RegisterFlagCompletionFunc("output", completeOutputFormat)
	viewCmd.ValidArgsFunction = completeRepositoryID
}
