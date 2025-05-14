package subnets

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	CreateFlagName        = "name"
	CreateFlagDescription = "description"
	CreateFlagVpc         = "vpc"
	CreateFlagCIDR        = "cidr"

	CreateFlagLabels      = "labels"
	CreateFlagAnnotations = "annotations"
)

var (
	createSubnetValues = iaas.CreateSubnet{}
)

// getCmd represents the get command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a subnet",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		tcclient, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if createSubnetValues.Name == "" {
			return fmt.Errorf("name is required")
		}
		if createSubnetValues.VpcIdentity == "" {
			return fmt.Errorf("vpc is required")
		}
		if createSubnetValues.Cidr == "" {
			return fmt.Errorf("cidr is required")
		}

		vpc, err := tcclient.IaaS().GetVpc(cmd.Context(), createSubnetValues.VpcIdentity)
		if err != nil {
			if client.IsNotFound(err) {
				vpcs, err := tcclient.IaaS().ListVpcs(cmd.Context())
				if err != nil {
					return err
				}
				for _, v := range vpcs {
					if v.Slug == createSubnetValues.VpcIdentity {
						createSubnetValues.VpcIdentity = v.Identity
						vpc = &v
						break
					}
				}
				if vpc == nil {
					return fmt.Errorf("vpc not found")
				}
			} else {
				return err
			}
		}
		createSubnetValues.VpcIdentity = vpc.Identity

		subnet, err := tcclient.IaaS().CreateSubnet(cmd.Context(), createSubnetValues)
		if err != nil {
			return err
		}
		body := make([][]string, 0, 1)
		body = append(body, []string{
			subnet.Identity,
			subnet.Name,
			vpc.Name,
			subnet.Cidr,
			formattime.FormatTime(subnet.CreatedAt.Local(), showExactTime),
		})
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "VPC", "CIDR", "Age"}, body)
		}
		return nil
	},
}

func init() {
	SubnetsCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createSubnetValues.Name, CreateFlagName, "", "Name of the subnet")
	createCmd.Flags().StringVar(&createSubnetValues.Description, CreateFlagDescription, "", "Description of the subnet")
	createCmd.Flags().StringVar(&createSubnetValues.VpcIdentity, CreateFlagVpc, "", "VPC of the subnet")
	createCmd.Flags().StringVar(&createSubnetValues.Cidr, CreateFlagCIDR, "", "CIDR of the subnet")
}
