package listeners

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	createName                  string
	createDescription           string
	createPort                  int
	createProtocol              string
	createTargetGroup           string
	createMaxConnections        uint32
	createConnectionIdleTimeout uint32
	createAllowedSources        []string
	createLabels                []string
	createAnnotations           []string
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a listener",
	Long:    "Create a listener on a load balancer.",
	Example: "tcloud networking loadbalancers listeners create --loadbalancer lb-123 --name http --port 80 --protocol tcp --target-group tg-123",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if loadbalancer == "" {
			return fmt.Errorf("--loadbalancer is required")
		}
		if createName == "" {
			return fmt.Errorf("name is required")
		}
		if createPort == 0 {
			return fmt.Errorf("port is required")
		}
		if createProtocol == "" {
			return fmt.Errorf("protocol is required")
		}
		if createTargetGroup == "" {
			return fmt.Errorf("target-group is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req := iaas.CreateListener{
			Name:           createName,
			Description:    createDescription,
			Port:           createPort,
			Protocol:       iaas.LoadbalancerProtocol(createProtocol),
			TargetGroup:    createTargetGroup,
			AllowedSources: createAllowedSources,
			Labels:         parseKeyValueSlice(createLabels),
			Annotations:    parseKeyValueSlice(createAnnotations),
		}
		if cmd.Flags().Changed("max-connections") {
			req.MaxConnections = &createMaxConnections
		}
		if cmd.Flags().Changed("connection-idle-timeout") {
			req.ConnectionIdleTimeout = &createConnectionIdleTimeout
		}

		listener, err := client.IaaS().CreateListener(cmd.Context(), loadbalancer, req)
		if err != nil {
			return err
		}

		fmt.Printf("Listener created successfully\n")
		fmt.Printf("ID: %s\n", listener.Identity)
		fmt.Printf("Name: %s\n", listener.Name)
		fmt.Printf("Port: %d/%s\n", listener.Port, listener.Protocol)
		return nil
	},
}

func init() {
	ListenersCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&loadbalancer, LoadbalancerFlag, "", "Load balancer identity")
	createCmd.Flags().StringVar(&createName, "name", "", "Name of the listener")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description of the listener")
	createCmd.Flags().IntVar(&createPort, "port", 0, "Listener port")
	createCmd.Flags().StringVar(&createProtocol, "protocol", "", "Listener protocol (tcp, udp)")
	createCmd.Flags().StringVar(&createTargetGroup, "target-group", "", "Target group identity")
	createCmd.Flags().Uint32Var(&createMaxConnections, "max-connections", 0, "Maximum connections")
	createCmd.Flags().Uint32Var(&createConnectionIdleTimeout, "connection-idle-timeout", 0, "Connection idle timeout in seconds")
	createCmd.Flags().StringSliceVar(&createAllowedSources, "allowed-sources", []string{}, "Allowed source CIDR blocks")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", []string{}, "Labels in key=value format")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", []string{}, "Annotations in key=value format")

	createCmd.MarkFlagRequired(LoadbalancerFlag)
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("port")
	createCmd.MarkFlagRequired("protocol")
	createCmd.MarkFlagRequired("target-group")

	createCmd.RegisterFlagCompletionFunc(LoadbalancerFlag, completeLoadbalancerID)
	createCmd.RegisterFlagCompletionFunc("protocol", completeLoadbalancerProtocol)
	createCmd.RegisterFlagCompletionFunc("target-group", completeTargetGroupID)
}
