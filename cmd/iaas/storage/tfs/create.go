package tfs

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/tfs"
)

var (
	createTfsName             string
	createTfsDescription      string
	createTfsRegion           string
	createTfsLabels           []string
	createTfsAnnotations      []string
	createTfsDeleteProtection bool
	createTfsWait             bool
	createTfsSizeGb           int

	createTfsVpc    string
	createTfsSubnet string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a TFS instance",
	Long:  "Create a new TFS (Thalassa File System) instance for shared file storage.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if createTfsName == "" {
			return fmt.Errorf("name is required")
		}
		if createTfsRegion == "" {
			return fmt.Errorf("region is required")
		}

		// Resolve region
		regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
		if err != nil {
			return fmt.Errorf("failed to get region: %w", err)
		}
		var region *iaas.Region
		for _, r := range regions {
			if r.Identity == createTfsRegion || r.Slug == createTfsRegion || r.Name == createTfsRegion {
				region = &r
				break
			}
		}
		if region == nil {
			return fmt.Errorf("region not found: %s", createTfsRegion)
		}

		// Resolve vpc
		vpc, err := client.IaaS().GetVpc(cmd.Context(), createTfsVpc)
		if err != nil {
			return fmt.Errorf("failed to get vpc: %w", err)
		}
		// vpc must be in the same region
		if vpc.CloudRegion.Identity != region.Identity {
			return fmt.Errorf("vpc must be in the same region")
		}

		// Resolve subnet
		subnet, err := client.IaaS().GetSubnet(cmd.Context(), createTfsSubnet)
		if err != nil {
			return fmt.Errorf("failed to get subnet: %w", err)
		}
		// subnet must be in the same vpc
		if subnet.Vpc.Identity != vpc.Identity {
			return fmt.Errorf("subnet must be in the same vpc")
		}

		// Parse labels from key=value format
		labels := make(map[string]string)
		for _, label := range createTfsLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Parse annotations from key=value format
		annotations := make(map[string]string)
		for _, annotation := range createTfsAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		req := tfs.CreateTfsInstanceRequest{
			Name:                createTfsName,
			Description:         createTfsDescription,
			CloudRegionIdentity: region.Identity,
			Labels:              labels,
			Annotations:         annotations,
			DeleteProtection:    createTfsDeleteProtection,

			SizeGB: createTfsSizeGb,
			// vpc
			VpcIdentity:    vpc.Identity,
			SubnetIdentity: subnet.Identity,
		}

		instance, err := client.Tfs().CreateTfsInstance(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create TFS instance: %w", err)
		}

		if createTfsWait {
			// Wait for instance to be available
			err = client.Tfs().WaitUntilTfsInstanceIsAvailable(cmd.Context(), instance.Identity)
			if err != nil {
				return fmt.Errorf("failed to wait for TFS instance to be available: %w", err)
			}
			// Refresh instance to get latest status
			instance, err = client.Tfs().GetTfsInstance(cmd.Context(), instance.Identity)
			if err != nil {
				return fmt.Errorf("failed to get TFS instance: %w", err)
			}
		}

		// Output in table format
		regionName := ""
		if instance.Region != nil {
			regionName = instance.Region.Name
		}
		body := [][]string{
			{
				instance.Identity,
				instance.Name,
				string(instance.Status),
				regionName,
				formattime.FormatTime(instance.CreatedAt.Local(), false),
			},
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Status", "Region", "Age"}, body)
		}

		return nil
	},
}

func init() {
	TfsCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createTfsName, "name", "", "Name of the TFS instance (required)")
	createCmd.Flags().StringVar(&createTfsDescription, "description", "", "Description of the TFS instance")
	createCmd.Flags().StringVar(&createTfsRegion, "region", "", "Region of the TFS instance (required)")
	createCmd.Flags().StringVar(&createTfsVpc, "vpc", "", "VPC of the TFS instance (required)")
	createCmd.Flags().StringVar(&createTfsSubnet, "subnet", "", "Subnet of the TFS instance (required)")
	createCmd.Flags().IntVar(&createTfsSizeGb, "size", 1, "Size of the TFS instance in GB (required)")
	createCmd.Flags().StringSliceVar(&createTfsLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	createCmd.Flags().StringSliceVar(&createTfsAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	createCmd.Flags().BoolVar(&createTfsDeleteProtection, "delete-protection", false, "Enable delete protection")
	createCmd.Flags().BoolVar(&createTfsWait, "wait", false, "Wait for the TFS instance to be available before returning")

	// Register completions
	createCmd.RegisterFlagCompletionFunc("vpc", completion.CompleteVPCID)
	createCmd.RegisterFlagCompletionFunc("subnet", completion.CompleteSubnetEnhanced)

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("region")
}
