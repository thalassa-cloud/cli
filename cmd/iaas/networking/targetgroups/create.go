package targetgroups

import (
	"fmt"

	"github.com/spf13/cobra"

	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	createName                string
	createDescription         string
	createVpc                 string
	createTargetPort          int
	createProtocol            string
	createTargetSelector      []string
	createLoadbalancingPolicy string
	createEnableProxyProtocol bool
	createLabels              []string
	createAnnotations         []string
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a target group",
	Long:    "Create a new load balancer target group.",
	Example: "tcloud networking target-groups create --name web --vpc vpc-123 --port 8080 --protocol http",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if createName == "" {
			return fmt.Errorf("name is required")
		}
		if createVpc == "" {
			return fmt.Errorf("vpc is required")
		}
		if createTargetPort == 0 {
			return fmt.Errorf("port is required")
		}
		if createProtocol == "" {
			return fmt.Errorf("protocol is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		vpc, err := iaasutil.GetVPCByIdentitySlugOrName(cmd.Context(), client.IaaS(), createVpc)
		if err != nil {
			return fmt.Errorf("failed to get VPC: %w", err)
		}

		req := iaas.CreateTargetGroup{
			Name:        createName,
			Description: createDescription,
			Vpc:         vpc.Identity,
			TargetPort:  createTargetPort,
			Protocol:    iaas.LoadbalancerProtocol(createProtocol),
			Labels:      parseKeyValueSlice(createLabels),
			Annotations: parseKeyValueSlice(createAnnotations),
		}

		if len(createTargetSelector) > 0 {
			req.TargetSelector = parseKeyValueSlice(createTargetSelector)
		}
		if createLoadbalancingPolicy != "" {
			policy := iaas.LoadbalancingPolicy(createLoadbalancingPolicy)
			req.LoadbalancingPolicy = &policy
		}
		if cmd.Flags().Changed("enable-proxy-protocol") {
			req.EnableProxyProtocol = &createEnableProxyProtocol
		}

		tg, err := client.IaaS().CreateTargetGroup(cmd.Context(), req)
		if err != nil {
			return err
		}

		fmt.Printf("Target group created successfully\n")
		fmt.Printf("ID: %s\n", tg.Identity)
		fmt.Printf("Name: %s\n", tg.Name)
		return nil
	},
}

func init() {
	TargetGroupsCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createName, "name", "", "Name of the target group")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description of the target group")
	createCmd.Flags().StringVar(&createVpc, "vpc", "", "VPC identity, slug, or name")
	createCmd.Flags().IntVar(&createTargetPort, "port", 0, "Target port")
	createCmd.Flags().StringVar(&createProtocol, "protocol", "", "Target protocol (tcp, udp, http, https, grpc, quic)")
	createCmd.Flags().StringSliceVar(&createTargetSelector, "target-selector", []string{}, "Label selector for automatic target assignment (key=value)")
	createCmd.Flags().StringVar(&createLoadbalancingPolicy, "loadbalancing-policy", "", "Load balancing policy (ROUND_ROBIN, RANDOM, MAGLEV)")
	createCmd.Flags().BoolVar(&createEnableProxyProtocol, "enable-proxy-protocol", false, "Enable proxy protocol")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", []string{}, "Labels in key=value format")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", []string{}, "Annotations in key=value format")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("vpc")
	createCmd.MarkFlagRequired("port")
	createCmd.MarkFlagRequired("protocol")

	createCmd.RegisterFlagCompletionFunc("vpc", completeVPCID)
	createCmd.RegisterFlagCompletionFunc("protocol", completeLoadbalancerProtocol)
	createCmd.RegisterFlagCompletionFunc("loadbalancing-policy", completeLoadbalancingPolicy)
}
