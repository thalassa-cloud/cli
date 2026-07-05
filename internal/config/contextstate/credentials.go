package contextstate

import "github.com/thalassa-cloud/cli/internal/credentials"

// ApplyPreferredCredentialStoreForNewLogin marks where freshly authenticated credentials should be stored.
func ApplyPreferredCredentialStoreForNewLogin(user *Users) {
	switch credentials.PreferredStore() {
	case credentials.StoreKeychain:
		user.CredentialStore = credentials.StoreKeychain
	case credentials.StoreFile:
		user.CredentialStore = credentials.StoreFile
	default:
		if credentials.ResolveStore(credentials.StoreKeychain).Available() {
			user.CredentialStore = credentials.StoreKeychain
		} else {
			user.CredentialStore = credentials.StoreFile
		}
	}
}
