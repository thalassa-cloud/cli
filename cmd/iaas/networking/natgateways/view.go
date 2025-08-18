package natgateways

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	"gopkg.in/yaml.v3"
)

var (
	outputFormat string
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View NAT gateway details",
	Long:    "View detailed information about a specific NAT gateway",
	Example: "tcloud networking natgateways view ngw-123\ntcloud networking natgateways view ngw-123 --output yaml",
	Aliases: []string{"show", "get", "describe"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		natGatewayIdentity := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		natGateway, err := client.IaaS().GetNatGateway(cmd.Context(), natGatewayIdentity)
		if err != nil {
			return fmt.Errorf("failed to get NAT gateway: %w", err)
		}

		if outputFormat == "yaml" {
			return outputYAML(*natGateway)
		}

		// Print basic information
		fmt.Printf("NAT Gateway Details:\n")
		fmt.Printf("  ID: %s\n", natGateway.Identity)
		fmt.Printf("  Name: %s\n", natGateway.Name)
		fmt.Printf("  Description: %s\n", natGateway.Description)
		fmt.Printf("  Status: %s\n", natGateway.Status)
		fmt.Printf("  Endpoint IP: %s\n", natGateway.EndpointIP)
		fmt.Printf("  IPv4 IP: %s\n", natGateway.V4IP)
		fmt.Printf("  Created: %s\n", formattime.FormatTime(natGateway.CreatedAt.Local(), false))
		fmt.Printf("  Updated: %s\n", formattime.FormatTime(natGateway.UpdatedAt.Local(), false))

		if natGateway.Vpc != nil {
			fmt.Printf("  VPC: %s (%s)\n", natGateway.Vpc.Name, natGateway.Vpc.Identity)
		}

		if natGateway.Subnet != nil {
			fmt.Printf("  Subnet: %s (%s)\n", natGateway.Subnet.Name, natGateway.Subnet.Identity)
		}

		return nil
	},
}

// outputYAML formats the NAT gateway as YAML
func outputYAML(natGateway iaas.VpcNatGateway) error {
	// clean up the NAT gateway
	natGateway.Organisation = nil

	// Marshal the original struct directly to YAML
	yamlData, err := yaml.Marshal(&natGateway)
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	fmt.Print(string(yamlData))
	return nil
}

func init() {
	NatGatewaysCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml)")
	
	// Add completion
	viewCmd.RegisterFlagCompletionFunc("output", completeOutputFormat)
	viewCmd.ValidArgsFunction = completeNatGatewayID
}


