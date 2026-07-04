package kubernetes

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/fzf"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"

	"k8s.io/utils/ptr"

	"github.com/blang/semver/v4"
)

var (
	upgradeClusterToVersion string
	upgradeClusterDryRun    bool
	upgradeAllClusters      bool
)

// Group candidates by major.minor to identify patch versions
type VersionGroup struct {
	MajorMinor string
	Versions   []kubernetes.KubernetesVersion
}

// selectUpgradeVersion selects the appropriate upgrade version based on the current version
// following Kubernetes versioning policy (max 1 minor version upgrade at a time)
func selectUpgradeVersion(currentVersion semver.Version, versions []kubernetes.KubernetesVersion, requestedVersion string) (*kubernetes.KubernetesVersion, error) {
	if requestedVersion != "" {
		return findAndValidateRequestedVersion(currentVersion, versions, requestedVersion)
	}

	newerVersions, upgradeCandidates := collectUpgradeCandidates(currentVersion, versions)
	if len(newerVersions) == 0 {
		return nil, fmt.Errorf("already at the latest available version (%s)", currentVersion)
	}
	if len(upgradeCandidates) == 0 {
		return nil, upgradePolicyBlockedError(currentVersion, newerVersions)
	}

	return selectBestUpgradeCandidate(currentVersion, upgradeCandidates)
}

var KubernetesUpgradeCmd = &cobra.Command{
	Use:     "upgrade <cluster>",
	Short:   "Upgrade a Kubernetes cluster",
	Aliases: []string{"u"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			fmt.Println("Error getting client:", err)
			return
		}
		versions, err := client.Kubernetes().ListKubernetesVersions(cmd.Context())
		if err != nil {
			fmt.Println("Error getting versions:", err)
			return
		}

		clusters := []kubernetes.KubernetesCluster{}

		if upgradeAllClusters {
			fmt.Println("Upgrading all clusters")
			clusters, err = client.Kubernetes().ListKubernetesClusters(cmd.Context(), &kubernetes.ListKubernetesClustersRequest{})
			if err != nil {
				fmt.Println("Error getting clusters:", err)
				return
			}
		} else {
			// get the cluster
			clusterIdentity, err := getSelectedCluster(args)
			if err != nil {
				fmt.Println("Error getting cluster:", err)
				return
			}

			cluster, err := client.Kubernetes().GetKubernetesCluster(cmd.Context(), clusterIdentity)
			if err != nil {
				fmt.Println("Error getting cluster:", err)
				return
			}
			clusters = append(clusters, *cluster)
		}

		for _, cluster := range clusters {
			if cluster.Status == "error" || cluster.Status == "deleting" {
				fmt.Println("Cluster is in error or deleting state...")
				continue
			}

			currentVersion, err := semver.Parse(cluster.ClusterVersion.KubernetesVersion)
			if err != nil {
				if !upgradeAllClusters {
					fmt.Println("Error parsing current cluster version:", err)
				}
				continue
			}

			upgradeToVersion, err := selectUpgradeVersion(currentVersion, versions, upgradeClusterToVersion)
			if err != nil {
				// Check if it's the "already at latest version" error
				if upgradeClusterToVersion == "" && strings.Contains(err.Error(), "already at the latest available version") {
					fmt.Println("Cluster", cluster.Name, "is already at the latest available version:", cluster.ClusterVersion.KubernetesVersion)
					continue
				}
				if !upgradeAllClusters {
					fmt.Println("Error selecting upgrade version:", err)
				}
				continue
			}

			if cluster.ClusterVersion.KubernetesVersion == upgradeToVersion.KubernetesVersion {
				fmt.Println("Cluster", cluster.Name, "is already at the desired version")
				continue
			}

			if upgradeClusterDryRun {
				fmt.Println("Dry run mode: would upgrade cluster", cluster.Name, "to version", upgradeToVersion.KubernetesVersion)
				continue
			}

			fmt.Println("Upgrading cluster", cluster.Name, "to version", upgradeToVersion.KubernetesVersion)

			// Create update request with the version slug
			updateRequest := kubernetes.UpdateKubernetesCluster{
				KubernetesVersionIdentity: ptr.To(upgradeToVersion.Identity),
			}
			// Call the API to upgrade the cluster
			_, err = client.Kubernetes().UpdateKubernetesCluster(cmd.Context(), cluster.Identity, updateRequest)
			if err != nil {
				fmt.Println("Error upgrading cluster:", err)
				continue
			}
			if !upgradeAllClusters {
				fmt.Println("Upgrade of cluster", cluster.Name, "initiated successfully. The cluster will be upgraded in the background.")
			}
		}

		if upgradeAllClusters {
			fmt.Println("Upgrade of all clusters initiated successfully. The clusters will be upgraded in the background.")
		}

	},
}

func init() {
	KubernetesCmd.AddCommand(KubernetesUpgradeCmd)

	KubernetesUpgradeCmd.Flags().StringVarP(&upgradeClusterToVersion, "version", "v", "",
		"the version to upgrade to. If not provided, the latest suitable version will be used following Kubernetes version policy (only +1 minor version or patch updates)")
	KubernetesUpgradeCmd.Flags().BoolVar(&upgradeClusterDryRun, "dry-run", true, "print the upgrade request without actually upgrading the cluster")
	KubernetesUpgradeCmd.Flags().BoolVar(&upgradeAllClusters, "all", false, "upgrade all clusters")
}

func getSelectedCluster(args []string) (string, error) {
	if contextstate.OrganisationFlag != "" {
		return contextstate.OrganisationFlag, nil
	}

	if len(args) == 0 && fzf.IsInteractiveMode(os.Stdout) {
		command := fmt.Sprintf("%s kubernetes clusters --no-header", os.Args[0])
		return fzf.InteractiveChoice(command)
	} else if len(args) == 1 {
		return args[0], nil
	} else {
		return "", errors.New("invalid organisation")
	}
}
