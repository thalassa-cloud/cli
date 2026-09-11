package kuberesolve

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalassa-cloud/cli/internal/dirconfig"
)

func TestPreferredClusterRef(t *testing.T) {
	t.Cleanup(dirconfig.Reset)

	tests := []struct {
		name     string
		explicit string
		env      map[string]string
		overlay  dirconfig.Config
		disabled bool
		want     string
		wantOK   bool
	}{
		{
			name:   "empty",
			wantOK: false,
		},
		{
			name:     "explicit wins",
			explicit: "from-flag",
			env: map[string]string{
				EnvClusterIdentity: "from-env",
			},
			overlay: dirconfig.Config{Kubernetes: dirconfig.KubernetesConfig{Cluster: "from-dir"}},
			want:    "from-flag",
			wantOK:  true,
		},
		{
			name: "env identity before directory",
			env: map[string]string{
				EnvClusterIdentity: "id-env",
			},
			overlay: dirconfig.Config{Kubernetes: dirconfig.KubernetesConfig{Cluster: "from-dir"}},
			want:    "id-env",
			wantOK:  true,
		},
		{
			name: "env slug when identity unset",
			env: map[string]string{
				EnvClusterSlug: "slug-env",
			},
			want:   "slug-env",
			wantOK: true,
		},
		{
			name: "env name last among env vars",
			env: map[string]string{
				EnvClusterName: "name-env",
			},
			want:   "name-env",
			wantOK: true,
		},
		{
			name:    "directory cluster",
			overlay: dirconfig.Config{Kubernetes: dirconfig.KubernetesConfig{Cluster: "from-dir"}},
			want:    "from-dir",
			wantOK:  true,
		},
		{
			name:     "disabled directory is ignored",
			overlay:  dirconfig.Config{Kubernetes: dirconfig.KubernetesConfig{Cluster: "from-dir"}},
			disabled: true,
			wantOK:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirconfig.Reset()
			t.Setenv(EnvClusterIdentity, "")
			t.Setenv(EnvClusterSlug, "")
			t.Setenv(EnvClusterName, "")
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			if !tt.overlay.Empty() {
				dirconfig.SetCurrentForTest(tt.overlay, "/repo/.thalassa")
			}
			if tt.disabled {
				dirconfig.Disable()
			}

			got, ok := PreferredClusterRef(tt.explicit)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRequireClusterRef(t *testing.T) {
	t.Cleanup(dirconfig.Reset)
	dirconfig.Reset()
	t.Setenv(EnvClusterIdentity, "")
	t.Setenv(EnvClusterSlug, "")
	t.Setenv(EnvClusterName, "")

	_, err := RequireClusterRef("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), ".thalassa")

	got, err := RequireClusterRef("c1")
	require.NoError(t, err)
	assert.Equal(t, "c1", got)
}
