package version

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

// ErrAlreadyAtLatestVersion is returned when no newer version is available.
var ErrAlreadyAtLatestVersion = errors.New("already at the latest available version")

// ErrAlreadyAtClusterVersion is returned when a node pool is not behind the cluster control plane.
var ErrAlreadyAtClusterVersion = errors.New("already at the cluster version")

// SelectUpgradeVersion selects the appropriate upgrade version based on the current version
// following Kubernetes versioning policy (max 1 minor version upgrade at a time).
func SelectUpgradeVersion(currentVersion semver.Version, versions []kubernetes.KubernetesVersion, requestedVersion string) (*kubernetes.KubernetesVersion, error) {
	if requestedVersion != "" {
		return findAndValidateRequestedVersion(currentVersion, versions, requestedVersion)
	}

	newerVersions, upgradeCandidates := collectUpgradeCandidates(currentVersion, versions)
	if len(newerVersions) == 0 {
		return nil, fmt.Errorf("%w (%s)", ErrAlreadyAtLatestVersion, currentVersion)
	}
	if len(upgradeCandidates) == 0 {
		return nil, upgradePolicyBlockedError(currentVersion, newerVersions)
	}

	return selectBestUpgradeCandidate(currentVersion, upgradeCandidates)
}

// SelectNodePoolUpgradeVersion selects a node pool version.
// Automatic selection follows the cluster upgrade policy and never moves past the cluster control plane.
// An explicit version must be enabled, must satisfy the same policy, and must not be newer than the cluster.
func SelectNodePoolUpgradeVersion(currentVersion, clusterVersion semver.Version, versions []kubernetes.KubernetesVersion, requestedVersion string) (*kubernetes.KubernetesVersion, error) {
	if requestedVersion != "" {
		if match, ok := findVersion(versions, requestedVersion); ok && !match.Enabled {
			return nil, fmt.Errorf("kubernetes version not found or not enabled: %s", requestedVersion)
		}
		selected, err := findAndValidateRequestedVersion(currentVersion, filterEnabled(versions), requestedVersion)
		if err != nil {
			return nil, err
		}
		requestedSemver, err := semver.Parse(selected.KubernetesVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid semver in requested version %s: %w", selected.KubernetesVersion, err)
		}
		if requestedSemver.GT(clusterVersion) {
			return nil, fmt.Errorf("cannot upgrade node pool past the cluster version %s (requested %s)", clusterVersion, requestedSemver)
		}
		return selected, nil
	}

	selected, err := SelectUpgradeVersion(currentVersion, filterNotNewerThan(filterEnabled(versions), clusterVersion), "")
	if err != nil {
		if errors.Is(err, ErrAlreadyAtLatestVersion) && !currentVersion.LT(clusterVersion) {
			return nil, fmt.Errorf("%w (%s)", ErrAlreadyAtClusterVersion, clusterVersion)
		}
		return nil, err
	}
	return selected, nil
}

func findAndValidateRequestedVersion(
	currentVersion semver.Version,
	versions []kubernetes.KubernetesVersion,
	requestedVersion string,
) (*kubernetes.KubernetesVersion, error) {
	for _, version := range versions {
		if !matchesVersion(version, requestedVersion) {
			continue
		}
		requestedSemver, err := semver.Parse(version.KubernetesVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid semver in requested version %s: %w", version.KubernetesVersion, err)
		}
		if requestedSemver.Major > currentVersion.Major {
			return nil, fmt.Errorf("cannot upgrade across major versions from %s to %s", currentVersion, requestedSemver)
		}
		if requestedSemver.Major == currentVersion.Major && requestedSemver.Minor > currentVersion.Minor+1 {
			return nil, fmt.Errorf("cannot skip minor versions from %s to %s", currentVersion, requestedSemver)
		}
		matched := version
		return &matched, nil
	}
	return nil, fmt.Errorf("requested version %s not found", requestedVersion)
}

func collectUpgradeCandidates(currentVersion semver.Version, versions []kubernetes.KubernetesVersion) ([]kubernetes.KubernetesVersion, []kubernetes.KubernetesVersion) {
	var newerVersions []kubernetes.KubernetesVersion
	var upgradeCandidates []kubernetes.KubernetesVersion

	for _, version := range versions {
		versionSemver, err := semver.Parse(version.KubernetesVersion)
		if err != nil {
			continue
		}
		if !versionSemver.GT(currentVersion) {
			continue
		}
		newerVersions = append(newerVersions, version)
		if versionSemver.Major > currentVersion.Major {
			continue
		}
		if versionSemver.Major == currentVersion.Major && versionSemver.Minor > currentVersion.Minor+1 {
			continue
		}
		upgradeCandidates = append(upgradeCandidates, version)
	}

	return newerVersions, upgradeCandidates
}

func upgradePolicyBlockedError(currentVersion semver.Version, newerVersions []kubernetes.KubernetesVersion) error {
	var latestAvailable semver.Version
	var highestVersion *kubernetes.KubernetesVersion

	for _, version := range newerVersions {
		versionSemver, err := semver.Parse(version.KubernetesVersion)
		if err != nil {
			continue
		}
		if highestVersion == nil || versionSemver.GT(latestAvailable) {
			latestAvailable = versionSemver
			matched := version
			highestVersion = &matched
		}
	}

	if highestVersion == nil {
		return fmt.Errorf("no suitable upgrade path found according to versioning policy (can only upgrade to next minor version). Current version: %s", currentVersion)
	}
	if latestAvailable.Major > currentVersion.Major {
		return fmt.Errorf("cannot upgrade directly to latest version %s (major version upgrade not allowed). Current version: %s",
			latestAvailable, currentVersion)
	}
	if latestAvailable.Minor > currentVersion.Minor+1 {
		nextMinorVersion := fmt.Sprintf("%d.%d.0", currentVersion.Major, currentVersion.Minor+1)
		return fmt.Errorf("cannot upgrade directly to latest version %s (can only upgrade one minor version at a time). Current version: %s. Try first upgrading to %s",
			latestAvailable, currentVersion, nextMinorVersion)
	}
	return fmt.Errorf("no suitable upgrade path found according to versioning policy (can only upgrade to next minor version). Current version: %s", currentVersion)
}

func selectBestUpgradeCandidate(currentVersion semver.Version, upgradeCandidates []kubernetes.KubernetesVersion) (*kubernetes.KubernetesVersion, error) {
	versionGroups := make(map[string][]kubernetes.KubernetesVersion)
	for _, candidate := range upgradeCandidates {
		candidateSemver, err := semver.Parse(candidate.KubernetesVersion)
		if err != nil {
			continue
		}
		majorMinor := fmt.Sprintf("%d.%d", candidateSemver.Major, candidateSemver.Minor)
		versionGroups[majorMinor] = append(versionGroups[majorMinor], candidate)
	}

	currentMajorMinor := fmt.Sprintf("%d.%d", currentVersion.Major, currentVersion.Minor)
	if patchVersions, hasPatchUpgrades := versionGroups[currentMajorMinor]; hasPatchUpgrades {
		sortVersionsDesc(patchVersions)
		selected := patchVersions[0]
		return &selected, nil
	}

	nextMinorMajorMinor := fmt.Sprintf("%d.%d", currentVersion.Major, currentVersion.Minor+1)
	if minorVersions, hasMinorUpgrade := versionGroups[nextMinorMajorMinor]; hasMinorUpgrade {
		sortVersionsDesc(minorVersions)
		selected := minorVersions[0]
		return &selected, nil
	}

	return nil, fmt.Errorf("no suitable upgrade version found")
}

func sortVersionsDesc(versions []kubernetes.KubernetesVersion) {
	sort.Slice(versions, func(i, j int) bool {
		iSemver, errI := semver.Parse(versions[i].KubernetesVersion)
		jSemver, errJ := semver.Parse(versions[j].KubernetesVersion)
		if errI != nil {
			return false
		}
		if errJ != nil {
			return true
		}
		return iSemver.GT(jSemver)
	})
}

func matchesVersion(version kubernetes.KubernetesVersion, requested string) bool {
	return strings.EqualFold(version.Identity, requested) ||
		strings.EqualFold(version.Name, requested) ||
		strings.EqualFold(version.Slug, requested) ||
		strings.EqualFold(version.KubernetesVersion, requested)
}

func findVersion(versions []kubernetes.KubernetesVersion, requested string) (kubernetes.KubernetesVersion, bool) {
	for _, version := range versions {
		if matchesVersion(version, requested) {
			return version, true
		}
	}
	return kubernetes.KubernetesVersion{}, false
}

func filterEnabled(versions []kubernetes.KubernetesVersion) []kubernetes.KubernetesVersion {
	enabled := make([]kubernetes.KubernetesVersion, 0, len(versions))
	for _, version := range versions {
		if version.Enabled {
			enabled = append(enabled, version)
		}
	}
	return enabled
}

func filterNotNewerThan(versions []kubernetes.KubernetesVersion, ceiling semver.Version) []kubernetes.KubernetesVersion {
	capped := make([]kubernetes.KubernetesVersion, 0, len(versions))
	for _, version := range versions {
		parsed, err := semver.Parse(version.KubernetesVersion)
		if err != nil {
			continue
		}
		if parsed.GT(ceiling) {
			continue
		}
		capped = append(capped, version)
	}
	return capped
}
