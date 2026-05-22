package loadbalancers

import (
	"fmt"

	"github.com/spf13/cobra"

	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateName             string
	updateDescription      string
	updateSubnet           string
	updateDeleteProtection bool
	updateLabels           []string
	updateAnnotations      []string
	updateSecurityGroups   []string
)

var updateCmd = &cobra.Command{
	Use:     "update LOADBALANCER",
	Short:   "Update a load balancer",
	Long:    "Update properties of an existing load balancer.",
	Example: "tcloud networking loadbalancers update lb-123 --name web-prod\ntcloud networking loadbalancers update lb-123 --delete-protection",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.IaaS().GetLoadbalancer(cmd.Context(), args[0])
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("load balancer not found: %s", args[0])
			}
			return fmt.Errorf("failed to get load balancer: %w", err)
		}

		req := iaas.UpdateLoadbalancer{
			Name:             current.Name,
			Description:      current.Description,
			Labels:           current.Labels,
			Annotations:      current.Annotations,
			DeleteProtection: current.DeleteProtection,
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
		if cmd.Flags().Changed("delete-protection") {
			req.DeleteProtection = updateDeleteProtection
		}
		if cmd.Flags().Changed("subnet") {
			subnet, err := iaasutil.GetSubnetByIdentitySlugOrName(cmd.Context(), client.IaaS(), updateSubnet)
			if err != nil {
				return fmt.Errorf("failed to get subnet: %w", err)
			}
			req.Subnet = &subnet.Identity
		}
		if cmd.Flags().Changed("security-groups") {
			req.SecurityGroupAttachments = updateSecurityGroups
		}

		lb, err := client.IaaS().UpdateLoadbalancer(cmd.Context(), current.Identity, req)
		if err != nil {
			return err
		}

		fmt.Printf("Load balancer updated successfully\n")
		fmt.Printf("ID: %s\n", lb.Identity)
		fmt.Printf("Name: %s\n", lb.Name)
		fmt.Printf("Status: %s\n", lb.Status)
		return nil
	},
}

func init() {
	LoadbalancersCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateName, "name", "", "Name of the load balancer")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description of the load balancer")
	updateCmd.Flags().StringVar(&updateSubnet, "subnet", "", "Subnet identity, slug, or name")
	updateCmd.Flags().BoolVar(&updateDeleteProtection, "delete-protection", false, "Enable delete protection")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", []string{}, "Labels in key=value format")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", []string{}, "Annotations in key=value format")
	updateCmd.Flags().StringSliceVar(&updateSecurityGroups, "security-groups", []string{}, "Security group identities to attach")

	updateCmd.ValidArgsFunction = completeLoadbalancerID
	updateCmd.RegisterFlagCompletionFunc("subnet", completeSubnetID)
	updateCmd.RegisterFlagCompletionFunc("security-groups", completeSecurityGroupID)
}
