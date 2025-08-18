package securitygroups

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	name                  string
	description           string
	vpcIdentity           string
	allowSameGroupTraffic bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a security group",
	Long:    "Create a new security group within your organisation",
	Example: "tcloud networking security-groups create --name my-sg --vpc vpc-123\ntcloud networking security-groups create --name my-sg --vpc vpc-123 --description 'My security group' --allow-same-group",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if name == "" {
			return fmt.Errorf("name is required")
		}
		if vpcIdentity == "" {
			return fmt.Errorf("vpc is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		createRequest := iaas.CreateSecurityGroupRequest{
			Name:                  name,
			Description:           description,
			VpcIdentity:           vpcIdentity,
			AllowSameGroupTraffic: allowSameGroupTraffic,
		}

		securityGroup, err := client.IaaS().CreateSecurityGroup(cmd.Context(), createRequest)
		if err != nil {
			return err
		}

		fmt.Printf("Security group created successfully\n")
		fmt.Printf("ID: %s\n", securityGroup.Identity)
		fmt.Printf("Name: %s\n", securityGroup.Name)
		fmt.Printf("Status: %s\n", securityGroup.Status)
		return nil
	},
}

func init() {
	SecurityGroupsCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&name, "name", "", "Name of the security group")
	createCmd.Flags().StringVar(&description, "description", "", "Description of the security group")
	createCmd.Flags().StringVar(&vpcIdentity, "vpc", "", "VPC identity where the security group will be created")
	createCmd.Flags().BoolVar(&allowSameGroupTraffic, "allow-same-group", false, "Allow traffic between instances in the same security group")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("vpc")

	// Add completion
	createCmd.RegisterFlagCompletionFunc("vpc", completeVPCID)
}
