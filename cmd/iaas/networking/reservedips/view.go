package reservedips

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var viewOutputFormat string

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View reserved IP details",
	Long:    "View detailed information about a specific reserved IP address.",
	Example: "tcloud networking reserved-ips view rip-123\ntcloud networking reserved-ips view rip-123 --output yaml",
	Aliases: []string{"show", "get", "describe"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		rip, err := client.IaaS().GetReservedIP(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get reserved IP: %w", err)
		}

		if viewOutputFormat == "yaml" {
			return outputYAML(*rip)
		}

		fmt.Printf("Reserved IP Details:\n")
		fmt.Printf("  ID: %s\n", rip.Identity)
		fmt.Printf("  Name: %s\n", rip.Name)
		fmt.Printf("  Description: %s\n", rip.Description)
		fmt.Printf("  Status: %s\n", rip.Status)
		fmt.Printf("  Region: %s\n", regionName(*rip))
		fmt.Printf("  IPv4: %s\n", dashIfEmpty(rip.IPv4Address))
		fmt.Printf("  IPv6: %s\n", dashIfEmpty(rip.IPv6Address))
		fmt.Printf("  Attached To: %s\n", attachedTo(*rip))
		fmt.Printf("  Created: %s\n", formattime.FormatTime(rip.CreatedAt.Local(), false))
		if rip.UpdatedAt != nil {
			fmt.Printf("  Updated: %s\n", formattime.FormatTime(rip.UpdatedAt.Local(), false))
		}
		return nil
	},
}

func outputYAML(rip iaas.ReservedIP) error {
	yamlData, err := yaml.Marshal(&rip)
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	fmt.Print(string(yamlData))
	return nil
}

func init() {
	ReservedIPsCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&viewOutputFormat, "output", "o", "", "Output format (yaml)")
	_ = viewCmd.RegisterFlagCompletionFunc("output", completeOutputFormat)
	viewCmd.ValidArgsFunction = completeReservedIPID
}
