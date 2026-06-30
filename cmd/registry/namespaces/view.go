package namespaces

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"gopkg.in/yaml.v3"
)

var outputFormat string

var viewCmd = &cobra.Command{
	Use:     "view NAMESPACE",
	Short:   "View namespace details",
	Aliases: []string{"show", "get", "describe"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ns, err := client.ContainerRegistry().GetContainerRegistryNamespace(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get namespace: %w", err)
		}

		if outputFormat == "yaml" {
			ns.Organisation = nil
			yamlData, err := yaml.Marshal(ns)
			if err != nil {
				return fmt.Errorf("failed to marshal to YAML: %w", err)
			}
			fmt.Print(string(yamlData))
			return nil
		}

		fmt.Printf("Namespace Details:\n")
		fmt.Printf("  ID: %s\n", ns.Identity)
		fmt.Printf("  Namespace: %s\n", ns.Namespace)
		fmt.Printf("  Description: %s\n", ns.Description)
		fmt.Printf("  Total Size: %s\n", formatBytes(ns.TotalSizeBytes))
		fmt.Printf("  Repositories: %d\n", len(ns.Repositories))
		if ns.Region != nil {
			fmt.Printf("  Region: %s (%s)\n", ns.Region.Name, ns.Region.Identity)
		}
		if ns.Configuration != nil {
			fmt.Printf("  Visibility: %s\n", ns.Configuration.Visibility)
			if ns.Configuration.RetentionPolicy != nil {
				fmt.Printf("  Retention Policy: enabled=%t\n", ns.Configuration.RetentionPolicy.Enabled)
			}
		}
		fmt.Printf("  Created: %s\n", formattime.FormatTime(ns.CreatedAt.Local(), false))
		fmt.Printf("  Updated: %s\n", formattime.FormatTime(ns.UpdatedAt.Local(), false))
		return nil
	},
}

func init() {
	NamespacesCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml)")
	viewCmd.RegisterFlagCompletionFunc("output", completeOutputFormat)
	viewCmd.ValidArgsFunction = completeNamespaceID
}
