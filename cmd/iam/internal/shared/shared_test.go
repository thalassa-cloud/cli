package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

func TestKeyValuePairsToMap(t *testing.T) {
	tests := []struct {
		name  string
		pairs []string
		want  map[string]string
	}{
		{name: "empty", pairs: nil, want: map[string]string{}},
		{name: "valid pairs", pairs: []string{"a=b", "c=d"}, want: map[string]string{"a": "b", "c": "d"}},
		{name: "trim", pairs: []string{" a = b "}, want: map[string]string{"a": "b"}},
		{name: "skip invalid", pairs: []string{"nope", "ok=yes"}, want: map[string]string{"ok": "yes"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KeyValuePairsToMap(tt.pairs)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPromptDestructiveUnlessForce_SkipsPromptWhenForce(t *testing.T) {
	proceed, err := PromptDestructiveUnlessForce(true, "will not print")
	assert.NoError(t, err)
	assert.True(t, proceed)
}

func TestParseAccessCredentialScopes(t *testing.T) {
	got, err := ParseAccessCredentialScopes([]string{"api:read", "kubernetes"})
	require.NoError(t, err)
	assert.Equal(t, []clientiam.AccessCredentialsScope{
		clientiam.AccessCredentialsScopeAPIRead,
		clientiam.AccessCredentialsScopeKubernetes,
	}, got)

	_, err = ParseAccessCredentialScopes([]string{"nope"})
	require.Error(t, err)
}
