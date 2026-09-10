package snapshotpolicies

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
)

func parseKeyValueSlice(items []string) map[string]string {
	result := make(map[string]string)
	for _, item := range items {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

func regionName(policy iaas.SnapshotPolicy) string {
	if policy.Region == nil {
		return ""
	}
	if policy.Region.Name != "" {
		return policy.Region.Name
	}
	if policy.Region.Slug != "" {
		return policy.Region.Slug
	}
	return policy.Region.Identity
}

func keepCountString(keepCount *int) string {
	if keepCount == nil {
		return "-"
	}
	return strconv.Itoa(*keepCount)
}

func targetSummary(target iaas.SnapshotPolicyTarget) string {
	switch target.Type {
	case iaas.SnapshotPolicyTargetTypeSelector:
		if len(target.Selector) == 0 {
			return "selector"
		}
		pairs := make([]string, 0, len(target.Selector))
		for k, v := range target.Selector {
			pairs = append(pairs, k+"="+v)
		}
		return "selector:" + strings.Join(pairs, ",")
	case iaas.SnapshotPolicyTargetTypeExplicit:
		if len(target.VolumeIdentities) == 0 {
			return "explicit"
		}
		return fmt.Sprintf("explicit:%d volumes", len(target.VolumeIdentities))
	default:
		if target.Type == "" {
			return "-"
		}
		return string(target.Type)
	}
}

func parseTTL(value string) (time.Duration, error) {
	if value == "" {
		return 0, fmt.Errorf("ttl is required")
	}
	ttl, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid ttl %q: %w", value, err)
	}
	if ttl <= 0 {
		return 0, fmt.Errorf("ttl must be greater than 0")
	}
	return ttl, nil
}

func buildTarget(targetType string, selectors []string, volumes []string) (iaas.SnapshotPolicyTarget, error) {
	target := iaas.SnapshotPolicyTarget{
		Type:             iaas.SnapshotPolicyTargetType(targetType),
		Selector:         parseKeyValueSlice(selectors),
		VolumeIdentities: volumes,
	}
	switch target.Type {
	case iaas.SnapshotPolicyTargetTypeSelector:
		if len(target.Selector) == 0 {
			return target, fmt.Errorf("--selector is required when --target-type is selector")
		}
	case iaas.SnapshotPolicyTargetTypeExplicit:
		if len(target.VolumeIdentities) == 0 {
			return target, fmt.Errorf("--volumes is required when --target-type is explicit")
		}
	case "":
		return target, fmt.Errorf("--target-type is required (selector or explicit)")
	default:
		return target, fmt.Errorf("invalid --target-type %q (must be selector or explicit)", targetType)
	}
	return target, nil
}
