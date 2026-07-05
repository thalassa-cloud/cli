package contextstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"

	"github.com/thalassa-cloud/cli/internal/credentials"
)

func TestApplyPreferredCredentialStoreForNewLoginUsesKeychainWhenAvailable(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(credentials.ResetKeyring)

	user := Users{Name: "default"}
	ApplyPreferredCredentialStoreForNewLogin(&user)
	assert.Equal(t, credentials.StoreKeychain, user.CredentialStore)
}

func TestApplyPreferredCredentialStoreForNewLoginRespectsFilePreference(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(credentials.ResetKeyring)
	t.Setenv(credentials.ThalassaCredentialStoreEnvVar, credentials.StoreFile)

	user := Users{Name: "default"}
	ApplyPreferredCredentialStoreForNewLogin(&user)
	assert.Equal(t, credentials.StoreFile, user.CredentialStore)
}
