package subnets

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/pkg/iaas"
	"github.com/thalassa-cloud/client-go/pkg/thalassa"
)

const (
	CreateFlagName        = "name"
	CreateFlagDescription = "description"
	CreateFlagVpc         = "vpc"
	CreateFlagCIDR        = "cidr"
	CreateFlagZone        = "zone"

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
		tcclient, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			return err
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
		if createSubnetValues.CloudZone == "" {
			return fmt.Errorf("zone is required")
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

		// validate the zone
		if len(vpc.CloudRegion.Zones) == 0 {
			return fmt.Errorf("no zones found in the region")
		}
		found := false
		for _, z := range vpc.CloudRegion.Zones {
			if z.Slug == createSubnetValues.CloudZone {
				found = true
				createSubnetValues.CloudZone = z.Slug
				break
			} else if z.Identity == createSubnetValues.CloudZone {
				found = true
				createSubnetValues.CloudZone = z.Slug
				break
			}
		}
		if !found {
			return fmt.Errorf("zone not found")
		}
		fmt.Printf("zone: %s\n", createSubnetValues.CloudZone)

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
	createCmd.Flags().StringVar(&createSubnetValues.CloudZone, CreateFlagZone, "", "Zone of the subnet")
}
