package securitygroups

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
	Short:   "Delete security group(s)",
	Long:    "Delete security group(s) by identity or label selector. This command will delete the security group(s) and all their rules.",
	Example: "tcloud networking security-groups delete sg-123\ntcloud networking security-groups delete sg-123 sg-456 --wait\ntcloud networking security-groups delete --selector environment=test --force",
	Aliases: []string{"d", "del", "remove"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && labelSelector == "" {
			return fmt.Errorf("either security group identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect security groups to delete
		securityGroupsToDelete := []iaas.SecurityGroup{}

		// If label selector is provided, filter by labels
		if labelSelector != "" {
			allSecurityGroups, err := client.IaaS().ListSecurityGroups(cmd.Context(), &iaas.ListSecurityGroupsRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(labelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list security groups: %w", err)
			}
			if len(allSecurityGroups) == 0 {
				fmt.Println("No security groups found matching the label selector")
				return nil
			}
			securityGroupsToDelete = append(securityGroupsToDelete, allSecurityGroups...)
		} else {
			// Get security groups by identity
			for _, sgIdentity := range args {
				sg, err := client.IaaS().GetSecurityGroup(cmd.Context(), sgIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("Security group %s not found\n", sgIdentity)
						continue
					}
					return fmt.Errorf("failed to get security group: %w", err)
				}
				securityGroupsToDelete = append(securityGroupsToDelete, *sg)
			}
		}

		if len(securityGroupsToDelete) == 0 {
			fmt.Println("No security groups to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !force {
			fmt.Printf("Are you sure you want to delete the following security group(s)?\n")
			for _, sg := range securityGroupsToDelete {
				fmt.Printf("  %s (%s)\n", sg.Name, sg.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each security group
		for _, sg := range securityGroupsToDelete {
			fmt.Printf("Deleting security group: %s (%s)\n", sg.Name, sg.Identity)
			err := client.IaaS().DeleteSecurityGroup(cmd.Context(), sg.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete security group: %w", err)
			}
			fmt.Printf("Security group %s deleted successfully\n", sg.Identity)
		}

		return nil
	},
}

func init() {
	SecurityGroupsCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&wait, "wait", false, "Wait for the security group(s) to be deleted")
	deleteCmd.Flags().BoolVar(&force, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVar(&labelSelector, "selector", "", "Label selector to filter security groups (format: key1=value1,key2=value2)")

	// Add completion
	deleteCmd.ValidArgsFunction = completeSecurityGroupID
}
