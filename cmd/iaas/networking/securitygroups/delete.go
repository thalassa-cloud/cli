package securitygroups

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	wait bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a security group",
	Long:    "Delete a security group. This command will delete the security group and all its rules.",
	Example: "tcloud networking security-groups delete sg-123\ntcloud networking security-groups delete sg-123 --wait",
	Aliases: []string{"d", "del", "remove"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		securityGroupIdentity := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// First, get the security group to show what we're deleting
		securityGroup, err := client.IaaS().GetSecurityGroup(cmd.Context(), securityGroupIdentity)
		if err != nil {
			return fmt.Errorf("failed to get security group: %w", err)
		}

		fmt.Printf("Deleting security group: %s (%s)\n", securityGroup.Name, securityGroup.Identity)

		err = client.IaaS().DeleteSecurityGroup(cmd.Context(), securityGroupIdentity)
		if err != nil {
			return fmt.Errorf("failed to delete security group: %w", err)
		}

		fmt.Println("Security group deleted successfully")
		return nil
	},
}

func init() {
	SecurityGroupsCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&wait, "wait", false, "Wait for the security group to be deleted")

	// Add completion
	deleteCmd.ValidArgsFunction = completeSecurityGroupID
}
