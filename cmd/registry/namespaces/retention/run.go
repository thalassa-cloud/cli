package retention

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var namespace string

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run the retention policy for a namespace",
	Example: "tcloud registry namespaces retention run --namespace crns-123",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			return fmt.Errorf("--namespace is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if err := client.ContainerRegistry().RunRetentionPolicy(cmd.Context(), namespace); err != nil {
			return err
		}

		fmt.Printf("Retention policy run started for namespace %s\n", namespace)
		return nil
	},
}

func init() {
	RetentionCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&namespace, NamespaceFlag, "", "Namespace identity")
	runCmd.MarkFlagRequired(NamespaceFlag)
	runCmd.RegisterFlagCompletionFunc(NamespaceFlag, completion.CompleteContainerRegistryNamespaceID)
}
