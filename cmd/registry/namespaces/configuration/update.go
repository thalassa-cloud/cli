package configuration

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/containerregistry"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateVisibility         string
	updateRetentionEnabled   bool
	updateDeleteUntagged     bool
	updateRetentionDays      int
	updateRetentionCount     int
	updateRetentionPolicyFile string
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update namespace configuration",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			return fmt.Errorf("--namespace is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.ContainerRegistry().GetNamespaceConfiguration(cmd.Context(), namespace)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("configuration not found for namespace: %s", namespace)
			}
			return fmt.Errorf("failed to get configuration: %w", err)
		}

		req := containerregistry.UpdateNamespaceConfigurationRequest{
			Visibility:      current.Visibility,
			RetentionPolicy: current.RetentionPolicy,
		}

		if cmd.Flags().Changed("visibility") {
			req.Visibility = containerregistry.NamespaceVisibility(updateVisibility)
		}
		if cmd.Flags().Changed("retention-enabled") || cmd.Flags().Changed("delete-untagged") ||
			cmd.Flags().Changed("retention-days") || cmd.Flags().Changed("retention-count") ||
			cmd.Flags().Changed("retention-policy-file") {
			retention, err := buildRetentionPolicy(updateRetentionEnabled, updateDeleteUntagged, updateRetentionDays, updateRetentionCount, updateRetentionPolicyFile)
			if err != nil {
				return err
			}
			req.RetentionPolicy = retention
		}

		cfg, err := client.ContainerRegistry().UpdateNamespaceConfiguration(cmd.Context(), namespace, req)
		if err != nil {
			return err
		}

		fmt.Printf("Configuration updated successfully\n")
		fmt.Printf("Visibility: %s\n", cfg.Visibility)
		return nil
	},
}

func init() {
	ConfigurationCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	updateCmd.Flags().StringVar(&updateVisibility, "visibility", string(containerregistry.NamespaceVisibilityPrivate), "Namespace visibility")
	updateCmd.Flags().BoolVar(&updateRetentionEnabled, "retention-enabled", false, "Enable retention policy")
	updateCmd.Flags().BoolVar(&updateDeleteUntagged, "delete-untagged", false, "Delete untagged images during retention runs")
	updateCmd.Flags().IntVar(&updateRetentionDays, "retention-days", 0, "Retain tags for this many days")
	updateCmd.Flags().IntVar(&updateRetentionCount, "retention-count", 0, "Retain this many recent tags")
	updateCmd.Flags().StringVar(&updateRetentionPolicyFile, "retention-policy-file", "", "Path to JSON file with full retention policy")

	updateCmd.MarkFlagRequired(NamespaceFlag)
	updateCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
}
