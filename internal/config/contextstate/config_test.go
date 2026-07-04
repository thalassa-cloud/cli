package contextstate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v3"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/credentials"
)

func TestConfigManagerMigratesPlaintextCredentialsToKeychain(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(credentials.ResetKeyring)
	t.Setenv(credentials.ThalassaCredentialStoreEnvVar, credentials.StoreKeychain)

	dir := t.TempDir()
	path := filepath.Join(dir, ".tcloud")
	require.NoError(t, os.WriteFile(path, []byte(`configVersion: v1
current-context: default
contexts:
  - name: default
    context:
      api: api.thalassa.cloud
      user: default
      organisation: acme
servers:
  - name: api.thalassa.cloud
    api:
      server: https://api.thalassa.cloud
users:
  - name: default
    user:
      token: tc_pat_secret
`), 0o644))

	manager := contextstate.NewConfigFileContextManager(path)
	require.NoError(t, manager.Load())

	ctx, err := manager.Get()
	require.NoError(t, err)
	assert.Equal(t, "tc_pat_secret", ctx.Users.User.Token)

	onDisk, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.NotContains(t, string(onDisk), "tc_pat_secret")

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	var saved struct {
		Users []struct {
			Name            string `yaml:"name"`
			CredentialStore string `yaml:"credentialStore"`
			User            struct {
				Token string `yaml:"token"`
			} `yaml:"user"`
		} `yaml:"users"`
	}
	require.NoError(t, yaml.Unmarshal(onDisk, &saved))
	require.Len(t, saved.Users, 1)
	assert.Equal(t, credentials.StoreKeychain, saved.Users[0].CredentialStore)
	assert.Empty(t, saved.Users[0].User.Token)

	store := credentials.ResolveStore(credentials.StoreKeychain)
	got, err := store.Get("default")
	require.NoError(t, err)
	assert.Equal(t, "tc_pat_secret", got.Token)
}

func TestConfigManagerSanitizedConfigRedactsSecrets(t *testing.T) {
	manager := contextstate.NewConfigFileContextManager(filepath.Join(t.TempDir(), ".tcloud"))
	require.NoError(t, manager.AddOrMergeContext(contextstate.Context{
		Name: "default",
		Users: contextstate.Users{
			Name: "default",
			User: contextstate.User{Token: "tc_pat_secret"},
		},
		Servers: contextstate.Servers{Name: "api", API: contextstate.API{Server: "https://api.example.com"}},
	}))
	require.NoError(t, manager.Set("default"))

	sanitized := manager.SanitizedConfig()
	require.Len(t, sanitized.Users, 1)
	assert.Empty(t, sanitized.Users[0].User.Token)

	ctx, err := manager.Get()
	require.NoError(t, err)
	assert.Equal(t, "tc_pat_secret", ctx.Users.User.Token)
}

func TestConfigManagerLoadDoesNotFixPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tcloud")
	require.NoError(t, os.WriteFile(path, []byte(`configVersion: v1
current-context: default
contexts:
  - name: default
    context:
      api: api
      user: default
servers:
  - name: api
    api:
      server: https://api.example.com
users:
  - name: default
    user: {}
`), 0o644))

	manager := contextstate.NewConfigFileContextManager(path)
	require.NoError(t, manager.Load())

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o644), info.Mode().Perm())
}

func TestConfigManagerFixPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tcloud")
	require.NoError(t, os.WriteFile(path, []byte("config"), 0o644))

	manager := contextstate.NewConfigFileContextManager(path)
	require.NoError(t, manager.FixPermissions())

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}
