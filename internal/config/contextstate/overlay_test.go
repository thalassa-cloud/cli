package contextstate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/thalassa-cloud/cli/internal/credentials"
	"github.com/thalassa-cloud/cli/internal/dirconfig"
)

func TestDirectoryOverlayPrecedence(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(credentials.ResetKeyring)
	t.Setenv(credentials.ThalassaCredentialStoreEnvVar, credentials.StoreKeychain)

	path := filepath.Join(t.TempDir(), ".tcloud")
	globalConfigManager = NewConfigFileContextManager(path)
	t.Cleanup(func() {
		dirconfig.Reset()
		resetAccessorFlags()
	})
	dirconfig.Reset()
	resetAccessorFlags()

	require.NoError(t, globalConfigManager.AddOrMergeContext(Context{
		Name:         "default",
		Organisation: "org-default",
		Project:      "proj-default",
		Users:        Users{Name: "default", User: User{}},
		Servers:      Servers{Name: "api", API: API{Server: "https://api.example.com"}},
	}))
	require.NoError(t, globalConfigManager.AddOrMergeContext(Context{
		Name:         "prod",
		Organisation: "org-prod",
		Project:      "proj-prod",
		Users:        Users{Name: "prod", User: User{}},
		Servers:      Servers{Name: "api-prod", API: API{Server: "https://api.prod.example.com"}},
	}))
	require.NoError(t, globalConfigManager.Set("default"))

	tests := []struct {
		name        string
		setup       func(t *testing.T)
		wantContext string
		wantOrg     string
		wantProject string
		wantServer  string
	}{
		{
			name:        "file current-context when no overlay",
			setup:       func(t *testing.T) {},
			wantContext: "default",
			wantOrg:     "org-default",
			wantProject: "proj-default",
			wantServer:  "https://api.example.com",
		},
		{
			name: "directory context name switches without mutating file",
			setup: func(t *testing.T) {
				dirconfig.SetCurrentForTest(dirconfig.Config{Context: "prod"}, "/repo/.thalassa")
			},
			wantContext: "prod",
			wantOrg:     "org-prod",
			wantProject: "proj-prod",
			wantServer:  "https://api.prod.example.com",
		},
		{
			name: "directory organisation overrides named context",
			setup: func(t *testing.T) {
				dirconfig.SetCurrentForTest(dirconfig.Config{Context: "prod", Organisation: "org-overlay"}, "/repo/.thalassa")
			},
			wantContext: "prod",
			wantOrg:     "org-overlay",
			wantProject: "proj-prod",
			wantServer:  "https://api.prod.example.com",
		},
		{
			name: "flag wins over directory overlay",
			setup: func(t *testing.T) {
				dirconfig.SetCurrentForTest(dirconfig.Config{Context: "prod", Organisation: "org-overlay", Project: "proj-overlay"}, "/repo/.thalassa")
				ContextFlag = "default"
				OrganisationFlag = "org-flag"
				ProjectFlag = "proj-flag"
			},
			wantContext: "default",
			wantOrg:     "org-flag",
			wantProject: "proj-flag",
			wantServer:  "https://api.example.com",
		},
		{
			name: "env wins over directory overlay",
			setup: func(t *testing.T) {
				dirconfig.SetCurrentForTest(dirconfig.Config{Organisation: "org-overlay", Project: "proj-overlay"}, "/repo/.thalassa")
				t.Setenv(ThalassaOrganisationIDEnvVar, "org-env")
				t.Setenv(ThalassaProjectIDEnvVar, "proj-env")
			},
			wantContext: "default",
			wantOrg:     "org-env",
			wantProject: "proj-env",
			wantServer:  "https://api.example.com",
		},
		{
			name: "disabled overlay is ignored",
			setup: func(t *testing.T) {
				dirconfig.SetCurrentForTest(dirconfig.Config{Context: "prod"}, "/repo/.thalassa")
				dirconfig.Disable()
			},
			wantContext: "default",
			wantOrg:     "org-default",
			wantProject: "proj-default",
			wantServer:  "https://api.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirconfig.Reset()
			resetAccessorFlags()
			t.Setenv(ThalassaOrganisationIDEnvVar, "")
			t.Setenv(ThalassaProjectIDEnvVar, "")
			tt.setup(t)

			assert.Equal(t, tt.wantContext, Name())
			assert.Equal(t, tt.wantOrg, Organisation())
			assert.Equal(t, tt.wantProject, Project())
			assert.Equal(t, tt.wantServer, Server())

			cfg, err := GetContextConfiguration()
			require.NoError(t, err)
			assert.Equal(t, tt.wantContext, cfg.Name)

			fileCtx, err := globalConfigManager.Get()
			require.NoError(t, err)
			assert.Equal(t, "default", fileCtx.Name, "directory overlay must not mutate ~/.tcloud current-context")
		})
	}
}

func TestGetByName(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".tcloud")
	manager := NewConfigFileContextManager(path)
	require.NoError(t, manager.AddOrMergeContext(Context{
		Name:         "default",
		Organisation: "org-default",
		Users:        Users{Name: "default", User: User{}},
		Servers:      Servers{Name: "api", API: API{Server: "https://api.example.com"}},
	}))
	require.NoError(t, manager.AddOrMergeContext(Context{
		Name:         "prod",
		Organisation: "org-prod",
		Users:        Users{Name: "prod", User: User{}},
		Servers:      Servers{Name: "api", API: API{Server: "https://api.example.com"}},
	}))
	require.NoError(t, manager.Set("default"))

	got, err := manager.GetByName("prod")
	require.NoError(t, err)
	assert.Equal(t, "prod", got.Name)
	assert.Equal(t, "org-prod", got.Organisation)

	fileCtx, err := manager.Get()
	require.NoError(t, err)
	assert.Equal(t, "default", fileCtx.Name)
}

func resetAccessorFlags() {
	OrganisationFlag = ""
	ProjectFlag = ""
	ContextFlag = ""
	EndpointFlag = ""
	AccessTokenFlag = ""
	PersonalAccessTokenFlag = ""
	OidcClientIDFlag = ""
	OidcClientSecretFlag = ""
	_ = os.Unsetenv(ThalassaOrganisationIDEnvVar)
	_ = os.Unsetenv(ThalassaProjectIDEnvVar)
}
