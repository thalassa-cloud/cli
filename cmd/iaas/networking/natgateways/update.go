package natgateways

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateName           string
	updateDescription    string
	updateLabels         []string
	updateAnnotations    []string
	updateSecurityGroups []string
	updateReservedIP     string
	updateDetachReserved bool
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a NAT gateway",
	Long:    "Update properties of an existing NAT gateway. Unspecified fields are preserved.",
	Example: "tcloud networking natgateways update ngw-123 --name egress-prod\ntcloud networking natgateways update ngw-123 --reserved-ip rip-456\ntcloud networking natgateways update ngw-123 --detach-reserved-ip",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateDetachReserved && cmd.Flags().Changed("reserved-ip") {
			return fmt.Errorf("--detach-reserved-ip and --reserved-ip are mutually exclusive")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.IaaS().GetNatGateway(cmd.Context(), args[0])
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("NAT gateway not found: %s", args[0])
			}
			return fmt.Errorf("failed to get NAT gateway: %w", err)
		}

		req := iaas.UpdateVpcNatGateway{
			Name:        current.Name,
			Description: current.Description,
			Labels:      current.Labels,
			Annotations: current.Annotations,
		}
		if len(current.SecurityGroups) > 0 {
			for _, sg := range current.SecurityGroups {
				req.SecurityGroupAttachments = append(req.SecurityGroupAttachments, sg.Identity)
			}
		}

		if cmd.Flags().Changed("name") {
			req.Name = updateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = parseKeyValueSlice(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = parseKeyValueSlice(updateAnnotations)
		}
		if cmd.Flags().Changed("security-groups") {
			req.SecurityGroupAttachments = updateSecurityGroups
		}
		if updateDetachReserved {
			empty := ""
			req.ReservedIpID = &empty
		} else if cmd.Flags().Changed("reserved-ip") {
			req.ReservedIpID = &updateReservedIP
		}

		ngw, err := client.IaaS().UpdateNatGateway(cmd.Context(), current.Identity, req)
		if err != nil {
			return err
		}

		fmt.Printf("NAT gateway updated successfully\n")
		fmt.Printf("ID: %s\n", ngw.Identity)
		fmt.Printf("Name: %s\n", ngw.Name)
		fmt.Printf("Status: %s\n", ngw.Status)
		return nil
	},
}

func init() {
	NatGatewaysCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateName, "name", "", "Name of the NAT gateway")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description of the NAT gateway")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", []string{}, "Labels in key=value format")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", []string{}, "Annotations in key=value format")
	updateCmd.Flags().StringSliceVar(&updateSecurityGroups, "security-groups", []string{}, "Security group identities to attach")
	updateCmd.Flags().StringVar(&updateReservedIP, "reserved-ip", "", "Reserved IP identity to attach or replace")
	updateCmd.Flags().BoolVar(&updateDetachReserved, "detach-reserved-ip", false, "Detach the currently associated reserved IP")

	updateCmd.ValidArgsFunction = completeNatGatewayID
	_ = updateCmd.RegisterFlagCompletionFunc("security-groups", completeSecurityGroupID)
	_ = updateCmd.RegisterFlagCompletionFunc("reserved-ip", completeReservedIPID)
}
