package completion

import (
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

	vpcs, err := client.IaaS().ListVpcs(cmd.Context(), &iaas.ListVpcsRequest{})
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

	regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
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

	subnets, err := client.IaaS().ListSubnets(cmd.Context(), &iaas.ListSubnetsRequest{})
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

	securityGroups, err := client.IaaS().ListSecurityGroups(cmd.Context(), &iaas.ListSecurityGroupsRequest{})
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

	natGateways, err := client.IaaS().ListNatGateways(cmd.Context(), &iaas.ListNatGatewaysRequest{})
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

	machines, err := client.IaaS().ListMachines(cmd.Context(), &iaas.ListMachinesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, machine := range machines {
		completions = append(completions, machine.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func CompleteSnapshotID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	snapshots, err := client.IaaS().ListSnapshots(cmd.Context(), &iaas.ListSnapshotsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var completions []string
	for _, snapshot := range snapshots {
		completions = append(completions, snapshot.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func CompleteVolumeID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	volumes, err := client.IaaS().ListVolumes(cmd.Context(), &iaas.ListVolumesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var completions []string
	for _, volume := range volumes {
		completions = append(completions, volume.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteOutputFormat provides completion for output format options
func CompleteOutputFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"yaml"}, cobra.ShellCompDirectiveNoFileComp
}
