package vpcs

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
	Short:   "Delete VPC(s)",
	Long:    "Delete VPC(s) by identity or label selector. This command will delete the VPC(s) and all associated resources.",
	Example: "tcloud networking vpcs delete vpc-123\ntcloud networking vpcs delete vpc-123 vpc-456 --wait\ntcloud networking vpcs delete --selector environment=test --force",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && labelSelector == "" {
			return fmt.Errorf("either VPC identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect VPCs to delete
		vpcsToDelete := []iaas.Vpc{}

		// If label selector is provided, filter by labels
		if labelSelector != "" {
			allVpcs, err := client.IaaS().ListVpcs(cmd.Context(), &iaas.ListVpcsRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(labelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list VPCs: %w", err)
			}
			if len(allVpcs) == 0 {
				fmt.Println("No VPCs found matching the label selector")
				return nil
			}
			vpcsToDelete = append(vpcsToDelete, allVpcs...)
		} else {
			// Get VPCs by identity
			for _, vpcIdentity := range args {
				vpc, err := client.IaaS().GetVpc(cmd.Context(), vpcIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("VPC %s not found\n", vpcIdentity)
						continue
					}
					return fmt.Errorf("failed to get VPC: %w", err)
				}
				vpcsToDelete = append(vpcsToDelete, *vpc)
			}
		}

		if len(vpcsToDelete) == 0 {
			fmt.Println("No VPCs to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !force {
			fmt.Printf("Are you sure you want to delete the following VPC(s)?\n")
			for _, vpc := range vpcsToDelete {
				fmt.Printf("  %s (%s)\n", vpc.Name, vpc.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each VPC
		for _, vpc := range vpcsToDelete {
			fmt.Printf("Deleting VPC: %s (%s)\n", vpc.Name, vpc.Identity)
			err := client.IaaS().DeleteVpc(cmd.Context(), vpc.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete VPC: %w", err)
			}

			if wait {
				if err := client.IaaS().WaitUntilVpcIsDeleted(cmd.Context(), vpc.Identity); err != nil {
					return fmt.Errorf("failed to wait for VPC to be deleted: %w", err)
				}
			}
			fmt.Printf("VPC %s deleted successfully\n", vpc.Identity)
		}

		return nil
	},
}

func init() {
	VpcsCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&wait, "wait", false, "Wait for the VPC(s) to be deleted")
	deleteCmd.Flags().BoolVar(&force, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVarP(&labelSelector, "selector", "l", "", "Label selector to filter VPCs (format: key1=value1,key2=value2)")
}
