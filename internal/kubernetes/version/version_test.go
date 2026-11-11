package version_test

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/kubernetes"

	"github.com/thalassa-cloud/cli/internal/kubernetes/version"
)

func TestSelectUpgradeVersion(t *testing.T) {
	versions := []kubernetes.KubernetesVersion{
		versionFixture("1.31.1", true),
		versionFixture("1.31.4", true),
		versionFixture("1.32.0", true),
		versionFixture("1.33.0", true),
	}

	tests := []struct {
		name      string
		current   string
		requested string
		want      string
		wantErr   string
	}{
		{
			name:    "prefers newest patch over next minor",
			current: "1.31.1",
			want:    "1.31.4",
		},
		{
			name:    "selects next minor when no patch is available",
			current: "1.31.4",
			want:    "1.32.0",
		},
		{
			name:      "accepts an explicit next minor",
			current:   "1.31.4",
			requested: "1.32.0",
			want:      "1.32.0",
		},
		{
			name:      "accepts an explicit version identity",
			current:   "1.31.1",
			requested: "id-1.31.4",
			want:      "1.31.4",
		},
		{
			name:      "rejects a skipped minor",
			current:   "1.31.1",
			requested: "1.33.0",
			wantErr:   "cannot skip minor versions",
		},
		{
			name:    "rejects a direct jump to a skipped minor",
			current: "1.31.4",
			wantErr: "can only upgrade one minor version at a time",
		},
		{
			name:    "reports the latest version when nothing newer exists",
			current: "1.33.0",
			wantErr: "already at the latest available version",
		},
		{
			name:      "rejects an unknown version",
			current:   "1.31.1",
			requested: "9.9.9",
			wantErr:   "requested version 9.9.9 not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catalog := versions
			if tt.name == "rejects a direct jump to a skipped minor" {
				catalog = []kubernetes.KubernetesVersion{
					versionFixture("1.31.4", true),
					versionFixture("1.33.0", true),
				}
			}

			got, err := version.SelectUpgradeVersion(semver.MustParse(tt.current), catalog, tt.requested)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
				if tt.wantErr == "already at the latest available version" {
					assert.ErrorIs(t, err, version.ErrAlreadyAtLatestVersion)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got.KubernetesVersion)
		})
	}
}

func TestSelectNodePoolUpgradeVersion(t *testing.T) {
	versions := []kubernetes.KubernetesVersion{
		versionFixture("1.31.1", true),
		versionFixture("1.31.4", true),
		versionFixture("1.32.0", true),
		versionFixture("1.32.1", false),
		versionFixture("1.32.2", true),
		versionFixture("1.33.0", true),
	}

	tests := []struct {
		name      string
		current   string
		cluster   string
		requested string
		want      string
		wantErr   string
		wantIs    error
	}{
		{
			name:    "prefers the newest patch within the cluster version",
			current: "1.31.1",
			cluster: "1.32.2",
			want:    "1.31.4",
		},
		{
			name:    "selects the next minor when the cluster allows it",
			current: "1.31.4",
			cluster: "1.32.2",
			want:    "1.32.2",
		},
		{
			name:    "does not move past the cluster version",
			current: "1.31.1",
			cluster: "1.31.4",
			want:    "1.31.4",
		},
		{
			name:      "accepts an explicit version by slug",
			current:   "1.31.1",
			cluster:   "1.32.2",
			requested: "slug-1.32.0",
			want:      "1.32.0",
		},
		{
			name:      "rejects an explicit version newer than the cluster",
			current:   "1.31.1",
			cluster:   "1.31.4",
			requested: "1.32.0",
			wantErr:   "cannot upgrade node pool past the cluster version",
		},
		{
			name:      "rejects a disabled version",
			current:   "1.32.0",
			cluster:   "1.33.0",
			requested: "1.32.1",
			wantErr:   "kubernetes version not found or not enabled",
		},
		{
			name:      "rejects skipping a minor",
			current:   "1.31.1",
			cluster:   "1.33.0",
			requested: "1.33.0",
			wantErr:   "cannot skip minor versions",
		},
		{
			name:    "stops at the cluster version when a newer catalog version exists",
			current: "1.31.4",
			cluster: "1.31.4",
			wantIs:  version.ErrAlreadyAtClusterVersion,
		},
		{
			name:    "reports the latest catalog version when the cluster is ahead of the catalog",
			current: "1.33.0",
			cluster: "1.34.0",
			wantIs:  version.ErrAlreadyAtLatestVersion,
		},
		{
			name:    "ignores a disabled patch and keeps the newest enabled patch",
			current: "1.32.0",
			cluster: "1.33.0",
			want:    "1.32.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := version.SelectNodePoolUpgradeVersion(
				semver.MustParse(tt.current),
				semver.MustParse(tt.cluster),
				versions,
				tt.requested,
			)
			if tt.wantErr != "" || tt.wantIs != nil {
				require.Error(t, err)
				if tt.wantErr != "" {
					assert.ErrorContains(t, err, tt.wantErr)
				}
				if tt.wantIs != nil {
					assert.ErrorIs(t, err, tt.wantIs)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got.KubernetesVersion)
		})
	}
}

func versionFixture(kubernetesVersion string, enabled bool) kubernetes.KubernetesVersion {
	return kubernetes.KubernetesVersion{
		Identity:          "id-" + kubernetesVersion,
		Name:              "name-" + kubernetesVersion,
		Slug:              "slug-" + kubernetesVersion,
		KubernetesVersion: kubernetesVersion,
		Enabled:           enabled,
	}
}
