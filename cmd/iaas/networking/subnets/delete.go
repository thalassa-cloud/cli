package subnets

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a subnet",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			return err
		}

		// get the subnets
		subnets, err := client.IaaS().ListSubnets(cmd.Context())
		if err != nil {
			return err
		}

		for _, subnetIdentityOrSlug := range args {
			subnetIdentity := subnetIdentityOrSlug
			var deleteSubnet *iaas.Subnet
			for _, subnet := range subnets {
				if subnetIdentity == subnet.Identity || subnetIdentity == subnet.Name || subnetIdentity == subnet.Slug {
					deleteSubnet = &subnet
					break
				}
			}

			if deleteSubnet == nil {
				fmt.Printf("Subnet %s not found\n", subnetIdentityOrSlug)
				continue
			}
			err := client.IaaS().DeleteSubnet(cmd.Context(), deleteSubnet.Identity)
			if err != nil {
				return err
			}
			fmt.Printf("Deleted subnet %s (%s)\n", deleteSubnet.Name, deleteSubnet.Identity)
		}

		return nil
	},
}

func init() {
	SubnetsCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	// TODO: implement filters
}
