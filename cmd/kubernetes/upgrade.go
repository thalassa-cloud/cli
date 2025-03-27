package kubernetes

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/fzf"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/pkg/kubernetesclient"
	"github.com/thalassa-cloud/client-go/pkg/thalassa"

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
	Versions   []kubernetesclient.KubernetesVersion
}

// selectUpgradeVersion selects the appropriate upgrade version based on the current version
// following Kubernetes versioning policy (max 1 minor version upgrade at a time)
func selectUpgradeVersion(currentVersion semver.Version, versions []kubernetesclient.KubernetesVersion, requestedVersion string) (*kubernetesclient.KubernetesVersion, error) {
	// If specific version requested, find and validate it
	if requestedVersion != "" {
		for _, version := range versions {
			if version.Name == requestedVersion || version.KubernetesVersion == requestedVersion || version.Slug == requestedVersion {
				requestedSemver, err := semver.Parse(version.KubernetesVersion)
				if err != nil {
					return nil, fmt.Errorf("invalid semver in requested version %s: %w", version.KubernetesVersion, err)
				}
				// Verify it follows the versioning policy (max 1 minor version upgrade)
				if requestedSemver.Major > currentVersion.Major {
					return nil, fmt.Errorf("cannot upgrade across major versions from %s to %s", currentVersion, requestedSemver)
				}
				if requestedSemver.Major == currentVersion.Major && requestedSemver.Minor > currentVersion.Minor+1 {
					return nil, fmt.Errorf("cannot skip minor versions from %s to %s", currentVersion, requestedSemver)
				}
				return &version, nil
			}
		}
		return nil, fmt.Errorf("requested version %s not found", requestedVersion)
	}

	// Check if there are any newer versions at all
	var newerVersions []kubernetesclient.KubernetesVersion
	for _, version := range versions {
		versionSemver, err := semver.Parse(version.KubernetesVersion)
		if err != nil {
			continue // Skip invalid versions
		}
		if versionSemver.GT(currentVersion) {
			newerVersions = append(newerVersions, version)
		}
	}

	if len(newerVersions) == 0 {
		return nil, fmt.Errorf("already at the latest available version (%s)", currentVersion)
	}

	// Convert versions to semver for proper comparison
	var upgradeCandidates []kubernetesclient.KubernetesVersion
	for _, version := range versions {
		versionSemver, err := semver.Parse(version.KubernetesVersion)
		if err != nil {
			// Skip invalid versions
			continue
		}
		if versionSemver.GT(currentVersion) {
			// Filter by versioning policy (max 1 minor version upgrade)
			if versionSemver.Major > currentVersion.Major {
				continue // Skip major version upgrades
			}
			if versionSemver.Major == currentVersion.Major && versionSemver.Minor > currentVersion.Minor+1 {
				continue // Skip minor versions beyond +1
			}
			upgradeCandidates = append(upgradeCandidates, version)
		}
	}

	if len(upgradeCandidates) == 0 {
		// We have newer versions but none that meet our policy
		var latestAvailable semver.Version
		var highestVersion *kubernetesclient.KubernetesVersion

		for _, version := range newerVersions {
			versionSemver, _ := semver.Parse(version.KubernetesVersion)
			if highestVersion == nil || versionSemver.GT(latestAvailable) {
				latestAvailable = versionSemver
				highestVersion = &version
			}
		}

		if highestVersion != nil {
			if latestAvailable.Major > currentVersion.Major {
				return nil, fmt.Errorf("cannot upgrade directly to latest version %s (major version upgrade not allowed). Current version: %s",
					latestAvailable, currentVersion)
			} else if latestAvailable.Minor > currentVersion.Minor+1 {
				// The cluster needs to be upgraded in multiple steps
				nextMinorVersion := fmt.Sprintf("%d.%d.0", currentVersion.Major, currentVersion.Minor+1)
				return nil, fmt.Errorf("cannot upgrade directly to latest version %s (can only upgrade one minor version at a time). Current version: %s. Try first upgrading to %s",
					latestAvailable, currentVersion, nextMinorVersion)
			}
		}

		return nil, fmt.Errorf("no suitable upgrade path found according to versioning policy (can only upgrade to next minor version). Current version: %s", currentVersion)
	}

	// Group candidates by major.minor to identify patch versions
	versionGroups := make(map[string][]kubernetesclient.KubernetesVersion)

	for _, candidate := range upgradeCandidates {
		candidateSemver, err := semver.Parse(candidate.KubernetesVersion)
		if err != nil {
			continue // Skip invalid versions
		}
		majorMinor := fmt.Sprintf("%d.%d", candidateSemver.Major, candidateSemver.Minor)
		versionGroups[majorMinor] = append(versionGroups[majorMinor], candidate)
	}

	// First, check if there are patch upgrades available for current major.minor
	currentMajorMinor := fmt.Sprintf("%d.%d", currentVersion.Major, currentVersion.Minor)
	if patchVersions, hasPatchUpgrades := versionGroups[currentMajorMinor]; hasPatchUpgrades {
		// Sort patch versions to get the latest one
		sort.Slice(patchVersions, func(i, j int) bool {
			iSemver, errI := semver.Parse(patchVersions[i].KubernetesVersion)
			jSemver, errJ := semver.Parse(patchVersions[j].KubernetesVersion)
			// If either parse fails, put it at the end
			if errI != nil {
				return false
			}
			if errJ != nil {
				return true
			}
			return iSemver.GT(jSemver)
		})
		// Use the latest patch version
		return &patchVersions[0], nil
	}

	// If no patch upgrades, find the next minor version (should be minor+1)
	nextMinorMajorMinor := fmt.Sprintf("%d.%d", currentVersion.Major, currentVersion.Minor+1)
	if minorVersions, hasMinorUpgrade := versionGroups[nextMinorMajorMinor]; hasMinorUpgrade {
		// Sort minor versions to get the latest patch within this minor
		sort.Slice(minorVersions, func(i, j int) bool {
			iSemver, errI := semver.Parse(minorVersions[i].KubernetesVersion)
			jSemver, errJ := semver.Parse(minorVersions[j].KubernetesVersion)
			// If either parse fails, put it at the end
			if errI != nil {
				return false
			}
			if errJ != nil {
				return true
			}
			return iSemver.GT(jSemver)
		})
		// Use the latest patch within next minor
		return &minorVersions[0], nil
	}

	// If we got here, no suitable upgrade was found
	return nil, fmt.Errorf("no suitable upgrade version found")
}

var KubernetesUpgradeCmd = &cobra.Command{
	Use:     "upgrade <cluster>",
	Short:   "Upgrade a Kubernetes cluster",
	Aliases: []string{"u"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			fmt.Println("Error creating client:", err)
			return
		}

		versions, err := client.Kubernetes().ListKubernetesVersions(cmd.Context())
		if err != nil {
			fmt.Println("Error getting versions:", err)
			return
		}

		clusters := []kubernetesclient.KubernetesCluster{}

		if upgradeAllClusters {
			fmt.Println("Upgrading all clusters")
			clusters, err = client.Kubernetes().ListKubernetesClusters(cmd.Context())
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
			updateRequest := kubernetesclient.UpdateKubernetesCluster{
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

	KubernetesUpgradeCmd.Flags().StringVarP(&upgradeClusterToVersion, "version", "v", "", "the version to upgrade to. If not provided, the latest suitable version will be used following Kubernetes version policy (only +1 minor version or patch updates)")
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
