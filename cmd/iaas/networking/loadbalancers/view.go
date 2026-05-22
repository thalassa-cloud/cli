package loadbalancers

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	"gopkg.in/yaml.v3"
)

var outputFormat string

var viewCmd = &cobra.Command{
	Use:     "view LOADBALANCER",
	Short:   "View load balancer details",
	Long:    "View detailed information about a load balancer.",
	Example: "tcloud networking loadbalancers view lb-123\ntcloud networking loadbalancers view lb-123 --output yaml",
	Aliases: []string{"show", "get", "describe"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		lb, err := client.IaaS().GetLoadbalancer(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get load balancer: %w", err)
		}

		if outputFormat == "yaml" {
			return outputLoadbalancerYAML(*lb)
		}

		fmt.Printf("Load Balancer Details:\n")
		fmt.Printf("  ID: %s\n", lb.Identity)
		fmt.Printf("  Name: %s\n", lb.Name)
		fmt.Printf("  Description: %s\n", lb.Description)
		fmt.Printf("  Status: %s\n", lb.Status)
		fmt.Printf("  Hostname: %s\n", lb.Hostname)
		fmt.Printf("  Delete Protection: %t\n", lb.DeleteProtection)
		fmt.Printf("  External IPs: %s\n", joinStrings(lb.ExternalIpAddresses))
		fmt.Printf("  Internal IPs: %s\n", joinStrings(lb.InternalIpAddresses))
		fmt.Printf("  Created: %s\n", formattime.FormatTime(lb.CreatedAt.Local(), false))
		fmt.Printf("  Updated: %s\n", formattime.FormatTime(lb.UpdatedAt.Local(), false))

		if lb.Vpc != nil {
			fmt.Printf("  VPC: %s (%s)\n", lb.Vpc.Name, lb.Vpc.Identity)
		}
		if lb.Subnet != nil {
			fmt.Printf("  Subnet: %s (%s)\n", lb.Subnet.Name, lb.Subnet.Identity)
		}
		if len(lb.SecurityGroups) > 0 {
			fmt.Printf("  Security Groups:\n")
			for _, sg := range lb.SecurityGroups {
				fmt.Printf("    - %s (%s)\n", sg.Name, sg.Identity)
			}
		}
		if len(lb.LoadbalancerListeners) > 0 {
			fmt.Printf("  Listeners:\n")
			for _, listener := range lb.LoadbalancerListeners {
				targetGroup := "-"
				if listener.TargetGroup != nil {
					targetGroup = listener.TargetGroup.Name
				}
				fmt.Printf("    - %s (%s) port %d/%s -> %s\n",
					listener.Name, listener.Identity, listener.Port, listener.Protocol, targetGroup)
			}
		}

		return nil
	},
}

func outputLoadbalancerYAML(lb iaas.VpcLoadbalancer) error {
	lb.Organisation = nil
	yamlData, err := yaml.Marshal(&lb)
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	fmt.Print(string(yamlData))
	return nil
}

func init() {
	LoadbalancersCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml)")
	viewCmd.RegisterFlagCompletionFunc("output", completeOutputFormat)
	viewCmd.ValidArgsFunction = completeLoadbalancerID
}
