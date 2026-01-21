package tfs

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/tfs"
)

var (
	deleteWait          bool
	deleteForce         bool
	deleteLabelSelector string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Delete TFS instance(s)",
	Long:              "Delete TFS instance(s) by identity or label selector.",
	Aliases:           []string{"d", "rm", "del", "remove"},
	Args:              cobra.MinimumNArgs(0),
	ValidArgsFunction: completion.CompleteTfsInstanceID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && deleteLabelSelector == "" {
			return fmt.Errorf("either TFS instance identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect instances to delete
		instancesToDelete := []tfs.TfsInstance{}

		// If label selector is provided, filter by labels
		if deleteLabelSelector != "" {
			allInstances, err := client.Tfs().ListTfsInstances(cmd.Context(), &tfs.ListTfsInstancesRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(deleteLabelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list TFS instances: %w", err)
			}
			if len(allInstances) == 0 {
				fmt.Println("No TFS instances found matching the label selector")
				return nil
			}
			instancesToDelete = append(instancesToDelete, allInstances...)
		} else {
			// Get instances by identity
			for _, instanceIdentity := range args {
				instance, err := client.Tfs().GetTfsInstance(cmd.Context(), instanceIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("TFS instance %s not found\n", instanceIdentity)
						continue
					}
					return fmt.Errorf("failed to get TFS instance: %w", err)
				}
				instancesToDelete = append(instancesToDelete, *instance)
			}
		}

		if len(instancesToDelete) == 0 {
			fmt.Println("No TFS instances to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !deleteForce {
			fmt.Printf("Are you sure you want to delete the following TFS instance(s)?\n")
			for _, instance := range instancesToDelete {
				fmt.Printf("  %s (%s)\n", instance.Name, instance.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each instance
		for _, instance := range instancesToDelete {
			fmt.Printf("Deleting TFS instance: %s (%s)\n", instance.Name, instance.Identity)
			err := client.Tfs().DeleteTfsInstance(cmd.Context(), instance.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete TFS instance: %w", err)
			}

			if deleteWait {
				if err := client.Tfs().WaitUntilTfsInstanceIsDeleted(cmd.Context(), instance.Identity); err != nil {
					return fmt.Errorf("failed to wait for TFS instance to be deleted: %w", err)
				}
			}
			fmt.Printf("TFS instance %s deleted successfully\n", instance.Identity)
		}

		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&deleteWait, "wait", false, "Wait for the TFS instance(s) to be deleted")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVarP(&deleteLabelSelector, "selector", "l", "", "Label selector to filter TFS instances (format: key1=value1,key2=value2)")

	TfsCmd.AddCommand(deleteCmd)
}
