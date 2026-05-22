package targetgroups

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
	Use:     "view TARGET_GROUP",
	Short:   "View target group details",
	Long:    "View detailed information about a target group.",
	Example: "tcloud networking target-groups view tg-123\ntcloud networking target-groups view tg-123 --output yaml",
	Aliases: []string{"show", "get", "describe"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		tg, err := client.IaaS().GetTargetGroup(cmd.Context(), iaas.GetTargetGroupRequest{
			Identity: args[0],
		})
		if err != nil {
			return fmt.Errorf("failed to get target group: %w", err)
		}

		if outputFormat == "yaml" {
			tg.Organisation = nil
			yamlData, err := yaml.Marshal(tg)
			if err != nil {
				return fmt.Errorf("failed to marshal to YAML: %w", err)
			}
			fmt.Print(string(yamlData))
			return nil
		}

		fmt.Printf("Target Group Details:\n")
		fmt.Printf("  ID: %s\n", tg.Identity)
		fmt.Printf("  Name: %s\n", tg.Name)
		fmt.Printf("  Description: %s\n", tg.Description)
		fmt.Printf("  Target Port: %d\n", tg.TargetPort)
		fmt.Printf("  Protocol: %s\n", tg.Protocol)
		if tg.LoadbalancingPolicy != nil {
			fmt.Printf("  Load Balancing Policy: %s\n", *tg.LoadbalancingPolicy)
		}
		if tg.EnableProxyProtocol != nil {
			fmt.Printf("  Proxy Protocol: %t\n", *tg.EnableProxyProtocol)
		}
		if tg.Vpc != nil {
			fmt.Printf("  VPC: %s (%s)\n", tg.Vpc.Name, tg.Vpc.Identity)
		}
		if len(tg.TargetSelector) > 0 {
			fmt.Printf("  Target Selector: %v\n", tg.TargetSelector)
		}
		if len(tg.LoadbalancerListeners) > 0 {
			fmt.Printf("  Listeners:\n")
			for _, listener := range tg.LoadbalancerListeners {
				fmt.Printf("    - %s (%s)\n", listener.Name, listener.Identity)
			}
		}
		if len(tg.LoadbalancerTargetGroupAttachments) > 0 {
			fmt.Printf("  Attachments:\n")
			for _, attachment := range tg.LoadbalancerTargetGroupAttachments {
				target := attachment.Identity
				if attachment.VirtualMachineInstance != nil {
					target = attachment.VirtualMachineInstance.Name + " (" + attachment.VirtualMachineInstance.Identity + ")"
				} else if attachment.Endpoint != nil {
					target = attachment.Endpoint.Name + " (" + attachment.Endpoint.Identity + ")"
				}
				fmt.Printf("    - %s -> %s\n", attachment.Identity, target)
			}
		}
		fmt.Printf("  Created: %s\n", formattime.FormatTime(tg.CreatedAt.Local(), false))
		fmt.Printf("  Updated: %s\n", formattime.FormatTime(tg.UpdatedAt.Local(), false))
		return nil
	},
}

func init() {
	TargetGroupsCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml)")
	viewCmd.RegisterFlagCompletionFunc("output", completeOutputFormat)
	viewCmd.ValidArgsFunction = completeTargetGroupID
}
