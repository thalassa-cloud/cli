package vpcs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"

	"github.com/thalassa-cloud/client-go/iaas"
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
	createVpcWait  bool
)

// getCmd represents the get command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a vpc",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
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

		regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
		if err != nil {
			return err
		}
		region, err := iaasutil.FindRegionByIdentitySlugOrNameWithError(regions, createVpcValues.CloudRegionIdentity)
		if err != nil {
			return err
		}
		createVpcValues.CloudRegionIdentity = region.Identity

		vpc, err := client.IaaS().CreateVpc(cmd.Context(), createVpcValues)
		if err != nil {
			return err
		}

		if createVpcWait {
			ctxWithTimeout, cancel := context.WithTimeout(cmd.Context(), 10*time.Minute)
			defer cancel()

			fmt.Println("Waiting for VPC to be ready...")
			for {
				vpc, err = client.IaaS().GetVpc(ctxWithTimeout, vpc.Identity)
				if err != nil {
					return fmt.Errorf("failed to get vpc: %w", err)
				}
				// VPC is ready when status is "ready" or "available"
				if strings.EqualFold(vpc.Status, "ready") || strings.EqualFold(vpc.Status, "available") {
					break
				}
				// Check for failed state
				if strings.EqualFold(vpc.Status, "failed") || strings.EqualFold(vpc.Status, "error") {
					return fmt.Errorf("vpc creation failed with status: %s", vpc.Status)
				}
				select {
				case <-ctxWithTimeout.Done():
					return fmt.Errorf("timeout waiting for vpc %s to be ready (current status: %s)", vpc.Identity, vpc.Status)
				case <-time.After(2 * time.Second):
					// Continue polling
				}
			}
			fmt.Println("VPC is ready")
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
	createCmd.Flags().BoolVar(&createVpcWait, "wait", false, "Wait for the VPC to be ready before returning")
	// createCmd.Flags().StringSliceVar(&createVpcValues.Labels, CreateFlagLabels, []string{}, "Labels of the vpc")
	// createCmd.Flags().StringSliceVar(&createVpcValues.Annotations, CreateFlagAnnotations, []string{}, "Annotations of the vpc")
}
