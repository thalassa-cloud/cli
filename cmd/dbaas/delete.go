package dbaas

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/filters"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	deleteWait          bool
	deleteForce         bool
	deleteLabelSelector string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Delete database cluster(s)",
	Long:              "Delete database cluster(s) by identity or label selector.",
	Aliases:           []string{"d", "rm", "del", "remove"},
	Args:              cobra.MinimumNArgs(0),
	ValidArgsFunction: completion.CompleteDbClusterID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && deleteLabelSelector == "" {
			return fmt.Errorf("either cluster identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect clusters to delete
		clustersToDelete := []dbaas.DbCluster{}

		// If label selector is provided, filter by labels
		if deleteLabelSelector != "" {
			allClusters, err := client.DBaaS().ListDbClusters(cmd.Context(), &dbaas.ListDbClustersRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(deleteLabelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list clusters: %w", err)
			}
			if len(allClusters) == 0 {
				fmt.Println("No database clusters found matching the label selector")
				return nil
			}
			clustersToDelete = append(clustersToDelete, allClusters...)
		} else {
			// Get clusters by identity
			for _, clusterIdentity := range args {
				cluster, err := client.DBaaS().GetDbCluster(cmd.Context(), clusterIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("Database cluster %s not found\n", clusterIdentity)
						continue
					}
					return fmt.Errorf("failed to get cluster: %w", err)
				}
				clustersToDelete = append(clustersToDelete, *cluster)
			}
		}

		if len(clustersToDelete) == 0 {
			fmt.Println("No database clusters to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !deleteForce {
			fmt.Printf("Are you sure you want to delete the following database cluster(s)?\n")
			for _, cluster := range clustersToDelete {
				fmt.Printf("  %s (%s)\n", cluster.Name, cluster.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each cluster
		for _, cluster := range clustersToDelete {
			fmt.Printf("Deleting database cluster: %s (%s)\n", cluster.Name, cluster.Identity)
			err := client.DBaaS().DeleteDbCluster(cmd.Context(), cluster.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete database cluster: %w", err)
			}

			if deleteWait {
				// Poll until cluster is deleted
				// Note: WaitUntilDbClusterIsDeleted may not be available, so we poll manually
				for {
					_, err := client.DBaaS().GetDbCluster(cmd.Context(), cluster.Identity)
					if err != nil {
						if tcclient.IsNotFound(err) {
							break
						}
						return fmt.Errorf("failed to check cluster deletion status: %w", err)
					}
					// Simple polling with sleep
					select {
					case <-cmd.Context().Done():
						return cmd.Context().Err()
					case <-time.After(5 * time.Second):
						// Continue polling
					}
				}
			}
			fmt.Printf("Database cluster %s deleted successfully\n", cluster.Identity)
		}

		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&deleteWait, "wait", false, "Wait for the database cluster(s) to be deleted")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVarP(&deleteLabelSelector, "selector", "l", "", "Label selector to filter clusters (format: key1=value1,key2=value2)")

	DbaasCmd.AddCommand(deleteCmd)
}
