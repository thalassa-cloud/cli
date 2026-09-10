package natgateways

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	createName                  string
	createDescription           string
	createSubnet                string
	createLabels                []string
	createAnnotations           []string
	createSecurityGroups        []string
	createConfigureDefaultRoute bool
	createReservedIP            string
	createWait                  bool
	createWaitTimeout           time.Duration
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a NAT gateway",
	Long:    "Create a new NAT gateway in the specified subnet.",
	Example: "tcloud networking natgateways create --name egress --subnet subnet-123\ntcloud networking natgateways create --name egress --subnet subnet-123 --configure-default-route --wait",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if createName == "" {
			return fmt.Errorf("name is required")
		}
		if createSubnet == "" {
			return fmt.Errorf("subnet is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		subnet, err := iaasutil.GetSubnetByIdentitySlugOrName(cmd.Context(), client.IaaS(), createSubnet)
		if err != nil {
			return fmt.Errorf("failed to get subnet: %w", err)
		}

		req := iaas.CreateVpcNatGateway{
			Name:                     createName,
			Description:              createDescription,
			Labels:                   parseKeyValueSlice(createLabels),
			Annotations:              parseKeyValueSlice(createAnnotations),
			SubnetIdentity:           subnet.Identity,
			SecurityGroupAttachments: createSecurityGroups,
			ConfigureDefaultRoute:    createConfigureDefaultRoute,
		}
		if createReservedIP != "" {
			req.ReservedIpID = &createReservedIP
		}

		ngw, err := client.IaaS().CreateNatGateway(cmd.Context(), req)
		if err != nil {
			return err
		}

		if createWait {
			ctxWithTimeout, cancel, err := shared.WaitContext(cmd.Context(), createWaitTimeout)
			if err != nil {
				return err
			}
			defer cancel()
			fmt.Println("Waiting for NAT gateway to have an endpoint...")
			ngw, err = client.IaaS().WaitUntilNatGatewayHasEndpoint(ctxWithTimeout, ngw.Identity)
			if err != nil {
				return fmt.Errorf("failed waiting for NAT gateway endpoint: %w", err)
			}
			fmt.Println("NAT gateway endpoint is ready")
		}

		fmt.Printf("NAT gateway created successfully\n")
		fmt.Printf("ID: %s\n", ngw.Identity)
		fmt.Printf("Name: %s\n", ngw.Name)
		fmt.Printf("Status: %s\n", ngw.Status)
		if ngw.EndpointIP != "" {
			fmt.Printf("Endpoint IP: %s\n", ngw.EndpointIP)
		}
		return nil
	},
}

func parseKeyValueSlice(items []string) map[string]string {
	result := make(map[string]string)
	for _, item := range items {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

func init() {
	NatGatewaysCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createName, "name", "", "Name of the NAT gateway")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description of the NAT gateway")
	createCmd.Flags().StringVar(&createSubnet, "subnet", "", "Subnet identity, slug, or name")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", []string{}, "Labels in key=value format")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", []string{}, "Annotations in key=value format")
	createCmd.Flags().StringSliceVar(&createSecurityGroups, "security-groups", []string{}, "Security group identities to attach")
	createCmd.Flags().BoolVar(&createConfigureDefaultRoute, "configure-default-route", false, "Configure the default route for the subnet route table")
	createCmd.Flags().StringVar(&createReservedIP, "reserved-ip", "", "Reserved IP identity to attach")
	createCmd.Flags().BoolVar(&createWait, "wait", false, "Wait for the NAT gateway to have an endpoint")
	createCmd.Flags().DurationVar(&createWaitTimeout, "wait-timeout", 20*time.Minute, "Maximum time to wait for the NAT gateway endpoint")

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("subnet")

	_ = createCmd.RegisterFlagCompletionFunc("subnet", completeSubnetID)
	_ = createCmd.RegisterFlagCompletionFunc("security-groups", completeSecurityGroupID)
	_ = createCmd.RegisterFlagCompletionFunc("reserved-ip", completeReservedIPID)
}
