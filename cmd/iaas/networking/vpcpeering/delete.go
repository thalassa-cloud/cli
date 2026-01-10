package vpcpeering

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
	deleteForce         bool
	deleteLabelSelector string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete VPC peering connection(s)",
	Long:    "Delete VPC peering connection(s) by identity or label selector",
	Example: "tcloud networking vpc-peering delete vpcpc-123\ntcloud networking vpc-peering delete vpcpc-123 vpcpc-456 --force\ntcloud networking vpc-peering delete --selector environment=test --force",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && deleteLabelSelector == "" {
			return fmt.Errorf("either connection identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect connections to delete
		connectionsToDelete := []iaas.VpcPeeringConnection{}

		// If label selector is provided, filter by labels
		if deleteLabelSelector != "" {
			allConnections, err := client.IaaS().ListVpcPeeringConnections(cmd.Context(), &iaas.ListVpcPeeringConnectionsRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(deleteLabelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list VPC peering connections: %w", err)
			}
			if len(allConnections) == 0 {
				fmt.Println("No VPC peering connections found matching the label selector")
				return nil
			}
			connectionsToDelete = append(connectionsToDelete, allConnections...)
		} else {
			// Get connections by identity
			for _, connectionIdentity := range args {
				connection, err := client.IaaS().GetVpcPeeringConnection(cmd.Context(), connectionIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("VPC peering connection %s not found\n", connectionIdentity)
						continue
					}
					return fmt.Errorf("failed to get VPC peering connection: %w", err)
				}
				connectionsToDelete = append(connectionsToDelete, *connection)
			}
		}

		if len(connectionsToDelete) == 0 {
			fmt.Println("No VPC peering connections to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !deleteForce {
			fmt.Printf("Are you sure you want to delete the following VPC peering connection(s)?\n")
			for _, conn := range connectionsToDelete {
				fmt.Printf("  %s (%s)\n", conn.Name, conn.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each connection
		for _, conn := range connectionsToDelete {
			fmt.Printf("Deleting VPC peering connection: %s (%s)\n", conn.Name, conn.Identity)
			err := client.IaaS().DeleteVpcPeeringConnection(cmd.Context(), conn.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete VPC peering connection: %w", err)
			}
			fmt.Printf("VPC peering connection %s deleted successfully\n", conn.Identity)
		}

		return nil
	},
}

func init() {
	VpcPeeringCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVarP(&deleteLabelSelector, "selector", "l", "", "Label selector to filter connections (format: key1=value1,key2=value2)")
}
