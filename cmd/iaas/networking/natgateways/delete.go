package natgateways

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
	wait          bool
	force         bool
	labelSelector string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete NAT gateway(s)",
	Long:    "Delete NAT gateway(s) by identity or label selector. This command will delete the NAT gateway(s) and all associated resources.",
	Example: "tcloud networking natgateways delete ngw-123\ntcloud networking natgateways delete ngw-123 ngw-456 --wait\ntcloud networking natgateways delete --selector environment=test --force",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && labelSelector == "" {
			return fmt.Errorf("either NAT gateway identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect NAT gateways to delete
		natGatewaysToDelete := []iaas.VpcNatGateway{}

		// If label selector is provided, filter by labels
		if labelSelector != "" {
			allNatGateways, err := client.IaaS().ListNatGateways(cmd.Context(), &iaas.ListNatGatewaysRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(labelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list NAT gateways: %w", err)
			}
			if len(allNatGateways) == 0 {
				fmt.Println("No NAT gateways found matching the label selector")
				return nil
			}
			natGatewaysToDelete = append(natGatewaysToDelete, allNatGateways...)
		} else {
			// Get NAT gateways by identity
			for _, ngwIdentity := range args {
				ngw, err := client.IaaS().GetNatGateway(cmd.Context(), ngwIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("NAT gateway %s not found\n", ngwIdentity)
						continue
					}
					return fmt.Errorf("failed to get NAT gateway: %w", err)
				}
				natGatewaysToDelete = append(natGatewaysToDelete, *ngw)
			}
		}

		if len(natGatewaysToDelete) == 0 {
			fmt.Println("No NAT gateways to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !force {
			fmt.Printf("Are you sure you want to delete the following NAT gateway(s)?\n")
			for _, ngw := range natGatewaysToDelete {
				fmt.Printf("  %s (%s)\n", ngw.Name, ngw.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each NAT gateway
		for _, ngw := range natGatewaysToDelete {
			fmt.Printf("Deleting NAT gateway: %s (%s)\n", ngw.Name, ngw.Identity)
			err := client.IaaS().DeleteNatGateway(cmd.Context(), ngw.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete NAT gateway: %w", err)
			}

			if wait {
				if err := client.IaaS().WaitUntilNatGatewayDeleted(cmd.Context(), ngw.Identity); err != nil {
					return fmt.Errorf("failed to wait for NAT gateway to be deleted: %w", err)
				}
			}
			fmt.Printf("NAT gateway %s deleted successfully\n", ngw.Identity)
		}

		return nil
	},
}

func init() {
	NatGatewaysCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&wait, "wait", false, "Wait for the NAT gateway(s) to be deleted")
	deleteCmd.Flags().BoolVar(&force, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVarP(&labelSelector, "selector", "l", "", "Label selector to filter NAT gateways (format: key1=value1,key2=value2)")

	// Add completion
	deleteCmd.ValidArgsFunction = completeNatGatewayID
}
