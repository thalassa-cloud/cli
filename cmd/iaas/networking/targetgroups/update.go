package targetgroups

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateName                string
	updateDescription         string
	updateTargetPort          int
	updateProtocol            string
	updateTargetSelector      []string
	updateLoadbalancingPolicy string
	updateEnableProxyProtocol bool
	updateLabels              []string
	updateAnnotations         []string
)

var updateCmd = &cobra.Command{
	Use:     "update TARGET_GROUP",
	Short:   "Update a target group",
	Long:    "Update properties of an existing target group.",
	Example: "tcloud networking target-groups update tg-123 --name web-prod\ntcloud networking target-groups update tg-123 --port 8443",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.IaaS().GetTargetGroup(cmd.Context(), iaas.GetTargetGroupRequest{
			Identity: args[0],
		})
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("target group not found: %s", args[0])
			}
			return fmt.Errorf("failed to get target group: %w", err)
		}

		req := iaas.UpdateTargetGroup{
			Name:            current.Name,
			Description:     current.Description,
			TargetPort:      current.TargetPort,
			Protocol:        current.Protocol,
			TargetSelector:  current.TargetSelector,
			Labels:          current.Labels,
			Annotations:     current.Annotations,
			LoadbalancingPolicy: current.LoadbalancingPolicy,
			EnableProxyProtocol: current.EnableProxyProtocol,
			HealthCheck:         current.HealthCheck,
		}

		if cmd.Flags().Changed("name") {
			req.Name = updateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("port") {
			req.TargetPort = updateTargetPort
		}
		if cmd.Flags().Changed("protocol") {
			req.Protocol = iaas.LoadbalancerProtocol(updateProtocol)
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = parseKeyValueSlice(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = parseKeyValueSlice(updateAnnotations)
		}
		if cmd.Flags().Changed("target-selector") {
			req.TargetSelector = parseKeyValueSlice(updateTargetSelector)
		}
		if cmd.Flags().Changed("loadbalancing-policy") {
			policy := iaas.LoadbalancingPolicy(updateLoadbalancingPolicy)
			req.LoadbalancingPolicy = &policy
		}
		if cmd.Flags().Changed("enable-proxy-protocol") {
			req.EnableProxyProtocol = &updateEnableProxyProtocol
		}

		tg, err := client.IaaS().UpdateTargetGroup(cmd.Context(), iaas.UpdateTargetGroupRequest{
			Identity:          current.Identity,
			UpdateTargetGroup: req,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Target group updated successfully\n")
		fmt.Printf("ID: %s\n", tg.Identity)
		fmt.Printf("Name: %s\n", tg.Name)
		return nil
	},
}

func init() {
	TargetGroupsCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateName, "name", "", "Name of the target group")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description of the target group")
	updateCmd.Flags().IntVar(&updateTargetPort, "port", 0, "Target port")
	updateCmd.Flags().StringVar(&updateProtocol, "protocol", "", "Target protocol")
	updateCmd.Flags().StringSliceVar(&updateTargetSelector, "target-selector", []string{}, "Label selector for automatic target assignment (key=value)")
	updateCmd.Flags().StringVar(&updateLoadbalancingPolicy, "loadbalancing-policy", "", "Load balancing policy (ROUND_ROBIN, RANDOM, MAGLEV)")
	updateCmd.Flags().BoolVar(&updateEnableProxyProtocol, "enable-proxy-protocol", false, "Enable proxy protocol")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", []string{}, "Labels in key=value format")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", []string{}, "Annotations in key=value format")

	updateCmd.ValidArgsFunction = completeTargetGroupID
	updateCmd.RegisterFlagCompletionFunc("protocol", completeLoadbalancerProtocol)
	updateCmd.RegisterFlagCompletionFunc("loadbalancing-policy", completeLoadbalancingPolicy)
}
