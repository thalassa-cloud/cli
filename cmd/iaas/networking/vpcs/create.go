package vpcs

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

const (
	CreateFlagName        = "name"
	CreateFlagDescription = "description"
	CreateFlagRegion      = "region"
	CreateFlagCIDRs       = "cidrs"

	CreateFlagLabels      = "labels"
	CreateFlagAnnotations = "annotations"
)

var (
	createVpcValues = iaas.CreateVpc{}
)

// getCmd represents the get command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a vpc",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			return err
		}

		if createVpcValues.Name == "" {
			return fmt.Errorf("name is required")
		}
		if createVpcValues.CloudRegionIdentity == "" {
			return fmt.Errorf("region is required")
		}

		if len(createVpcValues.VpcCidrs) == 0 {
			return fmt.Errorf("cidrs is required")
		}

		region, err := client.IaaS().GetRegion(cmd.Context(), createVpcValues.CloudRegionIdentity)
		if err != nil {
			return err
		}
		createVpcValues.CloudRegionIdentity = region.Identity

		vpc, err := client.IaaS().CreateVpc(cmd.Context(), createVpcValues)
		if err != nil {
			return err
		}
		body := make([][]string, 0, 1)
		body = append(body, []string{
			vpc.Identity,
			vpc.Name,
			vpc.CloudRegion.Name,

			strings.Join(vpc.CIDRs, ", "),
			formattime.FormatTime(vpc.CreatedAt.Local(), showExactTime),
		})
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Region", "CIDRs", "Age"}, body)
		}
		return nil
	},
}

func init() {
	VpcsCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createVpcValues.Name, CreateFlagName, "", "Name of the vpc")
	createCmd.Flags().StringVar(&createVpcValues.Description, CreateFlagDescription, "", "Description of the vpc")
	createCmd.Flags().StringVar(&createVpcValues.CloudRegionIdentity, CreateFlagRegion, "", "Region of the vpc")
	createCmd.Flags().StringSliceVar(&createVpcValues.VpcCidrs, CreateFlagCIDRs, []string{"10.0.0.0/16"}, "CIDRs of the vpc")
	// createCmd.Flags().StringSliceVar(&createVpcValues.Labels, CreateFlagLabels, []string{}, "Labels of the vpc")
	// createCmd.Flags().StringSliceVar(&createVpcValues.Annotations, CreateFlagAnnotations, []string{}, "Annotations of the vpc")
}
