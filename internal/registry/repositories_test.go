package registry

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryUnmarshalNamespace(t *testing.T) {
	pushed := time.Date(2026, 11, 10, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		payload       string
		wantIdentity  string
		wantNamespace string
		wantImage     string
	}{
		{
			name: "namespace is a string",
			payload: `{
				"identity": "repo-1",
				"image": "api",
				"full_name": "my-app/api",
				"namespace": "my-app",
				"tag_count": 2,
				"artifact_count": 3,
				"total_size_bytes": 1048576,
				"last_pushed_at": "` + pushed.Format(time.RFC3339) + `"
			}`,
			wantIdentity:  "my-app",
			wantNamespace: "my-app",
			wantImage:     "api",
		},
		{
			name: "namespace is an object",
			payload: `{
				"identity": "repo-1",
				"image": "api",
				"full_name": "my-app/api",
				"namespace": {"identity": "ns-1", "namespace": "my-app"}
			}`,
			wantIdentity:  "ns-1",
			wantNamespace: "my-app",
			wantImage:     "api",
		},
		{
			name: "namespace is null",
			payload: `{
				"identity": "repo-1",
				"image": "api",
				"namespace": null
			}`,
			wantImage: "api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var repo Repository
			require.NoError(t, json.Unmarshal([]byte(tt.payload), &repo))
			assert.Equal(t, tt.wantIdentity, repo.Namespace.Identity)
			assert.Equal(t, tt.wantNamespace, repo.Namespace.Namespace)
			assert.Equal(t, tt.wantImage, repo.Image)
			assert.Equal(t, "repo-1", repo.Identity)
		})
	}
}

func TestRepositoryUnmarshalRejectsInvalidNamespace(t *testing.T) {
	var repo Repository
	err := json.Unmarshal([]byte(`{"identity":"repo-1","namespace":1}`), &repo)
	require.Error(t, err)
	assert.ErrorContains(t, err, "namespace")
}
