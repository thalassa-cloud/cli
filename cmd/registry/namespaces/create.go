package namespaces

import (
	"fmt"

	"github.com/spf13/cobra"

	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/containerregistry"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	createNamespace   string
	createRegion      string
	createDescription string
	createLabels      []string
	createAnnotations []string
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a container registry namespace",
	Example: "tcloud registry namespaces create --namespace my-app --region eu-west-1",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if createNamespace == "" {
			return fmt.Errorf("namespace is required")
		}
		if createRegion == "" {
			return fmt.Errorf("region is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list regions: %w", err)
		}
		region, err := iaasutil.FindRegionByIdentitySlugOrNameWithError(regions, createRegion)
		if err != nil {
			return err
		}

		ns, err := client.ContainerRegistry().CreateContainerRegistryNamespace(cmd.Context(), containerregistry.CreateContainerRegistryNamespaceRequest{
			Namespace:   createNamespace,
			Region:      region.Identity,
			Description: createDescription,
			Labels:      parseKeyValueSlice(createLabels),
			Annotations: parseKeyValueSlice(createAnnotations),
		})
		if err != nil {
			return err
		}

		fmt.Printf("Namespace created successfully\n")
		fmt.Printf("ID: %s\n", ns.Identity)
		fmt.Printf("Namespace: %s\n", ns.Namespace)
		return nil
	},
}

func init() {
	NamespacesCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createNamespace, "namespace", "", "Registry namespace name")
	createCmd.Flags().StringVar(&createRegion, "region", "", "Region identity, slug, or name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", []string{}, "Labels in key=value format")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", []string{}, "Annotations in key=value format")

	createCmd.MarkFlagRequired("namespace")
	createCmd.MarkFlagRequired("region")
	createCmd.RegisterFlagCompletionFunc("region", completeRegion)
}
