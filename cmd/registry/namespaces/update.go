package namespaces

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/containerregistry"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateDescription string
	updateLabels        []string
	updateAnnotations   []string
)

var updateCmd = &cobra.Command{
	Use:     "update NAMESPACE",
	Short:   "Update a container registry namespace",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.ContainerRegistry().GetContainerRegistryNamespace(cmd.Context(), args[0])
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("namespace not found: %s", args[0])
			}
			return fmt.Errorf("failed to get namespace: %w", err)
		}

		req := containerregistry.UpdateContainerRegistryNamespaceRequest{
			Description: current.Description,
			Labels:      current.Labels,
			Annotations: current.Annotations,
		}

		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = parseKeyValueSlice(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = parseKeyValueSlice(updateAnnotations)
		}

		ns, err := client.ContainerRegistry().UpdateContainerRegistryNamespace(cmd.Context(), current.Identity, req)
		if err != nil {
			return err
		}

		fmt.Printf("Namespace updated successfully\n")
		fmt.Printf("ID: %s\n", ns.Identity)
		fmt.Printf("Namespace: %s\n", ns.Namespace)
		return nil
	},
}

func init() {
	NamespacesCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", []string{}, "Labels in key=value format")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", []string{}, "Annotations in key=value format")

	updateCmd.ValidArgsFunction = completeNamespaceID
}
