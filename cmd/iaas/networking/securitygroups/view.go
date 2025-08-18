package securitygroups

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	outputFormat string
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View security group details",
	Long:    "View detailed information about a specific security group",
	Example: "tcloud networking security-groups view sg-123\ntcloud networking security-groups view sg-123 --output yaml",
	Aliases: []string{"show", "get", "describe"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		securityGroupIdentity := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		securityGroup, err := client.IaaS().GetSecurityGroup(cmd.Context(), securityGroupIdentity)
		if err != nil {
			return fmt.Errorf("failed to get security group: %w", err)
		}

		if outputFormat == "yaml" {
			return outputYAML(*securityGroup)
		}

		// Print basic information
		fmt.Printf("Security Group Details:\n")
		fmt.Printf("  ID: %s\n", securityGroup.Identity)
		fmt.Printf("  Name: %s\n", securityGroup.Name)
		fmt.Printf("  Description: %s\n", securityGroup.Description)
		fmt.Printf("  Status: %s\n", securityGroup.Status)
		fmt.Printf("  Allow Same Group Traffic: %t\n", securityGroup.AllowSameGroupTraffic)
		fmt.Printf("  Created: %s\n", formattime.FormatTime(securityGroup.CreatedAt.Local(), false))
		fmt.Printf("  Updated: %s\n", formattime.FormatTime(securityGroup.UpdatedAt.Local(), false))

		if securityGroup.Vpc != nil {
			fmt.Printf("  VPC: %s (%s)\n", securityGroup.Vpc.Name, securityGroup.Vpc.Identity)
		}

		// Print ingress rules
		fmt.Printf("\nIngress Rules (%d):\n", len(securityGroup.IngressRules))
		if len(securityGroup.IngressRules) > 0 {
			ingressBody := make([][]string, 0, len(securityGroup.IngressRules))
			for _, rule := range securityGroup.IngressRules {
				portRange := ""
				if rule.PortRangeMin > 0 && rule.PortRangeMax > 0 {
					if rule.PortRangeMin == rule.PortRangeMax {
						portRange = strconv.Itoa(int(rule.PortRangeMin))
					} else {
						portRange = fmt.Sprintf("%d-%d", rule.PortRangeMin, rule.PortRangeMax)
					}
				}

				remote := ""
				if rule.RemoteAddress != nil {
					remote = *rule.RemoteAddress
				} else if rule.RemoteSecurityGroupIdentity != nil {
					remote = *rule.RemoteSecurityGroupIdentity
				}

				ingressBody = append(ingressBody, []string{
					rule.Name,
					string(rule.Protocol),
					portRange,
					remote,
					strconv.Itoa(int(rule.Priority)),
				})
			}
			table.Print([]string{"Name", "Protocol", "Ports", "Remote", "Priority"}, ingressBody)
		}

		// Print egress rules
		fmt.Printf("\nEgress Rules (%d):\n", len(securityGroup.EgressRules))
		if len(securityGroup.EgressRules) > 0 {
			egressBody := make([][]string, 0, len(securityGroup.EgressRules))
			for _, rule := range securityGroup.EgressRules {
				portRange := ""
				if rule.PortRangeMin > 0 && rule.PortRangeMax > 0 {
					if rule.PortRangeMin == rule.PortRangeMax {
						portRange = strconv.Itoa(int(rule.PortRangeMin))
					} else {
						portRange = fmt.Sprintf("%d-%d", rule.PortRangeMin, rule.PortRangeMax)
					}
				}

				remote := ""
				if rule.RemoteAddress != nil {
					remote = *rule.RemoteAddress
				} else if rule.RemoteSecurityGroupIdentity != nil {
					remote = *rule.RemoteSecurityGroupIdentity
				}

				egressBody = append(egressBody, []string{
					rule.Name,
					string(rule.Protocol),
					portRange,
					remote,
					strconv.Itoa(int(rule.Priority)),
				})
			}
			table.Print([]string{"Name", "Protocol", "Ports", "Remote", "Priority"}, egressBody)
		}

		return nil
	},
}

// outputYAML formats the security group as YAML
func outputYAML(securityGroup iaas.SecurityGroup) error {
	// clean up the security group
	securityGroup.Organisation = nil

	// Marshal the original struct directly to YAML
	yamlData, err := yaml.Marshal(&securityGroup)
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	fmt.Print(string(yamlData))
	return nil
}

func init() {
	SecurityGroupsCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml)")

	// Add completion
	viewCmd.RegisterFlagCompletionFunc("output", completeOutputFormat)
	viewCmd.ValidArgsFunction = completeSecurityGroupID
}
