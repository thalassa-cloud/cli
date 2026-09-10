package shared

import (
	"strings"
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
			assert.Equal(t, tt.want, got, "should match")
		})
	}
}

func TestPromptDestructiveUnlessForce_SkipsPromptWhenForce(t *testing.T) {
	proceed, err := PromptDestructiveUnlessForce(true, "will not print")
	assert.NoError(t, err, "should not error")
	assert.True(t, proceed, "should proceed")
}

func TestPromptStringFrom(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultValue string
		want         string
		wantErr      bool
	}{
		{name: "uses input", input: "restored-db\n", defaultValue: "default", want: "restored-db"},
		{name: "uses default on empty", input: "\n", defaultValue: "default", want: "default"},
		{name: "eof without newline", input: "", defaultValue: "default", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PromptStringFrom(strings.NewReader(tt.input), "Name", tt.defaultValue)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseAccessCredentialScopes(t *testing.T) {
	got, err := ParseAccessCredentialScopes([]string{"api:read", "kubernetes"})
	require.NoError(t, err, "should not error")
	assert.Equal(t, []clientiam.AccessCredentialsScope{
		clientiam.AccessCredentialsScopeAPIRead,
		clientiam.AccessCredentialsScopeKubernetes,
	}, got)

	_, err = ParseAccessCredentialScopes([]string{"nope"})
	require.Error(t, err, "should error")
}
