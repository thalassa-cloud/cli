package configuration

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete namespace configuration",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			return fmt.Errorf("--namespace is required")
		}

		if !deleteForce {
			fmt.Printf("Delete configuration for namespace %s?\n", namespace)
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

		if err := client.ContainerRegistry().DeleteNamespaceConfiguration(cmd.Context(), namespace); err != nil {
			return err
		}

		fmt.Printf("Configuration deleted for namespace %s\n", namespace)
		return nil
	},
}

func init() {
	ConfigurationCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation")
	deleteCmd.MarkFlagRequired(NamespaceFlag)
	deleteCmd.RegisterFlagCompletionFunc(NamespaceFlag, completeNamespaceID)
}
