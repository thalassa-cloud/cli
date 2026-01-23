package completion

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/tfs"
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

func CompleteTfsInstanceID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	instances, err := client.Tfs().ListTfsInstances(cmd.Context(), &tfs.ListTfsInstancesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var completions []string
	for _, instance := range instances {
		completions = append(completions, instance.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func CompleteDbClusterID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	clusters, err := client.DBaaS().ListDbClusters(cmd.Context(), &dbaas.ListDbClustersRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var completions []string
	for _, cluster := range clusters {
		completions = append(completions, cluster.Identity)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteOutputFormat provides completion for output format options
func CompleteOutputFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"yaml"}, cobra.ShellCompDirectiveNoFileComp
}

// CompleteRegionEnhanced provides enhanced completion for region names with identity, slug, and tab formatting
func CompleteRegionEnhanced(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(regions))
	for _, r := range regions {
		completions = append(completions, r.Identity+"\t"+r.Name)
		if r.Slug != "" {
			completions = append(completions, r.Slug+"\t"+r.Name)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteSubnetEnhanced provides enhanced completion for subnet IDs with identity, slug, and tab formatting
func CompleteSubnetEnhanced(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	subnets, err := client.IaaS().ListSubnets(cmd.Context(), &iaas.ListSubnetsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(subnets))
	for _, s := range subnets {
		completions = append(completions, s.Identity+"\t"+s.Name+" ("+s.Cidr+")")
		if s.Slug != "" {
			completions = append(completions, s.Slug+"\t"+s.Name+" ("+s.Cidr+")")
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteMachineType provides completion for machine types with descriptions
func CompleteMachineType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	machineTypeCategories, err := client.IaaS().ListMachineTypeCategories(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	for _, category := range machineTypeCategories {
		for _, mt := range category.MachineTypes {
			desc := fmt.Sprintf("%d vCPU, %d MB RAM", mt.Vcpus, mt.RamMb)
			completions = append(completions, mt.Name+"\t"+desc+" ("+category.Name+")")
			if mt.Slug != "" {
				completions = append(completions, mt.Slug+"\t"+desc+" ("+category.Name+")")
			}
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteKubernetesVersion provides completion for Kubernetes versions
func CompleteKubernetesVersion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	versions, err := client.Kubernetes().ListKubernetesVersions(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	for _, v := range versions {
		if !v.Enabled {
			continue
		}
		completions = append(completions, v.Identity+"\t"+v.Name)
		if v.Slug != "" {
			completions = append(completions, v.Slug+"\t"+v.Name)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteKubernetesCluster provides completion for Kubernetes cluster identities, names, and slugs
func CompleteKubernetesCluster(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	clusters, err := client.Kubernetes().ListKubernetesClusters(cmd.Context(), &kubernetes.ListKubernetesClustersRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	for _, c := range clusters {
		desc := fmt.Sprintf("%s (%s)", c.Name, c.Status)
		completions = append(completions, c.Identity+"\t"+desc)
		if c.Name != "" && c.Name != c.Identity {
			completions = append(completions, c.Name+"\t"+desc)
		}
		if c.Slug != "" && c.Slug != c.Identity && c.Slug != c.Name {
			completions = append(completions, c.Slug+"\t"+desc)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteVpcPeeringConnectionID provides completion for VPC peering connection identities, names, and slugs
func CompleteVpcPeeringConnectionID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	connections, err := client.IaaS().ListVpcPeeringConnections(cmd.Context(), &iaas.ListVpcPeeringConnectionsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	for _, conn := range connections {
		desc := fmt.Sprintf("%s (%s)", conn.Name, conn.Status)
		completions = append(completions, conn.Identity+"\t"+desc)
		if conn.Name != "" && conn.Name != conn.Identity {
			completions = append(completions, conn.Name+"\t"+desc)
		}
		if conn.Slug != "" && conn.Slug != conn.Identity && conn.Slug != conn.Name {
			completions = append(completions, conn.Slug+"\t"+desc)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteOrganisation provides completion for organisation identities and slugs
func CompleteOrganisation(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	organisations, err := client.Me().ListMyOrganisations(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	for _, org := range organisations {
		desc := org.Name
		completions = append(completions, org.Identity+"\t"+desc)
		if org.Slug != "" && org.Slug != org.Identity {
			completions = append(completions, org.Slug+"\t"+desc)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
