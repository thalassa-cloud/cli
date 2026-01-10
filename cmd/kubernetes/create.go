package kubernetes

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/kubernetes/nodepools"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/fzf"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/thalassa"
)

var (
	createName                  string
	createDescription           string
	createClusterType           string
	createRegion                string
	createSubnet                string
	createClusterVersion        string
	createNetworkingCNI         string
	createNetworkingServiceCIDR string
	createNetworkingPodCIDR     string
	createWait                  bool
	createDisablePublicEndpoint bool
	createLabels                []string
	createAnnotations           []string
	// Kube proxy flags
	createKubeProxyMode       string
	createKubeProxyDeployment string
	// Security flags
	createPodSecurityStandards string
	createAuditLogProfile      string
	createDefaultNetworkPolicy string
	// Upgrade schedule flags
	createUpgradeScheduleDay   string
	createUpgradeScheduleStart string
	// Node pool flags
	createNodePoolName           string
	createNodePoolMachineType    string
	createNodePoolReplicas       int
	createNodePoolEnableAS       bool
	createNodePoolMinNodes       int
	createNodePoolMaxNodes       int
	createNodePoolSubnet         string
	createNodePoolAZs            []string
	createNodePoolEnableAH       bool
	createNodePoolUpgradeStrat   string
	createNodePoolLabels         []string
	createNodePoolAnnotations    []string
	createNodePoolTaints         []string
	createNodePoolSecurityGroups []string
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c", "new"},
	Short:   "Create a Kubernetes cluster",
	Long: `Create a Kubernetes cluster in the Thalassa Cloud Platform.

This command creates a new Kubernetes cluster with sensible defaults.

Examples:
  # Create a managed cluster with minimal configuration
  tcloud kubernetes create my-cluster --subnet subnet-123

  # Create a cluster with custom networking
  tcloud kubernetes create my-cluster --subnet subnet-123 --pod-cidr 10.0.0.0/16 --service-cidr 172.16.0.0/18

  # Create a cluster and wait for it to be ready
  tcloud kubernetes create my-cluster --subnet subnet-123 --wait

  # Create a cluster with a node pool
  tcloud kubernetes create my-cluster --subnet subnet-123 --machine-type pgp-medium --num-nodes 3 --availability-zone nl-01a

  # Create a cluster with node pools in multiple availability zones
  tcloud kubernetes create my-cluster --subnet subnet-123 --machine-type pgp-medium --num-nodes 3 --availability-zone nl-01a --availability-zone nl-01b

  # Create a cluster with autoscaling node pool (auto-selects AZ if not provided)
  tcloud kubernetes create my-cluster --subnet subnet-123 --machine-type pgp-medium --enable-autoscaling --min-nodes 1 --max-nodes 5

  # Create a cluster with scheduled upgrades (maintenance window)
  tcloud kubernetes create my-cluster --subnet subnet-123 --maintenance-day monday --maintenance-start 02:00`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		clusterName := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Determine cluster type
		clusterType := kubernetes.KubernetesClusterType("managed")
		if createClusterType != "" {
			clusterType = kubernetes.KubernetesClusterType(createClusterType)
		}

		// Validate required fields based on cluster type
		if clusterType == kubernetes.KubernetesClusterType("hosted-control-plane") {
			if createRegion == "" {
				// Try interactive selection
				if fzf.IsInteractiveMode(os.Stdout) {
					command := fmt.Sprintf("%s iaas regions list --no-header", os.Args[0])
					selected, err := fzf.InteractiveChoice(command)
					if err != nil {
						return fmt.Errorf("region is required for hosted-control-plane clusters: %w", err)
					}
					createRegion = selected
				} else {
					return fmt.Errorf("region is required for hosted-control-plane clusters")
				}
			}
		} else {
			if createSubnet == "" {
				// Try interactive selection
				if fzf.IsInteractiveMode(os.Stdout) {
					command := fmt.Sprintf("%s networking subnets list --no-header", os.Args[0])
					selected, err := fzf.InteractiveChoice(command)
					if err != nil {
						return fmt.Errorf("subnet is required for managed clusters: %w", err)
					}
					createSubnet = selected
				} else {
					return fmt.Errorf("subnet is required for managed clusters")
				}
			}
		}

		// Resolve region if provided
		var regionIdentity string
		if createRegion != "" {
			regions, err := client.IaaS().ListRegions(ctx, &iaas.ListRegionsRequest{})
			if err != nil {
				return fmt.Errorf("failed to list regions: %w", err)
			}
			for _, r := range regions {
				if r.Identity == createRegion || r.Slug == createRegion || r.Name == createRegion {
					regionIdentity = r.Identity
					break
				}
			}
			if regionIdentity == "" {
				return fmt.Errorf("region not found: %s", createRegion)
			}
		}

		// Resolve subnet if provided
		var subnetIdentity string
		var subnet *iaas.Subnet
		if createSubnet != "" {
			subnets, err := client.IaaS().ListSubnets(ctx, &iaas.ListSubnetsRequest{})
			if err != nil {
				return fmt.Errorf("failed to list subnets: %w", err)
			}
			for _, s := range subnets {
				if s.Identity == createSubnet || s.Slug == createSubnet || s.Name == createSubnet {
					subnetIdentity = s.Identity
					subnet = &s
					break
				}
			}
			if subnetIdentity == "" {
				return fmt.Errorf("subnet not found: %s", createSubnet)
			}
		}

		// auto detect region based on subnet, if region is not provided
		if regionIdentity == "" && subnet != nil && subnet.Vpc != nil {
			if subnet.Vpc.CloudRegion != nil {
				regionIdentity = subnet.Vpc.CloudRegion.Identity
			} else {
				vpc, err := client.IaaS().GetVpc(ctx, subnet.Vpc.Identity)
				if err != nil {
					return fmt.Errorf("failed to get vpc: %w", err)
				}
				regionIdentity = vpc.CloudRegion.Identity
			}
		}

		// Resolve cluster version
		var versionIdentity string
		if createClusterVersion != "" {
			versions, err := client.Kubernetes().ListKubernetesVersions(ctx)
			if err != nil {
				return fmt.Errorf("failed to list kubernetes versions: %w", err)
			}
			for _, v := range versions {
				if !v.Enabled {
					continue
				}
				if v.Identity == createClusterVersion || v.Slug == createClusterVersion || v.Name == createClusterVersion {
					versionIdentity = v.Identity
					break
				}
			}
			if versionIdentity == "" {
				return fmt.Errorf("kubernetes version not found or not enabled: %s", createClusterVersion)
			}
		} else {
			// Get latest enabled version (prefer latest by sorting or just take first enabled)
			versions, err := client.Kubernetes().ListKubernetesVersions(ctx)
			if err != nil {
				return fmt.Errorf("failed to list kubernetes versions: %w", err)
			}
			// Find the first enabled version (versions are typically sorted with latest first)
			for _, v := range versions {
				if v.Enabled {
					versionIdentity = v.Identity
					break
				}
			}
			if versionIdentity == "" {
				return fmt.Errorf("no enabled kubernetes version found")
			}
		}

		// Parse labels
		labels := make(map[string]string)
		for _, label := range createLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Parse annotations
		annotations := make(map[string]string)
		for _, annotation := range createAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Set defaults for networking
		networkingCNI := "cilium"
		if createNetworkingCNI != "" {
			networkingCNI = createNetworkingCNI
		}

		serviceCIDR := "172.16.0.0/18"
		if createNetworkingServiceCIDR != "" {
			serviceCIDR = createNetworkingServiceCIDR
		}

		podCIDR := "192.168.0.0/16"
		if createNetworkingPodCIDR != "" {
			podCIDR = createNetworkingPodCIDR
		}

		// Build create request
		createReq := kubernetes.CreateKubernetesCluster{
			Name:                      clusterName,
			Description:               createDescription,
			Labels:                    labels,
			Annotations:               annotations,
			ClusterType:               clusterType,
			KubernetesVersionIdentity: versionIdentity,
			Subnet:                    subnetIdentity,
			RegionIdentity:            regionIdentity,
			DisablePublicEndpoint:     createDisablePublicEndpoint,
			Networking: kubernetes.KubernetesClusterNetworking{
				CNI:         networkingCNI,
				ServiceCIDR: serviceCIDR,
				PodCIDR:     podCIDR,
			},
		}

		// Set kube proxy configuration
		if createKubeProxyMode != "" {
			mode := kubernetes.KubernetesClusterKubeProxyMode(createKubeProxyMode)
			createReq.KubeProxyMode = &mode
		} else {
			// Default to iptables
			mode := kubernetes.KubernetesClusterKubeProxyModeIptables
			createReq.KubeProxyMode = &mode
		}

		if createKubeProxyDeployment != "" {
			deployment := kubernetes.KubeProxyDeployment(createKubeProxyDeployment)
			createReq.KubeProxyDeployment = &deployment
		} else {
			// Default to disabled
			deployment := kubernetes.KubeProxyDeploymentDisabled
			createReq.KubeProxyDeployment = &deployment
		}

		// Set security configuration
		if createPodSecurityStandards != "" {
			createReq.PodSecurityStandardsProfile = kubernetes.KubernetesClusterPodSecurityStandards(createPodSecurityStandards)
		} else {
			// Default to baseline
			createReq.PodSecurityStandardsProfile = kubernetes.KubernetesClusterPodSecurityStandardBaseline
		}

		if createAuditLogProfile != "" {
			createReq.AuditLogProfile = kubernetes.KubernetesClusterAuditLoggingProfile(createAuditLogProfile)
		} else {
			// Default to none
			createReq.AuditLogProfile = kubernetes.KubernetesClusterAuditLoggingProfileNone
		}

		if createDefaultNetworkPolicy != "" {
			createReq.DefaultNetworkPolicy = kubernetes.KubernetesDefaultNetworkPolicies(createDefaultNetworkPolicy)
		} else {
			// Default to allow-all
			createReq.DefaultNetworkPolicy = kubernetes.KubernetesDefaultNetworkPolicyAllowAll
		}

		// Set upgrade policy and schedule
		if createUpgradeScheduleDay != "" || createUpgradeScheduleStart != "" {
			// Upgrade schedule is provided, set auto-upgrade policy
			createReq.AutoUpgradePolicy = kubernetes.KubernetesClusterAutoUpgradePolicyLatestStable

			// Parse and set maintenance day if provided
			if createUpgradeScheduleDay != "" {
				day, err := parseDayOfWeek(createUpgradeScheduleDay)
				if err != nil {
					return fmt.Errorf("invalid maintenance day: %w", err)
				}
				createReq.MaintenanceDay = &day
			}

			// Parse and set maintenance start time if provided
			if createUpgradeScheduleStart != "" {
				startTime, err := parseTimeOfDay(createUpgradeScheduleStart)
				if err != nil {
					return fmt.Errorf("invalid maintenance start time: %w", err)
				}
				createReq.MaintenanceStartAt = &startTime
			}
		} else {
			// No schedule provided, default to "none"
			createReq.AutoUpgradePolicy = kubernetes.KubernetesClusterAutoUpgradePolicy("none")
		}

		cluster, err := client.Kubernetes().CreateKubernetesCluster(ctx, createReq)
		if err != nil {
			return fmt.Errorf("failed to create cluster: %w", err)
		}

		if createWait {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
			defer cancel()

			_, err = client.Kubernetes().WaitUntilKubernetesClusterReady(ctxWithTimeout, cluster.Identity)
			if err != nil {
				return fmt.Errorf("failed to wait for cluster to be ready: %w", err)
			}
			// Get fresh cluster data after wait
			cluster, err = client.Kubernetes().GetKubernetesCluster(ctx, cluster.Identity)
			if err != nil {
				return fmt.Errorf("failed to get cluster: %w", err)
			}
		}

		// Output in table format like IaaS commands
		vpcName := ""
		if cluster.VPC != nil {
			vpcName = cluster.VPC.Name
		}
		body := [][]string{
			{
				cluster.Identity,
				cluster.Name,
				vpcName,
				cluster.ClusterVersion.Name,
				string(cluster.ClusterType),
				cluster.Status,
				formattime.FormatTime(cluster.CreatedAt.Local(), false),
			},
		}
		table.Print([]string{"ID", "Name", "Vpc", "Version", "Type", "Status", "Age"}, body)

		// Create node pool if requested
		if createNodePoolMachineType != "" {
			// Ensure cluster is ready before creating node pool
			if !createWait {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
				defer cancel()

				_, err = client.Kubernetes().WaitUntilKubernetesClusterReady(ctxWithTimeout, cluster.Identity)
				if err != nil {
					return fmt.Errorf("failed to wait for cluster to be ready: %w", err)
				}
				// Get fresh cluster data
				cluster, err = client.Kubernetes().GetKubernetesCluster(ctx, cluster.Identity)
				if err != nil {
					return fmt.Errorf("failed to get cluster: %w", err)
				}
			}

			if err := createNodePool(ctx, client, cluster, clusterName); err != nil {
				return fmt.Errorf("failed to create node pool: %w", err)
			}
		}

		return nil
	},
}

func createNodePool(ctx context.Context, client thalassa.Client, cluster *kubernetes.KubernetesCluster, clusterName string) error {
	// Resolve machine type
	resolvedMachineType, err := nodepools.ResolveMachineType(ctx, client, createNodePoolMachineType)
	if err != nil {
		return err
	}

	// Validate autoscaling configuration
	if err := nodepools.ValidateAutoscalingConfig(createNodePoolEnableAS, createNodePoolMinNodes, createNodePoolMaxNodes, createNodePoolReplicas); err != nil {
		return err
	}

	// Determine node pool subnet
	nodePoolSubnetIdentity, err := nodepools.ResolveNodePoolSubnet(ctx, client, createNodePoolSubnet, cluster)
	if err != nil {
		return err
	}

	// Determine availability zones
	availabilityZones, err := nodepools.DetermineAvailabilityZones(ctx, client, createNodePoolAZs, cluster)
	if err != nil {
		return err
	}

	// Base node pool name
	baseNodePoolName := createNodePoolName
	if baseNodePoolName == "" {
		baseNodePoolName = "default-pool"
	}

	// Build node pool config
	config := nodepools.NodePoolConfig{
		Name:           baseNodePoolName,
		MachineType:    resolvedMachineType,
		Replicas:       createNodePoolReplicas,
		EnableAS:       createNodePoolEnableAS,
		MinNodes:       createNodePoolMinNodes,
		MaxNodes:       createNodePoolMaxNodes,
		EnableAH:       createNodePoolEnableAH,
		UpgradeStrat:   createNodePoolUpgradeStrat,
		Labels:         createNodePoolLabels,
		Annotations:    createNodePoolAnnotations,
		Taints:         createNodePoolTaints,
		SecurityGroups: createNodePoolSecurityGroups,
	}

	// Create node pool in each availability zone
	for _, az := range availabilityZones {
		// Determine node pool name
		nodePoolName := nodepools.BuildNodePoolName(baseNodePoolName, az, len(availabilityZones) > 1)

		// Build create request
		createNodePoolReq, err := nodepools.BuildNodePoolCreateRequest(config, cluster, az, nodePoolSubnetIdentity)
		if err != nil {
			return err
		}
		createNodePoolReq.Name = nodePoolName

		nodePool, err := client.Kubernetes().CreateKubernetesNodePool(ctx, cluster.Identity, createNodePoolReq)
		if err != nil {
			return fmt.Errorf("failed to create node pool in availability zone %s: %w", az, err)
		}

		if createWait {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
			defer cancel()

			_, err = client.Kubernetes().WaitUntilKubernetesNodePoolReady(ctxWithTimeout, cluster.Identity, nodePool.Identity)
			if err != nil {
				return fmt.Errorf("failed to wait for node pool to be ready: %w", err)
			}
			// Get fresh node pool data
			nodePool, err = client.Kubernetes().GetKubernetesNodePool(ctx, cluster.Identity, nodePool.Identity)
			if err != nil {
				return fmt.Errorf("failed to get node pool: %w", err)
			}
		}

		// Output in table format like IaaS commands
		replicas := formatReplicas(nodePool)
		body := [][]string{
			{
				nodePool.Identity,
				nodePool.Name,
				cluster.Name,
				replicas,
				nodePool.MachineType.Name,
				string(nodePool.Status),
				formattime.FormatTime(nodePool.CreatedAt.Local(), false),
			},
		}
		table.Print([]string{"ID", "Name", "Cluster", "Replicas", "Type", "Status", "Age"}, body)
	}

	return nil
}

// formatReplicas formats the replicas string based on autoscaling settings
func formatReplicas(np *kubernetes.KubernetesNodePool) string {
	if np.EnableAutoscaling {
		return fmt.Sprintf("%d-%d", np.MinReplicas, np.MaxReplicas)
	}
	return fmt.Sprintf("%d", np.Replicas)
}

func init() {
	KubernetesCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createName, "name", "", "Name of the cluster (deprecated: use positional argument)")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description of the cluster")
	createCmd.Flags().StringVar(&createClusterType, "cluster-type", "managed", "Cluster type: managed or hosted-control-plane")
	createCmd.Flags().StringVar(&createRegion, "region", "", "Region for hosted-control-plane clusters")
	createCmd.Flags().StringVar(&createSubnet, "subnet", "", "Subnet for managed clusters")
	createCmd.Flags().StringVar(&createClusterVersion, "cluster-version", "", "Kubernetes version (name, slug, or identity). Defaults to latest stable")
	createCmd.Flags().StringVar(&createNetworkingCNI, "cni", "cilium", "CNI plugin: cilium or custom")
	createCmd.Flags().StringVar(&createNetworkingServiceCIDR, "service-cidr", "172.16.0.0/18", "Service CIDR")
	createCmd.Flags().StringVar(&createNetworkingPodCIDR, "pod-cidr", "192.168.0.0/16", "Pod CIDR")
	createCmd.Flags().BoolVar(&createWait, "wait", false, "Wait for the cluster to be ready before returning")
	createCmd.Flags().BoolVar(&createDisablePublicEndpoint, "disable-public-endpoint", false, "Disable public API server endpoint")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	createCmd.Flags().StringVar(&createKubeProxyMode, "kube-proxy-mode", "", "Kube proxy mode: iptables or ipvs (default: iptables)")
	createCmd.Flags().StringVar(&createKubeProxyDeployment, "kube-proxy-deployment", "", "Kube proxy deployment: disabled, managed, or custom (default: disabled)")
	createCmd.Flags().StringVar(&createPodSecurityStandards, "pod-security-standards", "", "Pod security standards profile: baseline, restricted, or privileged (default: baseline)")
	createCmd.Flags().StringVar(&createAuditLogProfile, "audit-log-profile", "", "Audit log profile: none, basic, or metadata (default: none)")
	createCmd.Flags().StringVar(&createDefaultNetworkPolicy, "default-network-policy", "", "Default network policy: allow-all, deny-all, or none (default: allow-all)")

	// Node pool flags
	createCmd.Flags().StringVar(&createNodePoolName, "node-pool-name", "worker", "Name of the node pool")
	createCmd.Flags().StringVar(&createNodePoolMachineType, "machine-type", "", "Machine type for the node pool (required to create a node pool)")
	createCmd.Flags().IntVar(&createNodePoolReplicas, "num-nodes", 1, "Number of nodes in the node pool (ignored if --enable-autoscaling is set)")
	createCmd.Flags().BoolVar(&createNodePoolEnableAS, "enable-autoscaling", false, "Enable autoscaling for the node pool")
	createCmd.Flags().IntVar(&createNodePoolMinNodes, "min-nodes", 1, "Minimum number of nodes (required when autoscaling is enabled)")
	createCmd.Flags().IntVar(&createNodePoolMaxNodes, "max-nodes", 3, "Maximum number of nodes (required when autoscaling is enabled)")
	createCmd.Flags().StringVar(&createNodePoolSubnet, "node-pool-subnet", "", "Subnet for the node pool (defaults to cluster subnet)")
	createCmd.Flags().StringSliceVar(&createNodePoolAZs, "availability-zone", []string{}, "Availability zone for the node pool (can be specified multiple times to create node pools in multiple AZs). If not specified, a random AZ from the cluster's region will be selected.")
	createCmd.Flags().BoolVar(&createNodePoolEnableAH, "enable-autohealing", false, "Enable autohealing for the node pool")
	createCmd.Flags().StringVar(&createNodePoolUpgradeStrat, "upgrade-strategy", "auto", "Upgrade strategy: manual, auto, always, on-delete, inplace, or never")
	createCmd.Flags().StringSliceVar(&createNodePoolLabels, "node-labels", []string{}, "Node labels in key=value format (applied to Kubernetes nodes)")
	createCmd.Flags().StringSliceVar(&createNodePoolAnnotations, "node-annotations", []string{}, "Node annotations in key=value format (applied to Kubernetes nodes)")
	createCmd.Flags().StringSliceVar(&createNodePoolTaints, "node-taints", []string{}, "Node taints in key=value:effect or key:effect format (e.g., 'dedicated=gpu:NoSchedule')")
	createCmd.Flags().StringSliceVar(&createNodePoolSecurityGroups, "security-groups", []string{}, "Security group identities to attach to node pool machines")

	// Register completions
	createCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegionEnhanced)
	createCmd.RegisterFlagCompletionFunc("subnet", completion.CompleteSubnetEnhanced)
	createCmd.RegisterFlagCompletionFunc("cluster-version", completion.CompleteKubernetesVersion)
	createCmd.RegisterFlagCompletionFunc("cluster-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"managed", "hosted-control-plane"}, cobra.ShellCompDirectiveNoFileComp
	})
	createCmd.RegisterFlagCompletionFunc("cni", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"cilium", "custom"}, cobra.ShellCompDirectiveNoFileComp
	})
	createCmd.RegisterFlagCompletionFunc("machine-type", completion.CompleteMachineType)
	createCmd.RegisterFlagCompletionFunc("upgrade-strategy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"manual", "auto", "always", "on-delete", "inplace", "never"}, cobra.ShellCompDirectiveNoFileComp
	})
	createCmd.RegisterFlagCompletionFunc("maintenance-day", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"0", "sunday", "1", "monday", "2", "tuesday", "3", "wednesday", "4", "thursday", "5", "friday", "6", "saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	createCmd.RegisterFlagCompletionFunc("kube-proxy-mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"iptables", "ipvs"}, cobra.ShellCompDirectiveNoFileComp
	})
	createCmd.RegisterFlagCompletionFunc("kube-proxy-deployment", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"disabled", "managed", "custom"}, cobra.ShellCompDirectiveNoFileComp
	})
	createCmd.RegisterFlagCompletionFunc("pod-security-standards", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"baseline", "restricted", "privileged"}, cobra.ShellCompDirectiveNoFileComp
	})
	createCmd.RegisterFlagCompletionFunc("audit-log-profile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"none", "basic", "metadata"}, cobra.ShellCompDirectiveNoFileComp
	})
	createCmd.RegisterFlagCompletionFunc("default-network-policy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"allow-all", "deny-all", "none"}, cobra.ShellCompDirectiveNoFileComp
	})
}

var (
	dayMap = map[string]uint{
		"sunday":    0,
		"sun":       0,
		"0":         0,
		"monday":    1,
		"mon":       1,
		"1":         1,
		"tuesday":   2,
		"tue":       2,
		"2":         2,
		"wednesday": 3,
		"wed":       3,
		"3":         3,
		"thursday":  4,
		"thu":       4,
		"4":         4,
		"friday":    5,
		"fri":       5,
		"5":         5,
		"saturday":  6,
		"sat":       6,
		"6":         6,
	}
)

// parseDayOfWeek converts a day of week string to a uint (0=Sunday, 1=Monday, ..., 6=Saturday)
func parseDayOfWeek(dayStr string) (uint, error) {
	dayStr = strings.ToLower(strings.TrimSpace(dayStr))
	if day, ok := dayMap[dayStr]; ok {
		return day, nil
	}
	return 0, fmt.Errorf("invalid day of week: %s (expected 0-6, Sunday-Saturday, or day name)", dayStr)
}

// parseTimeOfDay converts a time string (HH:MM) to minutes since midnight (uint)
func parseTimeOfDay(timeStr string) (uint, error) {
	timeStr = strings.TrimSpace(timeStr)

	// Parse HH:MM format
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format: %s (expected HH:MM, e.g., '02:00' or '14:30')", timeStr)
	}

	var hour, minute uint
	_, err := fmt.Sscanf(parts[0], "%d", &hour)
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %s", parts[0])
	}

	_, err = fmt.Sscanf(parts[1], "%d", &minute)
	if err != nil {
		return 0, fmt.Errorf("invalid minute: %s", parts[1])
	}

	if hour > 23 {
		return 0, fmt.Errorf("hour must be between 0 and 23, got: %d", hour)
	}
	if minute > 59 {
		return 0, fmt.Errorf("minute must be between 0 and 59, got: %d", minute)
	}

	// Convert to minutes since midnight
	return hour*60 + minute, nil
}
