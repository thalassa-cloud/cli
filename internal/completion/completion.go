package completion

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

// CompleteVPCID provides completion for VPC IDs
func CompleteVPCID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	vpcs, err := client.IaaS().ListVpcs(context.Background(), &iaas.ListVpcsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, vpc := range vpcs {
		completions = append(completions, vpc.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteRegion provides completion for region names
func CompleteRegion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	regions, err := client.IaaS().ListRegions(context.Background(), &iaas.ListRegionsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, region := range regions {
		completions = append(completions, region.Name)
		completions = append(completions, region.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteSubnetID provides completion for subnet IDs
func CompleteSubnetID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	subnets, err := client.IaaS().ListSubnets(context.Background(), &iaas.ListSubnetsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, subnet := range subnets {
		completions = append(completions, subnet.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteSecurityGroupID provides completion for security group IDs
func CompleteSecurityGroupID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	securityGroups, err := client.IaaS().ListSecurityGroups(context.Background(), &iaas.ListSecurityGroupsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, sg := range securityGroups {
		completions = append(completions, sg.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteNatGatewayID provides completion for NAT gateway IDs
func CompleteNatGatewayID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	natGateways, err := client.IaaS().ListNatGateways(context.Background(), &iaas.ListNatGatewaysRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, ngw := range natGateways {
		completions = append(completions, ngw.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteMachineID provides completion for machine IDs
func CompleteMachineID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	machines, err := client.IaaS().ListMachines(context.Background(), &iaas.ListMachinesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, machine := range machines {
		completions = append(completions, machine.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteOutputFormat provides completion for output format options
func CompleteOutputFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"yaml"}, cobra.ShellCompDirectiveNoFileComp
}
