package registry

import "strings"

func ParseKeyValueSlice(items []string) map[string]string {
	result := make(map[string]string)
	for _, item := range items {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}
