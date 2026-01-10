package kubernetes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var (
	updateName                  *string
	updateDescription           *string
	updateClusterVersion        *string
	updateDisablePublicEndpoint *bool
	updateKubeProxyMode         *string
	updateKubeProxyDeployment   *string
	updatePodSecurityStandards  *string
	updateAuditLogProfile       *string
	updateDefaultNetworkPolicy  *string
	updateUpgradeScheduleDay    *string
	updateUpgradeScheduleStart  *string
	updateWait                  bool
)

var updateCmd = &cobra.Command{
	Use:     "update <cluster>",
	Aliases: []string{"u", "edit", "modify"},
	Short:   "Update a Kubernetes cluster",
	Long: `Update an existing Kubernetes cluster.

Only fields specified with flags will be updated. All other fields will remain unchanged.

Examples:
  # Update cluster description
  tcloud kubernetes update my-cluster --description "Production cluster"

  # Update Kubernetes version
  tcloud kubernetes update my-cluster --cluster-version 1.28.0

  # Update maintenance window
  tcloud kubernetes update my-cluster --maintenance-day monday --maintenance-start 02:00

  # Update multiple fields
  tcloud kubernetes update my-cluster --description "Updated description" --disable-public-endpoint

Note: Use 'tcloud kubernetes label' and 'tcloud kubernetes annotate' commands
to manage labels and annotations separately.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		clusterIdentifier := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Resolve cluster
		clusters, err := client.Kubernetes().ListKubernetesClusters(ctx, &kubernetes.ListKubernetesClustersRequest{})
		if err != nil {
			return fmt.Errorf("failed to list clusters: %w", err)
		}

		var cluster *kubernetes.KubernetesCluster
		for _, c := range clusters {
			if strings.EqualFold(c.Identity, clusterIdentifier) ||
				strings.EqualFold(c.Name, clusterIdentifier) ||
				strings.EqualFold(c.Slug, clusterIdentifier) {
				cluster = &c
				break
			}
		}

		if cluster == nil {
			return fmt.Errorf("cluster not found: %s", clusterIdentifier)
		}

		// Get fresh cluster data
		cluster, err = client.Kubernetes().GetKubernetesCluster(ctx, cluster.Identity)
		if err != nil {
			return fmt.Errorf("failed to get cluster: %w", err)
		}

		// Build update request
		updateReq := kubernetes.UpdateKubernetesCluster{}

		// Update name
		if updateName != nil && *updateName != "" {
			updateReq.Name = updateName
		}

		// Update description
		if updateDescription != nil {
			updateReq.Description = updateDescription
		}

		// Update Kubernetes version
		if updateClusterVersion != nil && *updateClusterVersion != "" {
			versions, err := client.Kubernetes().ListKubernetesVersions(ctx)
			if err != nil {
				return fmt.Errorf("failed to list kubernetes versions: %w", err)
			}

			var versionIdentity string
			for _, v := range versions {
				if !v.Enabled {
					continue
				}
				if strings.EqualFold(v.Identity, *updateClusterVersion) ||
					strings.EqualFold(v.Slug, *updateClusterVersion) ||
					strings.EqualFold(v.Name, *updateClusterVersion) {
					versionIdentity = v.Identity
					break
				}
			}

			if versionIdentity == "" {
				return fmt.Errorf("kubernetes version not found or not enabled: %s", *updateClusterVersion)
			}

			updateReq.KubernetesVersionIdentity = &versionIdentity
		}

		// Update public endpoint
		if updateDisablePublicEndpoint != nil {
			updateReq.DisablePublicEndpoint = updateDisablePublicEndpoint
		}

		// Update kube proxy mode
		if updateKubeProxyMode != nil && *updateKubeProxyMode != "" {
			mode := kubernetes.KubernetesClusterKubeProxyMode(*updateKubeProxyMode)
			updateReq.KubeProxyMode = &mode
		}

		// Update kube proxy deployment
		if updateKubeProxyDeployment != nil && *updateKubeProxyDeployment != "" {
			deployment := kubernetes.KubeProxyDeployment(*updateKubeProxyDeployment)
			updateReq.KubeProxyDeployment = &deployment
		}

		// Update pod security standards
		if updatePodSecurityStandards != nil && *updatePodSecurityStandards != "" {
			profile := kubernetes.KubernetesClusterPodSecurityStandards(*updatePodSecurityStandards)
			updateReq.PodSecurityStandardsProfile = &profile
		}

		// Update audit log profile
		if updateAuditLogProfile != nil && *updateAuditLogProfile != "" {
			profile := kubernetes.KubernetesClusterAuditLoggingProfile(*updateAuditLogProfile)
			updateReq.AuditLogProfile = &profile
		}

		// Update default network policy
		if updateDefaultNetworkPolicy != nil && *updateDefaultNetworkPolicy != "" {
			policy := kubernetes.KubernetesDefaultNetworkPolicies(*updateDefaultNetworkPolicy)
			updateReq.DefaultNetworkPolicy = &policy
		}

		// Update upgrade schedule
		if updateUpgradeScheduleDay != nil || updateUpgradeScheduleStart != nil {
			// Parse and set maintenance day if provided
			if updateUpgradeScheduleDay != nil && *updateUpgradeScheduleDay != "" {
				day, err := parseDayOfWeek(*updateUpgradeScheduleDay)
				if err != nil {
					return fmt.Errorf("invalid maintenance day: %w", err)
				}
				updateReq.MaintenanceDay = &day
			}

			// Parse and set maintenance start time if provided
			if updateUpgradeScheduleStart != nil && *updateUpgradeScheduleStart != "" {
				startTime, err := parseTimeOfDay(*updateUpgradeScheduleStart)
				if err != nil {
					return fmt.Errorf("invalid maintenance start time: %w", err)
				}
				updateReq.MaintenanceStartAt = &startTime
			}

			// Set auto-upgrade policy if schedule is provided
			updateReq.AutoUpgradePolicy = kubernetes.KubernetesClusterAutoUpgradePolicyLatestStable
		}

		// Check if there are any updates
		if updateReq.Name == nil &&
			updateReq.Description == nil &&
			updateReq.KubernetesVersionIdentity == nil &&
			updateReq.DisablePublicEndpoint == nil &&
			updateReq.KubeProxyMode == nil &&
			updateReq.KubeProxyDeployment == nil &&
			updateReq.PodSecurityStandardsProfile == nil &&
			updateReq.AuditLogProfile == nil &&
			updateReq.DefaultNetworkPolicy == nil &&
			updateReq.MaintenanceDay == nil &&
			updateReq.MaintenanceStartAt == nil {
			return fmt.Errorf("no updates specified")
		}

		updatedCluster, err := client.Kubernetes().UpdateKubernetesCluster(ctx, cluster.Identity, updateReq)
		if err != nil {
			return fmt.Errorf("failed to update cluster: %w", err)
		}

		if updateWait {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
			defer cancel()

			_, err = client.Kubernetes().WaitUntilKubernetesClusterReady(ctxWithTimeout, updatedCluster.Identity)
			if err != nil {
				return fmt.Errorf("failed to wait for cluster to be ready: %w", err)
			}
			// Get fresh cluster data
			updatedCluster, err = client.Kubernetes().GetKubernetesCluster(ctx, updatedCluster.Identity)
			if err != nil {
				return fmt.Errorf("failed to get cluster: %w", err)
			}
		}

		// Output in table format like IaaS commands
		vpcName := ""
		if updatedCluster.VPC != nil {
			vpcName = updatedCluster.VPC.Name
		}
		body := [][]string{
			{
				updatedCluster.Identity,
				updatedCluster.Name,
				vpcName,
				updatedCluster.ClusterVersion.Name,
				string(updatedCluster.ClusterType),
				updatedCluster.Status,
				formattime.FormatTime(updatedCluster.CreatedAt.Local(), false),
			},
		}
		table.Print([]string{"ID", "Name", "Vpc", "Version", "Type", "Status", "Age"}, body)

		return nil
	},
}

func init() {
	KubernetesCmd.AddCommand(updateCmd)

	updateName = updateCmd.Flags().String("name", "", "Name of the cluster")
	updateDescription = updateCmd.Flags().String("description", "", "Description of the cluster")
	updateClusterVersion = updateCmd.Flags().String("cluster-version", "", "Kubernetes version (name, slug, or identity)")
	updateDisablePublicEndpoint = updateCmd.Flags().Bool("disable-public-endpoint", false, "Disable public API server endpoint")
	updateKubeProxyMode = updateCmd.Flags().String("kube-proxy-mode", "", "Kube proxy mode: iptables or ipvs")
	updateKubeProxyDeployment = updateCmd.Flags().String("kube-proxy-deployment", "", "Kube proxy deployment: disabled, managed, or custom")
	updatePodSecurityStandards = updateCmd.Flags().String("pod-security-standards", "", "Pod security standards profile: baseline, restricted, or privileged")
	updateAuditLogProfile = updateCmd.Flags().String("audit-log-profile", "", "Audit log profile: none, basic, or metadata")
	updateDefaultNetworkPolicy = updateCmd.Flags().String("default-network-policy", "", "Default network policy: allow-all, deny-all, or none")
	updateUpgradeScheduleDay = updateCmd.Flags().String("maintenance-day", "", "Maintenance day: 0-6, Sunday-Saturday, or day name")
	updateUpgradeScheduleStart = updateCmd.Flags().String("maintenance-start", "", "Maintenance start time: HH:MM format (e.g., '02:00' or '14:30')")
	updateCmd.Flags().BoolVar(&updateWait, "wait", false, "Wait for the cluster update to complete")

	updateCmd.RegisterFlagCompletionFunc("cluster-version", completion.CompleteKubernetesVersion)
	updateCmd.RegisterFlagCompletionFunc("kube-proxy-mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"iptables", "ipvs"}, cobra.ShellCompDirectiveNoFileComp
	})
	updateCmd.RegisterFlagCompletionFunc("kube-proxy-deployment", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"disabled", "managed", "custom"}, cobra.ShellCompDirectiveNoFileComp
	})
	updateCmd.RegisterFlagCompletionFunc("pod-security-standards", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"baseline", "restricted", "privileged"}, cobra.ShellCompDirectiveNoFileComp
	})
	updateCmd.RegisterFlagCompletionFunc("audit-log-profile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"none", "basic", "metadata"}, cobra.ShellCompDirectiveNoFileComp
	})
	updateCmd.RegisterFlagCompletionFunc("default-network-policy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"allow-all", "deny-all", "none"}, cobra.ShellCompDirectiveNoFileComp
	})
	updateCmd.RegisterFlagCompletionFunc("maintenance-day", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"0", "sunday", "1", "monday", "2", "tuesday", "3", "wednesday", "4", "thursday", "5", "friday", "6", "saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
}
