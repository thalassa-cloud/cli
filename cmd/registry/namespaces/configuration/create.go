package configuration

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/containerregistry"
)

var (
	createVisibility         string
	createRetentionEnabled   bool
	createDeleteUntagged     bool
	createRetentionDays      int
	createRetentionCount     int
	createRetentionPolicyFile string
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create namespace configuration",
	Example: "tcloud registry namespaces configuration create --namespace crns-123 --visibility private --retention-enabled --retention-days 30",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			return fmt.Errorf("--namespace is required")
		}

		retention, err := buildRetentionPolicy(createRetentionEnabled, createDeleteUntagged, createRetentionDays, createRetentionCount, createRetentionPolicyFile)
		if err != nil {
			return err
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		cfg, err := client.ContainerRegistry().CreateNamespaceConfiguration(cmd.Context(), namespace, containerregistry.CreateNamespaceConfigurationRequest{
			Visibility:      containerregistry.NamespaceVisibility(createVisibility),
			RetentionPolicy: retention,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Configuration created successfully\n")
		fmt.Printf("Visibility: %s\n", cfg.Visibility)
		return nil
	},
}

func init() {
	ConfigurationCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	createCmd.Flags().StringVar(&createVisibility, "visibility", string(containerregistry.NamespaceVisibilityPrivate), "Namespace visibility")
	createCmd.Flags().BoolVar(&createRetentionEnabled, "retention-enabled", false, "Enable retention policy")
	createCmd.Flags().BoolVar(&createDeleteUntagged, "delete-untagged", false, "Delete untagged images during retention runs")
	createCmd.Flags().IntVar(&createRetentionDays, "retention-days", 0, "Retain tags for this many days (simple single-rule policy)")
	createCmd.Flags().IntVar(&createRetentionCount, "retention-count", 0, "Retain this many recent tags (simple single-rule policy)")
	createCmd.Flags().StringVar(&createRetentionPolicyFile, "retention-policy-file", "", "Path to JSON file with full retention policy")

	createCmd.MarkFlagRequired(NamespaceFlag)
	createCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
}
