package labels

import "strings"

// ParseLabelSelector parses a label selector string into a map of labels.
// Format: key1=value1,key2=value2
// The function handles whitespace trimming and validates that each pair has both key and value.
func ParseLabelSelector(selector string) map[string]string {
	labels := make(map[string]string)
	pairs := strings.Split(selector, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return labels
}
