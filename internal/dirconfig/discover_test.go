package dirconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestDiscover(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(t *testing.T, root string) string
		want      Config
		wantFile  string
		wantErr   string
		wantEmpty bool
	}{
		{
			name: "missing file is empty",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				return filepath.Join(root, "repo")
			},
			wantEmpty: true,
		},
		{
			name: "file form",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				dir := filepath.Join(root, "repo")
				require.NoError(t, os.MkdirAll(dir, 0o755))
				writeThalassa(t, filepath.Join(dir, FileName), `context: prod
kubernetes:
  cluster: prod-cluster
`)
				return dir
			},
			want: Config{
				Context:    "prod",
				Kubernetes: KubernetesConfig{Cluster: "prod-cluster"},
			},
			wantFile: FileName,
		},
		{
			name: "directory form",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				dir := filepath.Join(root, "repo")
				cfgDir := filepath.Join(dir, FileName)
				require.NoError(t, os.MkdirAll(cfgDir, 0o755))
				writeThalassa(t, filepath.Join(cfgDir, NestedConfigFile), "context: staging\n")
				return dir
			},
			want: Config{
				Context: "staging",
			},
			wantFile: filepath.Join(FileName, NestedConfigFile),
		},
		{
			name: "nearest child wins over parent",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				parent := filepath.Join(root, "repo")
				child := filepath.Join(parent, "services", "api")
				require.NoError(t, os.MkdirAll(child, 0o755))
				writeThalassa(t, filepath.Join(parent, FileName), "context: parent\n")
				writeThalassa(t, filepath.Join(child, FileName), "context: child\n")
				return child
			},
			want:     Config{Context: "child"},
			wantFile: FileName,
		},
		{
			name: "walks up to parent when child has none",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				parent := filepath.Join(root, "repo")
				child := filepath.Join(parent, "nested")
				require.NoError(t, os.MkdirAll(child, 0o755))
				writeThalassa(t, filepath.Join(parent, FileName), "organisation: acme\n")
				return child
			},
			want:     Config{Organisation: "acme"},
			wantFile: FileName,
		},
		{
			name: "directory without config.yaml continues walk",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				parent := filepath.Join(root, "repo")
				child := filepath.Join(parent, "app")
				require.NoError(t, os.MkdirAll(filepath.Join(child, FileName), 0o755))
				writeThalassa(t, filepath.Join(parent, FileName), "project: platform\n")
				return child
			},
			want:     Config{Project: "platform"},
			wantFile: FileName,
		},
		{
			name: "invalid yaml fails closed with path",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				dir := filepath.Join(root, "repo")
				require.NoError(t, os.MkdirAll(dir, 0o755))
				writeThalassa(t, filepath.Join(dir, FileName), "context: [\n")
				return dir
			},
			wantErr: FileName,
		},
		{
			name: "oversized file fails closed",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				dir := filepath.Join(root, "repo")
				require.NoError(t, os.MkdirAll(dir, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, FileName), []byte(strings.Repeat("a", MaxSize+1)), 0o644))
				return dir
			},
			wantErr: "larger than",
		},
		{
			name: "unknown secret keys are ignored",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				dir := filepath.Join(root, "repo")
				require.NoError(t, os.MkdirAll(dir, 0o755))
				writeThalassa(t, filepath.Join(dir, FileName), `context: prod
token: tc_pat_should_not_be_used
accessToken: leaked
kubernetes:
  cluster: prod-cluster
`)
				return dir
			},
			want: Config{
				Context:    "prod",
				Kubernetes: KubernetesConfig{Cluster: "prod-cluster"},
			},
			wantFile: FileName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			cwd := tt.setup(t, root)
			require.NoError(t, os.MkdirAll(cwd, 0o755))

			got, path, err := Discover(cwd)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			if tt.wantEmpty {
				assert.True(t, got.Empty())
				assert.Empty(t, path)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Truef(t, strings.HasSuffix(path, filepath.FromSlash(tt.wantFile)), "path %q should end with %q", path, tt.wantFile)
		})
	}
}

func TestLoadRespectsDisableEnv(t *testing.T) {
	dir := t.TempDir()
	writeThalassa(t, filepath.Join(dir, FileName), "context: prod\n")
	t.Chdir(dir)
	t.Setenv(EnvDisable, "0")
	t.Cleanup(Reset)
	Reset()

	require.NoError(t, Load())
	assert.True(t, Disabled())
	assert.Empty(t, Current().Path)
	assert.True(t, Current().Config.Empty())
}

func TestLoadCachesFile(t *testing.T) {
	dir := t.TempDir()
	writeThalassa(t, filepath.Join(dir, FileName), "context: prod\nkubernetes:\n  cluster: c1\n")
	t.Chdir(dir)
	t.Cleanup(Reset)
	Reset()

	require.NoError(t, Load())
	got := Current()
	assert.Equal(t, "prod", got.Config.Context)
	assert.Equal(t, "c1", got.Config.Kubernetes.Cluster)
	assert.Equal(t, filepath.Join(dir, FileName), got.Path)
}

func TestDisableClearsCurrent(t *testing.T) {
	t.Cleanup(Reset)
	Reset()
	SetCurrentForTest(Config{Context: "prod"}, "/tmp/.thalassa")
	assert.Equal(t, "prod", Current().Config.Context)

	Disable()
	assert.True(t, Disabled())
	assert.True(t, Current().Config.Empty())
}

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)
	cfg := Config{Context: "prod", Kubernetes: KubernetesConfig{Cluster: "c1"}}

	require.NoError(t, WriteFile(path, cfg, false))
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var got Config
	require.NoError(t, yaml.Unmarshal(data, &got))
	assert.Equal(t, cfg, got)

	err = WriteFile(path, cfg, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	cfg.Kubernetes.Cluster = "c2"
	require.NoError(t, WriteFile(path, cfg, true))
}

func TestWriteFileRefusesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)
	require.NoError(t, os.Mkdir(path, 0o755))
	err := WriteFile(path, Config{Context: "prod"}, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is a directory")
}

func TestWriteFileRejectsEmpty(t *testing.T) {
	err := WriteFile(filepath.Join(t.TempDir(), FileName), Config{}, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no values")
}

func writeThalassa(t *testing.T, path, contents string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o644))
}
