package listeners

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateName                  string
	updateDescription           string
	updatePort                  int
	updateProtocol              string
	updateTargetGroup           string
	updateMaxConnections        uint32
	updateConnectionIdleTimeout uint32
	updateAllowedSources        []string
	updateLabels                []string
	updateAnnotations           []string
)

var updateCmd = &cobra.Command{
	Use:     "update LISTENER",
	Short:   "Update a listener",
	Long:    "Update properties of an existing load balancer listener.",
	Example: "tcloud networking loadbalancers listeners update listener-123 --loadbalancer lb-123 --port 8080",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if loadbalancer == "" {
			return fmt.Errorf("--loadbalancer is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.IaaS().GetListener(cmd.Context(), iaas.GetLoadbalancerListenerRequest{
			Loadbalancer: loadbalancer,
			Listener:     args[0],
		})
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("listener not found: %s", args[0])
			}
			return fmt.Errorf("failed to get listener: %w", err)
		}

		req := iaas.UpdateListener{
			Name:            current.Name,
			Description:     current.Description,
			Port:            current.Port,
			Protocol:        current.Protocol,
			AllowedSources:  current.AllowedSources,
			Labels:          current.Labels,
			Annotations:     current.Annotations,
			MaxConnections:  current.MaxConnections,
			ConnectionIdleTimeout: current.ConnectionIdleTimeout,
		}
		if current.TargetGroup != nil {
			req.TargetGroup = current.TargetGroup.Identity
		}

		if cmd.Flags().Changed("name") {
			req.Name = updateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("port") {
			req.Port = updatePort
		}
		if cmd.Flags().Changed("protocol") {
			req.Protocol = iaas.LoadbalancerProtocol(updateProtocol)
		}
		if cmd.Flags().Changed("target-group") {
			req.TargetGroup = updateTargetGroup
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = parseKeyValueSlice(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = parseKeyValueSlice(updateAnnotations)
		}
		if cmd.Flags().Changed("allowed-sources") {
			req.AllowedSources = updateAllowedSources
		}
		if cmd.Flags().Changed("max-connections") {
			req.MaxConnections = &updateMaxConnections
		}
		if cmd.Flags().Changed("connection-idle-timeout") {
			req.ConnectionIdleTimeout = &updateConnectionIdleTimeout
		}

		listener, err := client.IaaS().UpdateListener(cmd.Context(), loadbalancer, current.Identity, req)
		if err != nil {
			return err
		}

		fmt.Printf("Listener updated successfully\n")
		fmt.Printf("ID: %s\n", listener.Identity)
		fmt.Printf("Name: %s\n", listener.Name)
		return nil
	},
}

func init() {
	ListenersCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&loadbalancer, LoadbalancerFlag, "", "Load balancer identity")
	updateCmd.Flags().StringVar(&updateName, "name", "", "Name of the listener")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description of the listener")
	updateCmd.Flags().IntVar(&updatePort, "port", 0, "Listener port")
	updateCmd.Flags().StringVar(&updateProtocol, "protocol", "", "Listener protocol")
	updateCmd.Flags().StringVar(&updateTargetGroup, "target-group", "", "Target group identity")
	updateCmd.Flags().Uint32Var(&updateMaxConnections, "max-connections", 0, "Maximum connections")
	updateCmd.Flags().Uint32Var(&updateConnectionIdleTimeout, "connection-idle-timeout", 0, "Connection idle timeout in seconds")
	updateCmd.Flags().StringSliceVar(&updateAllowedSources, "allowed-sources", []string{}, "Allowed source CIDR blocks")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", []string{}, "Labels in key=value format")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", []string{}, "Annotations in key=value format")

	updateCmd.MarkFlagRequired(LoadbalancerFlag)
	updateCmd.ValidArgsFunction = completeLoadbalancerListenerID
	updateCmd.RegisterFlagCompletionFunc(LoadbalancerFlag, completeLoadbalancerID)
	updateCmd.RegisterFlagCompletionFunc("protocol", completeLoadbalancerProtocol)
	updateCmd.RegisterFlagCompletionFunc("target-group", completeTargetGroupID)
}
