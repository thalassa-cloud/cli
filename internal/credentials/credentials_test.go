package credentials_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/thalassa-cloud/cli/internal/credentials"
)

func TestKeyringStoreRoundTrip(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(credentials.ResetKeyring)

	store := credentials.ResolveStore(credentials.StoreKeychain)
	require.True(t, store.Available())

	secrets := credentials.Secrets{
		Token:        "tc_pat_test",
		AccessToken:  "access",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}
	require.NoError(t, store.Set("default", secrets))

	got, err := store.Get("default")
	require.NoError(t, err)
	assert.Equal(t, secrets, got)

	require.NoError(t, store.Delete("default"))
	_, err = store.Get("default")
	require.Error(t, err)
}

func TestUseKeychainRespectsPreference(t *testing.T) {
	t.Setenv(credentials.ThalassaCredentialStoreEnvVar, credentials.StoreFile)
	assert.False(t, credentials.UseKeychain(credentials.PreferredStore(), credentials.StoreKeychain))

	t.Setenv(credentials.ThalassaCredentialStoreEnvVar, credentials.StoreKeychain)
	assert.True(t, credentials.UseKeychain(credentials.PreferredStore(), credentials.StoreFile))
}

func TestSecretsEmpty(t *testing.T) {
	assert.True(t, credentials.Secrets{}.Empty())
	assert.False(t, credentials.Secrets{Token: "x"}.Empty())
}
