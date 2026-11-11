package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/registry"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/containerregistry"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/thalassa"
)

var (
	noHeader          bool
	showExactTime     bool
	listLabelSelector string
	namespace         string
	namespaceID       string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List container registry repositories",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		targets, err := namespacesForList(cmd.Context(), client)
		if err != nil {
			return err
		}

		f := filters.Filters{}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		body := make([][]string, 0)
		for _, target := range targets {
			repos, err := registry.ListRepositories(cmd.Context(), client.ContainerRegistry(), target.Identity, f)
			if err != nil {
				return fmt.Errorf("failed to list repositories in namespace %s: %w", target.label(), err)
			}
			for _, repo := range repos {
				body = append(body, []string{
					repo.Identity,
					repo.Image,
					fmt.Sprintf("%d", repo.TagCount),
					fmt.Sprintf("%d", repo.ArtifactCount),
					formatRepositorySize(repo.TotalSizeBytes),
					formatOptionalTime(repo.LastPushedAt, showExactTime),
					formatOptionalTime(repo.LastPulledAt, showExactTime),
				})
			}
		}

		if len(body) == 0 {
			fmt.Println("No repositories found")
			return nil
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Image", "Tags", "Artifacts", "Size", "Last Pushed", "Last Pulled"}, body)
		}
		return nil
	},
}

func init() {
	RepositoriesCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace name. Omit to list repositories in every namespace")
	listCmd.Flags().StringVar(&namespaceID, "id", "", "Namespace identity")
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector (format: key1=value1,key2=value2)")

	listCmd.MarkFlagsMutuallyExclusive(NamespaceFlag, "id")
	_ = listCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceName)
	_ = listCmd.RegisterFlagCompletionFunc("id", completeNamespaceID)
}

type namespaceListTarget struct {
	Identity string
	Name     string
}

func (t namespaceListTarget) label() string {
	if t.Name != "" {
		return t.Name
	}
	return t.Identity
}

func namespacesForList(ctx context.Context, client thalassa.Client) ([]namespaceListTarget, error) {
	if namespaceID != "" {
		return []namespaceListTarget{{Identity: namespaceID}}, nil
	}

	namespaces, err := client.ContainerRegistry().ListContainerRegistryNamespaces(ctx, &containerregistry.ListContainerRegistryNamespacesRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}
	if namespace != "" {
		ns, err := ResolveNamespaceByName(namespaces, namespace)
		if err != nil {
			return nil, err
		}
		return []namespaceListTarget{{Identity: ns.Identity, Name: ns.Namespace}}, nil
	}

	targets := make([]namespaceListTarget, 0, len(namespaces))
	for _, ns := range namespaces {
		if ns.Identity == "" {
			continue
		}
		targets = append(targets, namespaceListTarget{Identity: ns.Identity, Name: ns.Namespace})
	}
	return targets, nil
}

func formatOptionalTime(value *time.Time, exact bool) string {
	if value == nil || value.IsZero() {
		return "-"
	}
	return formattime.FormatTime(value.Local(), exact)
}

func formatRepositorySize(size int64) string {
	if size <= 0 {
		return "-"
	}
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit && exp < len("KMGTPE")-1; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
