package listeners

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
	Use:     "view LISTENER",
	Short:   "View listener details",
	Long:    "View detailed information about a load balancer listener.",
	Example: "tcloud networking loadbalancers listeners view listener-123 --loadbalancer lb-123",
	Aliases: []string{"show", "get", "describe"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if loadbalancer == "" {
			return fmt.Errorf("--loadbalancer is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		listener, err := client.IaaS().GetListener(cmd.Context(), iaas.GetLoadbalancerListenerRequest{
			Loadbalancer: loadbalancer,
			Listener:     args[0],
		})
		if err != nil {
			return fmt.Errorf("failed to get listener: %w", err)
		}

		if outputFormat == "yaml" {
			yamlData, err := yaml.Marshal(listener)
			if err != nil {
				return fmt.Errorf("failed to marshal to YAML: %w", err)
			}
			fmt.Print(string(yamlData))
			return nil
		}

		fmt.Printf("Listener Details:\n")
		fmt.Printf("  ID: %s\n", listener.Identity)
		fmt.Printf("  Name: %s\n", listener.Name)
		fmt.Printf("  Description: %s\n", listener.Description)
		fmt.Printf("  Port: %d\n", listener.Port)
		fmt.Printf("  Protocol: %s\n", listener.Protocol)
		if listener.TargetGroup != nil {
			fmt.Printf("  Target Group: %s (%s)\n", listener.TargetGroup.Name, listener.TargetGroup.Identity)
		}
		if listener.MaxConnections != nil {
			fmt.Printf("  Max Connections: %d\n", *listener.MaxConnections)
		}
		if listener.ConnectionIdleTimeout != nil {
			fmt.Printf("  Connection Idle Timeout: %d\n", *listener.ConnectionIdleTimeout)
		}
		if len(listener.AllowedSources) > 0 {
			fmt.Printf("  Allowed Sources: %v\n", listener.AllowedSources)
		}
		fmt.Printf("  Created: %s\n", formattime.FormatTime(listener.CreatedAt.Local(), false))
		fmt.Printf("  Updated: %s\n", formattime.FormatTime(listener.UpdatedAt.Local(), false))
		return nil
	},
}

func init() {
	ListenersCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVar(&loadbalancer, LoadbalancerFlag, "", "Load balancer identity")
	viewCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml)")

	viewCmd.MarkFlagRequired(LoadbalancerFlag)
	viewCmd.RegisterFlagCompletionFunc(LoadbalancerFlag, completeLoadbalancerID)
	viewCmd.RegisterFlagCompletionFunc("output", completeOutputFormat)
	viewCmd.ValidArgsFunction = completeLoadbalancerListenerID
}
