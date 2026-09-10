package reservedips

import (
	"fmt"

	"github.com/spf13/cobra"

	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	createName        string
	createDescription string
	createRegion      string
	createLabels      []string
	createAnnotations []string
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a reserved IP address",
	Long:    "Create a new reserved public IP address in the specified region.",
	Example: "tcloud networking reserved-ips create --name my-ip --region nl-ams\ntcloud networking reserved-ips create --name lb-ip --region nl-ams --description 'for production LB'",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if createName == "" {
			return fmt.Errorf("name is required")
		}
		if createRegion == "" {
			return fmt.Errorf("region is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list regions: %w", err)
		}
		region, err := iaasutil.FindRegionByIdentitySlugOrNameWithError(regions, createRegion)
		if err != nil {
			return err
		}

		req := iaas.CreateReservedIpRequest{
			Name:        createName,
			Description: createDescription,
			Labels:      parseKeyValueSlice(createLabels),
			Annotations: parseKeyValueSlice(createAnnotations),
			Region:      region.Identity,
		}

		rip, err := client.IaaS().CreateReservedIP(cmd.Context(), req)
		if err != nil {
			return err
		}

		body := [][]string{{
			rip.Identity,
			rip.Name,
			regionName(*rip),
			string(rip.Status),
			dashIfEmpty(rip.IPv4Address),
			dashIfEmpty(rip.IPv6Address),
			attachedTo(*rip),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Region", "Status", "IPv4", "IPv6", "AttachedTo"}, body)
		}
		return nil
	},
}

func init() {
	ReservedIPsCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createName, "name", "", "Name of the reserved IP")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description of the reserved IP")
	createCmd.Flags().StringVar(&createRegion, "region", "", "Region for the reserved IP")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", []string{}, "Labels in key=value format")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", []string{}, "Annotations in key=value format")

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("region")
	_ = createCmd.RegisterFlagCompletionFunc("region", completeRegion)
}
