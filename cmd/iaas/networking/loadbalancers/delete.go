package loadbalancers

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	deleteWait          bool
	deleteForce         bool
	deleteLabelSelector string
)

var deleteCmd = &cobra.Command{
	Use:     "delete [LOADBALANCER...]",
	Short:   "Delete load balancer(s)",
	Long:    "Delete load balancer(s) by identity or label selector.",
	Example: "tcloud networking loadbalancers delete lb-123\ntcloud networking loadbalancers delete lb-123 --wait\ntcloud networking loadbalancers delete --selector env=test --force",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && deleteLabelSelector == "" {
			return fmt.Errorf("either load balancer identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		loadbalancersToDelete := []iaas.VpcLoadbalancer{}

		if deleteLabelSelector != "" {
			all, err := client.IaaS().ListLoadbalancers(cmd.Context(), &iaas.ListLoadbalancersRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(deleteLabelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list load balancers: %w", err)
			}
			if len(all) == 0 {
				fmt.Println("No load balancers found matching the label selector")
				return nil
			}
			loadbalancersToDelete = append(loadbalancersToDelete, all...)
		} else {
			for _, lbIdentity := range args {
				lb, err := client.IaaS().GetLoadbalancer(cmd.Context(), lbIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("Load balancer %s not found\n", lbIdentity)
						continue
					}
					return fmt.Errorf("failed to get load balancer: %w", err)
				}
				loadbalancersToDelete = append(loadbalancersToDelete, *lb)
			}
		}

		if len(loadbalancersToDelete) == 0 {
			fmt.Println("No load balancers to delete")
			return nil
		}

		if !deleteForce {
			fmt.Printf("Are you sure you want to delete the following load balancer(s)?\n")
			for _, lb := range loadbalancersToDelete {
				fmt.Printf("  %s (%s)\n", lb.Name, lb.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		for _, lb := range loadbalancersToDelete {
			fmt.Printf("Deleting load balancer: %s (%s)\n", lb.Name, lb.Identity)
			if err := client.IaaS().DeleteLoadbalancer(cmd.Context(), lb.Identity); err != nil {
				return fmt.Errorf("failed to delete load balancer: %w", err)
			}
			if deleteWait {
				if err := client.IaaS().WaitUntilLoadbalancerIsDeleted(cmd.Context(), lb.Identity); err != nil {
					return fmt.Errorf("failed to wait for load balancer deletion: %w", err)
				}
			}
			fmt.Printf("Load balancer %s deleted successfully\n", lb.Identity)
		}

		return nil
	},
}

func init() {
	LoadbalancersCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteWait, "wait", false, "Wait for the load balancer(s) to be deleted")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Force deletion and skip confirmation")
	deleteCmd.Flags().StringVarP(&deleteLabelSelector, "selector", "l", "", "Label selector (format: key1=value1,key2=value2)")

	deleteCmd.ValidArgsFunction = completeLoadbalancerID
}
