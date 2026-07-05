package contextstate

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/thalassa-cloud/cli/internal/credentials"
)

func TestAccessTokenReadsFromLoadedContext(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(credentials.ResetKeyring)
	t.Setenv(credentials.ThalassaCredentialStoreEnvVar, credentials.StoreKeychain)

	path := filepath.Join(t.TempDir(), ".tcloud")
	globalConfigManager = NewConfigFileContextManager(path)

	require.NoError(t, globalConfigManager.AddOrMergeContext(Context{
		Name: "default",
		Users: Users{
			Name: "default",
			User: User{},
		},
		Servers: Servers{Name: "api", API: API{Server: "https://api.example.com"}},
	}))
	require.NoError(t, globalConfigManager.Set("default"))

	ctx, err := globalConfigManager.Get()
	require.NoError(t, err)
	ctx.Users.User.SetBrowserTokens("browser-access-token", "refresh", time.Now().Add(time.Hour))
	require.NoError(t, CombineConfigContext(ctx))
	require.NoError(t, Save())

	AccessTokenFlag = ""
	t.Setenv(ThalassaAccessTokenEnvVar, "")

	assert.Equal(t, "browser-access-token", AccessToken())
}
