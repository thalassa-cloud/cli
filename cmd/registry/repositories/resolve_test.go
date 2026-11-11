package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/containerregistry"
)

func TestResolveNamespaceByName(t *testing.T) {
	namespaces := []containerregistry.ContainerRegistryNamespace{
		{Identity: "id-app", Namespace: "my-app"},
		{Identity: "id-other", Namespace: "other"},
		{Identity: "id-app-2", Namespace: "My-App"},
	}

	tests := []struct {
		name      string
		input     string
		items     []containerregistry.ContainerRegistryNamespace
		wantID    string
		wantError string
	}{
		{
			name:   "matches the namespace name",
			input:  "other",
			items:  namespaces,
			wantID: "id-other",
		},
		{
			name:   "matches a name without regard to case when there is one match",
			input:  "OTHER",
			items:  []containerregistry.ContainerRegistryNamespace{{Identity: "id-other", Namespace: "other"}},
			wantID: "id-other",
		},
		{
			name:      "rejects an identity passed as a name",
			input:     "id-app",
			items:     namespaces,
			wantError: "use --id",
		},
		{
			name:      "reports an unknown name",
			input:     "missing",
			items:     namespaces,
			wantError: "namespace not found",
		},
		{
			name:  "rejects an ambiguous name",
			input: "my-app",
			items: []containerregistry.ContainerRegistryNamespace{
				{Identity: "id-a", Namespace: "my-app"},
				{Identity: "id-b", Namespace: "my-app"},
			},
			wantError: "more than one registry namespace",
		},
		{
			name:      "requires a name",
			input:     "  ",
			items:     namespaces,
			wantError: "namespace name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveNamespaceByName(tt.items, tt.input)
			if tt.wantError != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.wantID, got.Identity)
		})
	}
}
