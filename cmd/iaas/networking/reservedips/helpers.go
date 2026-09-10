package reservedips

import (
	"strings"

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

func regionName(rip iaas.ReservedIP) string {
	if rip.Region == nil {
		return ""
	}
	if rip.Region.Name != "" {
		return rip.Region.Name
	}
	if rip.Region.Slug != "" {
		return rip.Region.Slug
	}
	return rip.Region.Identity
}

func attachedTo(rip iaas.ReservedIP) string {
	if rip.AttachedToResourceIdentity == "" {
		return "-"
	}
	if rip.AttachedToResourceType != "" {
		return string(rip.AttachedToResourceType) + ":" + rip.AttachedToResourceIdentity
	}
	return rip.AttachedToResourceIdentity
}

func dashIfEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
