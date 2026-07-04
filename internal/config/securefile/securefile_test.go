package securefile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalassa-cloud/cli/internal/config/securefile"
)

func TestWriteEnforcesPermissionsOnExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")

	require.NoError(t, os.WriteFile(path, []byte("old"), 0o644))
	require.NoError(t, securefile.Write(path, []byte("new")))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	assert.Equal(t, "new", readFile(t, path))
}

func TestCheckPermissionsDetectsPermissiveMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")

	require.NoError(t, os.WriteFile(path, []byte("secret"), 0o644))

	mode, permissive, err := securefile.CheckPermissions(path)
	require.NoError(t, err)
	assert.True(t, permissive)
	assert.Equal(t, os.FileMode(0o644), mode)

	require.NoError(t, securefile.EnsurePermissions(path))
	mode, permissive, err = securefile.CheckPermissions(path)
	require.NoError(t, err)
	assert.False(t, permissive)
	assert.Equal(t, os.FileMode(0o600), mode)
}

func TestEnsurePermissionsFixesExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")

	require.NoError(t, os.WriteFile(path, []byte("secret"), 0o644))
	require.NoError(t, securefile.EnsurePermissions(path))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}
