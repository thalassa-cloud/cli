package repositories

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/containerregistry"
	"github.com/thalassa-cloud/client-go/filters"
)

var (
	noHeader          bool
	showExactTime     bool
	listLabelSelector string
	namespace         string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List repositories in a namespace",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			return fmt.Errorf("--namespace is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := filters.Filters{}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		repos, err := client.ContainerRegistry().ListContainerRegistryRepositories(cmd.Context(), namespace, &containerregistry.ListContainerRegistryRepositoriesRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(repos))
		for _, repo := range repos {
			lastPushed := "-"
			if repo.LastPushedAt != nil {
				lastPushed = formattime.FormatTime(repo.LastPushedAt.Local(), showExactTime)
			}
			body = append(body, []string{
				repo.Identity,
				repo.Image,
				repo.FullName,
				fmt.Sprintf("%d", repo.TagCount),
				fmt.Sprintf("%d", repo.ArtifactCount),
				formatSizeMB(repo.TotalSizeBytes),
				lastPushed,
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Image", "Full Name", "Tags", "Artifacts", "Size", "Last Pushed"}, body)
		}
		return nil
	},
}

func init() {
	RepositoriesCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector (format: key1=value1,key2=value2)")

	listCmd.MarkFlagRequired(NamespaceFlag)
	listCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
}

func formatSizeMB(bytes int64) string {
	if bytes == 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f MB", float64(bytes)/(1024*1024))
}
