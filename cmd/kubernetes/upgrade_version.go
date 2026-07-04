package kubernetes

import (
	"fmt"
	"sort"

	"github.com/blang/semver/v4"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

func findAndValidateRequestedVersion(
	currentVersion semver.Version,
	versions []kubernetes.KubernetesVersion,
	requestedVersion string,
) (*kubernetes.KubernetesVersion, error) {
	for _, version := range versions {
		if version.Name != requestedVersion && version.KubernetesVersion != requestedVersion && version.Slug != requestedVersion {
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
		return &version, nil
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
		versionSemver, _ := semver.Parse(version.KubernetesVersion)
		if highestVersion == nil || versionSemver.GT(latestAvailable) {
			latestAvailable = versionSemver
			highestVersion = &version
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
		return &patchVersions[0], nil
	}

	nextMinorMajorMinor := fmt.Sprintf("%d.%d", currentVersion.Major, currentVersion.Minor+1)
	if minorVersions, hasMinorUpgrade := versionGroups[nextMinorMajorMinor]; hasMinorUpgrade {
		sortVersionsDesc(minorVersions)
		return &minorVersions[0], nil
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
