package completion

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/projectresolve"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/containerregistry"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/dns"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/kms"
	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/observability"
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

// CompleteReservedIPID provides completion for reserved IP IDs.
func CompleteReservedIPID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	ips, err := client.IaaS().ListReservedIPs(cmd.Context(), &iaas.ListReservedIPsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(ips))
	for _, ip := range ips {
		desc := ip.Name
		if ip.IPv4Address != "" {
			desc = fmt.Sprintf("%s (%s)", ip.Name, ip.IPv4Address)
		}
		completions = append(completions, ip.Identity+"\t"+desc)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteRouteTableID provides completion for route table IDs.
func CompleteRouteTableID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	tables, err := client.IaaS().ListRouteTables(cmd.Context(), &iaas.ListRouteTablesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(tables))
	for _, table := range tables {
		completions = append(completions, table.Identity+"\t"+table.Name)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteSnapshotPolicyID provides completion for snapshot policy IDs.
func CompleteSnapshotPolicyID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	policies, err := client.IaaS().ListSnapshotPolicies(cmd.Context(), &iaas.ListSnapshotPoliciesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(policies))
	for _, policy := range policies {
		completions = append(completions, policy.Identity+"\t"+policy.Name)
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

// CompleteDbBackupID provides completion for DBaaS backup IDs
func CompleteDbBackupID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	backups, err := client.DBaaS().ListDbBackupsForOrganisation(cmd.Context(), &dbaas.ListDbBackupsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, backup := range backups {
		// Skip backups that are already in args to avoid duplicates
		alreadyAdded := slices.Contains(args, backup.Identity)
		if !alreadyAdded {
			completions = append(completions, backup.Identity)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteDbBackupStoreID provides completion for DBaaS backup store IDs.
func CompleteDbBackupStoreID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	stores, err := client.DBaaS().ListDbObjectStores(cmd.Context(), &dbaas.ListDbObjectStoresRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(stores))
	for _, store := range stores {
		alreadyAdded := slices.Contains(args, store.Identity)
		if alreadyAdded {
			continue
		}
		desc := store.Name
		if store.Status != "" {
			desc = fmt.Sprintf("%s (%s)", store.Name, store.Status)
		}
		completions = append(completions, store.Identity+"\t"+desc)
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

// CompleteDbEngineVersion provides completion for DBaaS engine versions
// It requires the --engine flag to be set to determine which engine versions to return
func CompleteDbEngineVersion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Get the engine flag value
	engineFlag, err := cmd.Flags().GetString("engine")
	if err != nil || engineFlag == "" {
		// If engine is not set, return empty (user needs to set engine first)
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// Parse the engine type
	engine := dbaas.DbClusterDatabaseEngine(engineFlag)

	versions, err := client.DBaaS().ListEngineVersions(cmd.Context(), engine, &dbaas.ListEngineVersionsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	for _, v := range versions {
		desc := fmt.Sprintf("%s (%d.%d)", v.EngineVersion, v.MajorVersion, v.MinorVersion)
		completions = append(completions, v.Identity+"\t"+desc)
		completions = append(completions, v.EngineVersion+"\t"+desc)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteDbInstanceType provides completion for DBaaS instance types
func CompleteDbInstanceType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	instanceTypes, err := client.DBaaS().ListDatabaseInstanceTypes(cmd.Context(), &dbaas.ListDatabaseInstanceTypesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	for _, it := range instanceTypes {
		desc := fmt.Sprintf("%d vCPU, %d GB RAM (%s)", it.Cpus, it.Memory, it.CategorySlug)
		completions = append(completions, it.Identity+"\t"+desc)
		completions = append(completions, it.Name+"\t"+desc)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteKubernetesNodePool provides completion for Kubernetes node pools
// It requires the --cluster flag to be set to determine which node pools to return
func CompleteKubernetesNodePool(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Get the cluster flag value
	clusterFlag, err := cmd.Flags().GetString("cluster")
	if err != nil || clusterFlag == "" {
		// If cluster is not set, return empty (user needs to set cluster first)
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// Resolve cluster by identity, name, or slug
	clusters, err := client.Kubernetes().ListKubernetesClusters(cmd.Context(), &kubernetes.ListKubernetesClustersRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var cluster *kubernetes.KubernetesCluster
	for _, c := range clusters {
		if strings.EqualFold(c.Identity, clusterFlag) || strings.EqualFold(c.Name, clusterFlag) || strings.EqualFold(c.Slug, clusterFlag) {
			cluster = &c
			break
		}
	}

	if cluster == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// List node pools for the cluster
	nodePools, err := client.Kubernetes().ListKubernetesNodePools(cmd.Context(), cluster.Identity, &kubernetes.ListKubernetesNodePoolsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	for _, np := range nodePools {
		desc := fmt.Sprintf("%s (%s)", np.Name, np.Status)
		completions = append(completions, np.Identity+"\t"+desc)
		if np.Name != "" && np.Name != np.Identity {
			completions = append(completions, np.Name+"\t"+desc)
		}
		if np.Slug != "" && np.Slug != np.Identity && np.Slug != np.Name {
			completions = append(completions, np.Slug+"\t"+desc)
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

// CompleteProject provides completion for project identities and slugs.
func CompleteProject(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	projects, err := client.Projects().ListProjects(cmd.Context(), nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := []string{projectresolve.RootRef + "\tOrganisation root (no project)"}
	for _, project := range projects {
		desc := project.Name
		completions = append(completions, project.Identity+"\t"+desc)
		if project.Slug != "" && project.Slug != project.Identity {
			completions = append(completions, project.Slug+"\t"+desc)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteLoadbalancerID provides completion for load balancer IDs.
func CompleteLoadbalancerID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	loadbalancers, err := client.IaaS().ListLoadbalancers(cmd.Context(), &iaas.ListLoadbalancersRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(loadbalancers))
	for _, lb := range loadbalancers {
		desc := fmt.Sprintf("%s (%s)", lb.Name, lb.Status)
		completions = append(completions, lb.Identity+"\t"+desc)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteTargetGroupID provides completion for target group IDs.
func CompleteTargetGroupID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	targetGroups, err := client.IaaS().ListTargetGroups(cmd.Context(), &iaas.ListTargetGroupsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(targetGroups))
	for _, tg := range targetGroups {
		desc := fmt.Sprintf("%s port %d/%s", tg.Name, tg.TargetPort, tg.Protocol)
		completions = append(completions, tg.Identity+"\t"+desc)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteLoadbalancerListenerID provides completion for load balancer listener IDs.
// Requires the --loadbalancer flag to be set.
func CompleteLoadbalancerListenerID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	loadbalancerFlag, err := cmd.Flags().GetString("loadbalancer")
	if err != nil || loadbalancerFlag == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	listeners, err := client.IaaS().ListListeners(cmd.Context(), &iaas.ListLoadbalancerListenersRequest{
		Loadbalancer: loadbalancerFlag,
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(listeners))
	for _, listener := range listeners {
		desc := fmt.Sprintf("%s port %d/%s", listener.Name, listener.Port, listener.Protocol)
		completions = append(completions, listener.Identity+"\t"+desc)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteLoadbalancerProtocol provides completion for load balancer protocols.
func CompleteLoadbalancerProtocol(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		string(iaas.ProtocolTCP),
		string(iaas.ProtocolUDP),
		string(iaas.ProtocolHTTP),
		string(iaas.ProtocolHTTPS),
		string(iaas.ProtocolGRPC),
		string(iaas.ProtocolQUIC),
	}, cobra.ShellCompDirectiveNoFileComp
}

// CompleteLoadbalancingPolicy provides completion for target group load balancing policies.
func CompleteLoadbalancingPolicy(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		string(iaas.LoadbalancingPolicyRoundRobin),
		string(iaas.LoadbalancingPolicyRandom),
		string(iaas.LoadbalancingPolicyMagLev),
	}, cobra.ShellCompDirectiveNoFileComp
}

// CompleteContainerRegistryNamespaceID provides completion for container registry namespace IDs.
func CompleteContainerRegistryNamespaceID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	namespaces, err := client.ContainerRegistry().ListContainerRegistryNamespaces(cmd.Context(), &containerregistry.ListContainerRegistryNamespacesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(namespaces))
	for _, ns := range namespaces {
		completions = append(completions, ns.Identity+"\t"+ns.Namespace)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteContainerRegistryRepositoryID provides completion for container registry repository IDs.
// Requires the --namespace flag to be set.
func CompleteContainerRegistryRepositoryID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	namespaceFlag, err := cmd.Flags().GetString("namespace")
	if err != nil || namespaceFlag == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	repos, err := client.ContainerRegistry().ListContainerRegistryRepositories(cmd.Context(), namespaceFlag, &containerregistry.ListContainerRegistryRepositoriesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(repos))
	for _, repo := range repos {
		desc := repo.FullName
		if desc == "" {
			desc = repo.Image
		}
		completions = append(completions, repo.Identity+"\t"+desc)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompletePrometheusTenantID provides completion for observability workspace IDs.
// Kept for compatibility with older command names that referred to Prometheus tenants.
func CompletePrometheusTenantID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return CompleteObservabilityWorkspaceID(cmd, args, toComplete)
}

// CompleteObservabilityWorkspaceID provides completion for observability workspace IDs.
func CompleteObservabilityWorkspaceID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	workspaces, err := client.Observability().ListObservabilityWorkspaces(cmd.Context(), &observability.ListObservabilityWorkspacesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(workspaces))
	for _, workspace := range workspaces {
		desc := fmt.Sprintf("%s (%s)", workspace.Name, workspace.Status)
		completions = append(completions, workspace.Identity+"\t"+desc)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteDnsZoneIdentity provides completion for DNS zone identities.
func CompleteDnsZoneIdentity(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	zones, err := client.DNS().ListZones(cmd.Context(), &dns.ListZonesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(zones))
	for _, zone := range zones {
		completions = append(completions, zone.Identity+"\t"+zone.Name)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completionFlagString returns a flag value from local or inherited flags.
// Prefer Value.String so shell completion sees values typed earlier on the line.
func completionFlagString(cmd *cobra.Command, name string) string {
	if cmd == nil {
		return ""
	}
	if f := cmd.Flags().Lookup(name); f != nil {
		if v := strings.TrimSpace(f.Value.String()); v != "" {
			return v
		}
	}
	if f := cmd.InheritedFlags().Lookup(name); f != nil {
		if v := strings.TrimSpace(f.Value.String()); v != "" {
			return v
		}
	}
	return ""
}

// CompleteKmsKeyIdentity provides completion for KMS key identities.
// Requires the --region flag to be set.
func CompleteKmsKeyIdentity(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	regionFlag := completionFlagString(cmd, "region")
	if regionFlag == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	keys, err := client.KMS().ListKeys(cmd.Context(), regionFlag, &kms.ListKeysRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(keys))
	for _, key := range keys {
		completions = append(completions, key.Identity+"\t"+key.Name)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteKmsKeyVersion provides completion for KMS key versions.
// Requires --region and --key (or a positional key identity).
func CompleteKmsKeyVersion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	regionFlag := completionFlagString(cmd, "region")
	keyFlag := completionFlagString(cmd, "key")
	if keyFlag == "" && len(args) > 0 {
		keyFlag = strings.TrimSpace(args[0])
	}
	if regionFlag == "" || keyFlag == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	key, err := client.KMS().GetKey(cmd.Context(), regionFlag, keyFlag)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(key.Versions))
	for _, version := range key.Versions {
		desc := version.Status
		if desc == "" {
			desc = "version"
		}
		completions = append(completions, fmt.Sprintf("%d\t%s", version.Version, desc))
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
