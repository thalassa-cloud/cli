package labels

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLabelSelector(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		expected map[string]string
	}{
		{
			name:     "single label",
			selector: "key1=value1",
			expected: map[string]string{
				"key1": "value1",
			},
		},
		{
			name:     "multiple labels",
			selector: "key1=value1,key2=value2",
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "multiple labels with spaces",
			selector: "key1=value1, key2=value2",
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "labels with spaces around key and value",
			selector: " key1 = value1 , key2 = value2 ",
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "empty string",
			selector: "",
			expected: map[string]string{},
		},
		{
			name:     "value with equals sign",
			selector: "key1=value=with=equals",
			expected: map[string]string{
				"key1": "value=with=equals",
			},
		},
		{
			name:     "value with spaces",
			selector: "key1=value with spaces",
			expected: map[string]string{
				"key1": "value with spaces",
			},
		},
		{
			name:     "multiple labels with various spacing",
			selector: "env=production, app=backend, version=1.0.0",
			expected: map[string]string{
				"env":     "production",
				"app":     "backend",
				"version": "1.0.0",
			},
		},
		{
			name:     "invalid pair without equals",
			selector: "key1=value1,invalidpair,key2=value2",
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "invalid pair with only key",
			selector: "key1=value1,key2",
			expected: map[string]string{
				"key1": "value1",
			},
		},
		{
			name:     "invalid pair with only equals",
			selector: "key1=value1,=",
			expected: map[string]string{
				"":     "",
				"key1": "value1",
			},
		},
		{
			name:     "empty key",
			selector: "=value1",
			expected: map[string]string{
				"": "value1",
			},
		},
		{
			name:     "empty value",
			selector: "key1=",
			expected: map[string]string{
				"key1": "",
			},
		},
		{
			name:     "whitespace only",
			selector: "   ,  ,  ",
			expected: map[string]string{},
		},
		{
			name:     "duplicate keys (last wins)",
			selector: "key1=value1,key1=value2",
			expected: map[string]string{
				"key1": "value2",
			},
		},
		{
			name:     "real-world example",
			selector: "environment=test,team=backend,region=us-west-1",
			expected: map[string]string{
				"environment": "test",
				"team":        "backend",
				"region":      "us-west-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLabelSelector(tt.selector)
			assert.Equal(t, tt.expected, result, "ParseLabelSelector(%q) should return expected map", tt.selector)
		})
	}
}

