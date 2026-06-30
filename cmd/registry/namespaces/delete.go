package namespaces

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/containerregistry"
	"github.com/thalassa-cloud/client-go/filters"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	deleteForce         bool
	deleteLabelSelector string
)

var deleteCmd = &cobra.Command{
	Use:     "delete [NAMESPACE...]",
	Short:   "Delete container registry namespace(s)",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && deleteLabelSelector == "" {
			return fmt.Errorf("either namespace identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		toDelete := []containerregistry.ContainerRegistryNamespace{}

		if deleteLabelSelector != "" {
			all, err := client.ContainerRegistry().ListContainerRegistryNamespaces(cmd.Context(), &containerregistry.ListContainerRegistryNamespacesRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(deleteLabelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list namespaces: %w", err)
			}
			toDelete = append(toDelete, all...)
		} else {
			for _, id := range args {
				ns, err := client.ContainerRegistry().GetContainerRegistryNamespace(cmd.Context(), id)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("Namespace %s not found\n", id)
						continue
					}
					return fmt.Errorf("failed to get namespace: %w", err)
				}
				toDelete = append(toDelete, *ns)
			}
		}

		if len(toDelete) == 0 {
			fmt.Println("No namespaces to delete")
			return nil
		}

		if !deleteForce {
			fmt.Printf("Are you sure you want to delete %d namespace(s)?\n", len(toDelete))
			var confirm string
			fmt.Print("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		for _, ns := range toDelete {
			fmt.Printf("Deleting namespace: %s (%s)\n", ns.Namespace, ns.Identity)
			if err := client.ContainerRegistry().DeleteContainerRegistryNamespace(cmd.Context(), ns.Identity); err != nil {
				return fmt.Errorf("failed to delete namespace: %w", err)
			}
			fmt.Printf("Namespace %s deleted successfully\n", ns.Identity)
		}
		return nil
	},
}

func init() {
	NamespacesCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation")
	deleteCmd.Flags().StringVarP(&deleteLabelSelector, "selector", "l", "", "Label selector (format: key1=value1,key2=value2)")

	deleteCmd.ValidArgsFunction = completeNamespaceID
}
