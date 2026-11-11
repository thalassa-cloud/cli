package nodepools

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/thalassa"
)

// NodePoolConfig holds configuration for creating a node pool
type NodePoolConfig struct {
	Name           string
	MachineType    string
	Replicas       int
	EnableAS       bool
	MinNodes       int
	MaxNodes       int
	Subnet         string
	AZs            []string
	EnableAH       bool
	UpgradeStrat   string
	Labels         []string
	Annotations    []string
	Taints         []string
	SecurityGroups []string
}

// ResolveMachineType resolves a machine type identifier to its name
func ResolveMachineType(ctx context.Context, client thalassa.Client, machineTypeIdentifier string) (string, error) {
	machineTypeCategories, err := client.IaaS().ListMachineTypeCategories(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list machine types: %w", err)
	}

	for _, category := range machineTypeCategories {
		for _, mt := range category.MachineTypes {
			if strings.EqualFold(mt.Name, machineTypeIdentifier) ||
				strings.EqualFold(mt.Slug, machineTypeIdentifier) ||
				strings.EqualFold(mt.Identity, machineTypeIdentifier) {
				return mt.Name, nil
			}
		}
	}

	return "", fmt.Errorf("machine type not found: %s", machineTypeIdentifier)
}

// ValidateAutoscalingConfig validates autoscaling configuration
func ValidateAutoscalingConfig(enableAS bool, minNodes, maxNodes, replicas int) error {
	if enableAS {
		if minNodes < 0 {
			return fmt.Errorf("min-nodes must be at least 0")
		}
		if maxNodes < 1 {
			return fmt.Errorf("max-nodes must be at least 1")
		}
		if minNodes > maxNodes {
			return fmt.Errorf("min-nodes (%d) cannot be greater than max-nodes (%d)", minNodes, maxNodes)
		}
	} else {
		if replicas < 1 {
			return fmt.Errorf("num-nodes must be at least 1 when autoscaling is disabled")
		}
	}
	return nil
}

// ResolveNodePoolSubnet resolves the subnet for a node pool
func ResolveNodePoolSubnet(ctx context.Context, client thalassa.Client, subnetIdentifier string, cluster *kubernetes.KubernetesCluster) (string, error) {
	if subnetIdentifier != "" {
		subnets, err := client.IaaS().ListSubnets(ctx, &iaas.ListSubnetsRequest{})
		if err != nil {
			return "", fmt.Errorf("failed to list subnets: %w", err)
		}

		for _, s := range subnets {
			if strings.EqualFold(s.Identity, subnetIdentifier) ||
				strings.EqualFold(s.Slug, subnetIdentifier) ||
				strings.EqualFold(s.Name, subnetIdentifier) {
				return s.Identity, nil
			}
		}

		return "", fmt.Errorf("node pool subnet not found: %s", subnetIdentifier)
	}

	// Use cluster subnet if available (for managed clusters)
	if cluster.Subnet != nil {
		return cluster.Subnet.Identity, nil
	}

	return "", fmt.Errorf("node-pool-subnet is required for hosted-control-plane clusters")
}

// DetermineAvailabilityZones determines availability zones for node pools
func DetermineAvailabilityZones(ctx context.Context, client thalassa.Client, requestedAZs []string, cluster *kubernetes.KubernetesCluster) ([]string, error) {
	if len(requestedAZs) > 0 {
		return requestedAZs, nil
	}

	// Get availability zones from the region
	var regionIdentity string
	if cluster.Region != nil {
		regionIdentity = cluster.Region.Identity
	} else {
		return nil, fmt.Errorf("cannot determine availability zones: cluster region is not available")
	}

	// Get region details to find available zones
	regions, err := client.IaaS().ListRegions(ctx, &iaas.ListRegionsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list regions: %w", err)
	}

	var availableZones []string
	for _, r := range regions {
		if r.Identity == regionIdentity {
			for _, zone := range r.Zones {
				if zone.Name != "" {
					availableZones = append(availableZones, zone.Name)
				}
			}
			break
		}
	}

	if len(availableZones) == 0 {
		return nil, fmt.Errorf("no availability zones found for region %s", regionIdentity)
	}

	// Select a random availability zone
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	selectedAZ := availableZones[rng.Intn(len(availableZones))]

	return []string{selectedAZ}, nil
}

// CalculateReplicas calculates the number of replicas based on autoscaling configuration
func CalculateReplicas(enableAS bool, replicas, minNodes int) int {
	if enableAS {
		return minNodes
	}
	return replicas
}

// ParseLabels parses label strings into a map
func ParseLabels(labelStrings []string) map[string]string {
	labels := make(map[string]string)
	for _, label := range labelStrings {
		parts := strings.SplitN(label, "=", 2)
		if len(parts) == 2 {
			labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return labels
}

// ParseAnnotations parses annotation strings into a map
func ParseAnnotations(annotationStrings []string) map[string]string {
	annotations := make(map[string]string)
	for _, annotation := range annotationStrings {
		parts := strings.SplitN(annotation, "=", 2)
		if len(parts) == 2 {
			annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return annotations
}

// ParseTaints parses taint strings into NodeTaint structs
func ParseTaints(taintStrings []string) ([]kubernetes.NodeTaint, error) {
	nodeTaints := []kubernetes.NodeTaint{}
	for _, taint := range taintStrings {
		// Format: key=value:effect or key:effect
		parts := strings.SplitN(taint, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid taint format: %s (expected key=value:effect or key:effect)", taint)
		}

		keyValue := strings.TrimSpace(parts[0])
		effect := strings.TrimSpace(parts[1])

		keyParts := strings.SplitN(keyValue, "=", 2)
		key := keyParts[0]
		value := ""
		if len(keyParts) == 2 {
			value = keyParts[1]
		}

		operator := "Equal"
		if value == "" {
			operator = "Exists"
		}

		nodeTaints = append(nodeTaints, kubernetes.NodeTaint{
			Key:      key,
			Value:    value,
			Effect:   effect,
			Operator: operator,
		})
	}
	return nodeTaints, nil
}

// ParseUpgradeStrategy parses upgrade strategy string
func ParseUpgradeStrategy(strategy string) kubernetes.KubernetesNodePoolUpgradeStrategy {
	if strategy == "" {
		return kubernetes.KubernetesNodePoolUpgradeStrategyAuto
	}
	return kubernetes.KubernetesNodePoolUpgradeStrategy(strategy)
}

// NodePoolUpdateOverrides holds fields that replace the current node pool values.
// Nil pointers and false replace flags keep the current value.
type NodePoolUpdateOverrides struct {
	MachineType               *string
	Replicas                  *int
	MinReplicas               *int
	MaxReplicas               *int
	EnableAutoscaling         *bool
	EnableAutoHealing         *bool
	UpgradeStrategy           *kubernetes.KubernetesNodePoolUpgradeStrategy
	KubernetesVersionIdentity *string
	Taints                    []kubernetes.NodeTaint
	ReplaceTaints             bool
	SecurityGroupAttachments  []string
	ReplaceSecurityGroups     bool
}

// FindCluster returns the cluster matching an identity, name, or slug.
func FindCluster(clusters []kubernetes.KubernetesCluster, ref string) (*kubernetes.KubernetesCluster, error) {
	for _, cluster := range clusters {
		if strings.EqualFold(cluster.Identity, ref) ||
			strings.EqualFold(cluster.Name, ref) ||
			strings.EqualFold(cluster.Slug, ref) {
			match := cluster
			return &match, nil
		}
	}
	return nil, fmt.Errorf("cluster not found: %s", ref)
}

// FindNodePool returns the node pool matching an identity, name, or slug.
func FindNodePool(nodePools []kubernetes.KubernetesNodePool, ref string) (*kubernetes.KubernetesNodePool, error) {
	for _, nodePool := range nodePools {
		if strings.EqualFold(nodePool.Identity, ref) ||
			strings.EqualFold(nodePool.Name, ref) ||
			strings.EqualFold(nodePool.Slug, ref) {
			match := nodePool
			return &match, nil
		}
	}
	return nil, fmt.Errorf("node pool not found: %s", ref)
}

// ResolveKubernetesVersion resolves an enabled Kubernetes version identity, slug, name, or version string.
func ResolveKubernetesVersion(ctx context.Context, client thalassa.Client, versionIdentifier string) (string, error) {
	versions, err := client.Kubernetes().ListKubernetesVersions(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list kubernetes versions: %w", err)
	}

	for _, version := range versions {
		if !version.Enabled {
			continue
		}
		if strings.EqualFold(version.Identity, versionIdentifier) ||
			strings.EqualFold(version.Slug, versionIdentifier) ||
			strings.EqualFold(version.Name, versionIdentifier) ||
			strings.EqualFold(version.KubernetesVersion, versionIdentifier) {
			if version.Identity == "" {
				return "", fmt.Errorf("kubernetes version %s has no identity", versionIdentifier)
			}
			return version.Identity, nil
		}
	}

	return "", fmt.Errorf("kubernetes version not found or not enabled: %s", versionIdentifier)
}

// BuildNodePoolUpdateRequest builds a full node pool update from the current node pool.
// The API replaces the node pool, so every field is populated. Overrides replace only the fields that are set.
func BuildNodePoolUpdateRequest(current *kubernetes.KubernetesNodePool, overrides NodePoolUpdateOverrides) (kubernetes.UpdateKubernetesNodePool, error) {
	if current == nil {
		return kubernetes.UpdateKubernetesNodePool{}, fmt.Errorf("node pool is required")
	}

	machineType := current.MachineType.Name
	if overrides.MachineType != nil {
		machineType = *overrides.MachineType
	}
	if machineType == "" {
		return kubernetes.UpdateKubernetesNodePool{}, fmt.Errorf("node pool machine type is required")
	}

	versionIdentity := ""
	if current.KubernetesVersion != nil {
		versionIdentity = current.KubernetesVersion.Identity
	}
	if overrides.KubernetesVersionIdentity != nil {
		versionIdentity = *overrides.KubernetesVersionIdentity
	}
	if versionIdentity == "" {
		return kubernetes.UpdateKubernetesNodePool{}, fmt.Errorf("node pool kubernetes version is required")
	}

	replicas := current.Replicas
	if overrides.Replicas != nil {
		replicas = *overrides.Replicas
	}
	minReplicas := current.MinReplicas
	if overrides.MinReplicas != nil {
		minReplicas = *overrides.MinReplicas
	}
	maxReplicas := current.MaxReplicas
	if overrides.MaxReplicas != nil {
		maxReplicas = *overrides.MaxReplicas
	}

	enableAutoscaling := current.EnableAutoscaling
	if overrides.EnableAutoscaling != nil {
		enableAutoscaling = *overrides.EnableAutoscaling
	}
	enableAutoHealing := current.EnableAutoHealing
	if overrides.EnableAutoHealing != nil {
		enableAutoHealing = *overrides.EnableAutoHealing
	}

	upgradeStrategy := current.UpgradeStrategy
	if overrides.UpgradeStrategy != nil {
		upgradeStrategy = *overrides.UpgradeStrategy
	}

	nodeSettings := kubernetes.KubernetesNodeSettings{
		Labels:      cloneStringMap(current.NodeSettings.Labels),
		Annotations: cloneStringMap(current.NodeSettings.Annotations),
		Taints:      cloneTaints(current.NodeSettings.Taints),
	}
	if overrides.ReplaceTaints {
		nodeSettings.Taints = cloneTaints(overrides.Taints)
	}

	securityGroups, err := nodePoolSecurityGroupIdentities(current.SecurityGroups)
	if err != nil {
		return kubernetes.UpdateKubernetesNodePool{}, err
	}
	if overrides.ReplaceSecurityGroups {
		securityGroups = cloneStrings(overrides.SecurityGroupAttachments)
	}

	return kubernetes.UpdateKubernetesNodePool{
		Description:               current.Description,
		Labels:                    cloneStringMap(current.Labels),
		Annotations:               cloneStringMap(current.Annotations),
		MachineType:               machineType,
		Replicas:                  &replicas,
		MinReplicas:               &minReplicas,
		MaxReplicas:               &maxReplicas,
		KubernetesVersionIdentity: &versionIdentity,
		AvailabilityZone:          current.AvailabilityZone,
		UpgradeStrategy:           &upgradeStrategy,
		EnableAutoHealing:         &enableAutoHealing,
		EnableAutoscaling:         &enableAutoscaling,
		ManageNodeAllocatable:     current.ManageNodeAllocatable,
		NodeSettings:              &nodeSettings,
		SecurityGroupAttachments:  securityGroups,
	}, nil
}

func nodePoolSecurityGroupIdentities(groups []iaas.SecurityGroup) ([]string, error) {
	if len(groups) == 0 {
		return nil, nil
	}
	identities := make([]string, 0, len(groups))
	for _, group := range groups {
		if group.Identity == "" {
			return nil, fmt.Errorf("node pool security group %q has no identity", group.Name)
		}
		identities = append(identities, group.Identity)
	}
	return identities, nil
}

func cloneStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}
	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}

func cloneTaints(taints []kubernetes.NodeTaint) []kubernetes.NodeTaint {
	if taints == nil {
		return nil
	}
	cloned := make([]kubernetes.NodeTaint, len(taints))
	copy(cloned, taints)
	return cloned
}

// BuildNodePoolName builds the node pool name, optionally appending AZ suffix
func BuildNodePoolName(baseName string, az string, multipleAZs bool) string {
	if multipleAZs {
		return fmt.Sprintf("%s-%s", baseName, az)
	}
	return baseName
}

// BuildNodePoolCreateRequest builds a CreateKubernetesNodePool request
func BuildNodePoolCreateRequest(config NodePoolConfig, cluster *kubernetes.KubernetesCluster, az string, subnetIdentity string) (kubernetes.CreateKubernetesNodePool, error) {
	if cluster.ClusterVersion.Identity == "" {
		return kubernetes.CreateKubernetesNodePool{}, fmt.Errorf("cluster does not have a valid Kubernetes version")
	}

	kubernetesVersionIdentity := cluster.ClusterVersion.Identity
	replicas := CalculateReplicas(config.EnableAS, config.Replicas, config.MinNodes)
	nodePoolName := BuildNodePoolName(config.Name, az, false) // AZ suffix handled separately

	labels := ParseLabels(config.Labels)
	annotations := ParseAnnotations(config.Annotations)
	taints, err := ParseTaints(config.Taints)
	if err != nil {
		return kubernetes.CreateKubernetesNodePool{}, err
	}

	upgradeStrategy := ParseUpgradeStrategy(config.UpgradeStrat)

	req := kubernetes.CreateKubernetesNodePool{
		Name:                      nodePoolName,
		MachineType:               config.MachineType,
		KubernetesVersionIdentity: &kubernetesVersionIdentity,
		Replicas:                  replicas,
		EnableAutoscaling:         config.EnableAS,
		MinReplicas:               config.MinNodes,
		MaxReplicas:               config.MaxNodes,
		AvailabilityZone:          az,
		EnableAutoHealing:         config.EnableAH,
		UpgradeStrategy:           &upgradeStrategy,
		SubnetIdentity:            &subnetIdentity,
		NodeSettings: kubernetes.KubernetesNodeSettings{
			Labels:      labels,
			Annotations: annotations,
			Taints:      taints,
		},
	}

	if len(config.SecurityGroups) > 0 {
		req.SecurityGroupAttachments = config.SecurityGroups
	}

	return req, nil
}
