package subnets

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	wait          bool
	force         bool
	labelSelector string
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete subnet(s)",
	Long:    "Delete subnet(s) by identity or label selector.",
	Example: "tcloud networking subnets delete subnet-123\ntcloud networking subnets delete subnet-123 subnet-456 --wait\ntcloud networking subnets delete --selector environment=test --force",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && labelSelector == "" {
			return fmt.Errorf("either subnet identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect subnets to delete
		subnetsToDelete := []iaas.Subnet{}

		// If label selector is provided, filter by labels
		if labelSelector != "" {
			allSubnets, err := client.IaaS().ListSubnets(cmd.Context(), &iaas.ListSubnetsRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(labelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list subnets: %w", err)
			}
			if len(allSubnets) == 0 {
				fmt.Println("No subnets found matching the label selector")
				return nil
			}
			subnetsToDelete = append(subnetsToDelete, allSubnets...)
		} else {
			// Get subnets by identity (support identity, name, or slug)
			allSubnets, err := client.IaaS().ListSubnets(cmd.Context(), &iaas.ListSubnetsRequest{})
			if err != nil {
				return fmt.Errorf("failed to list subnets: %w", err)
			}

			for _, subnetIdentityOrSlug := range args {
				var deleteSubnet *iaas.Subnet
				for i := range allSubnets {
					subnet := allSubnets[i]
					if subnetIdentityOrSlug == subnet.Identity || subnetIdentityOrSlug == subnet.Name || subnetIdentityOrSlug == subnet.Slug {
						deleteSubnet = &subnet
						break
					}
				}

				if deleteSubnet == nil {
					fmt.Printf("Subnet %s not found\n", subnetIdentityOrSlug)
					continue
				}
				subnetsToDelete = append(subnetsToDelete, *deleteSubnet)
			}
		}

		if len(subnetsToDelete) == 0 {
			fmt.Println("No subnets to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !force {
			fmt.Printf("Are you sure you want to delete the following subnet(s)?\n")
			for _, subnet := range subnetsToDelete {
				fmt.Printf("  %s (%s)\n", subnet.Name, subnet.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each subnet
		for _, subnet := range subnetsToDelete {
			fmt.Printf("Deleting subnet: %s (%s)\n", subnet.Name, subnet.Identity)
			err := client.IaaS().DeleteSubnet(cmd.Context(), subnet.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete subnet: %w", err)
			}

			if wait {
				if err := client.IaaS().WaitUntilSubnetDeleted(cmd.Context(), subnet.Identity); err != nil {
					return fmt.Errorf("failed to wait for subnet to be deleted: %w", err)
				}
			}
			fmt.Printf("Subnet %s deleted successfully\n", subnet.Identity)
		}

		return nil
	},
}

func init() {
	SubnetsCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&wait, "wait", false, "Wait for the subnet(s) to be deleted")
	deleteCmd.Flags().BoolVar(&force, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVarP(&labelSelector, "selector", "l", "", "Label selector to filter subnets (format: key1=value1,key2=value2)")
}
