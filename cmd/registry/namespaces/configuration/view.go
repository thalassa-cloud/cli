package configuration

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"gopkg.in/yaml.v3"
)

var (
	namespace    string
	outputFormat string
)

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View namespace configuration",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			return fmt.Errorf("--namespace is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		cfg, err := client.ContainerRegistry().GetNamespaceConfiguration(cmd.Context(), namespace)
		if err != nil {
			return fmt.Errorf("failed to get configuration: %w", err)
		}

		if outputFormat == "yaml" {
			yamlData, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("failed to marshal to YAML: %w", err)
			}
			fmt.Print(string(yamlData))
			return nil
		}

		fmt.Printf("Configuration:\n")
		fmt.Printf("  Visibility: %s\n", cfg.Visibility)
		if cfg.RetentionPolicy != nil {
			fmt.Printf("  Retention Enabled: %t\n", cfg.RetentionPolicy.Enabled)
			fmt.Printf("  Delete Untagged: %t\n", cfg.RetentionPolicy.DeleteUntaggedImages)
			fmt.Printf("  Rules: %d\n", len(cfg.RetentionPolicy.Rules))
		}
		return nil
	},
}

func init() {
	ConfigurationCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	viewCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml)")
	viewCmd.MarkFlagRequired(NamespaceFlag)
	viewCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
}
