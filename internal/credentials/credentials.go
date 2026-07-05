package credentials

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

const (
	StoreAuto     = "auto"
	StoreKeychain = "keychain"
	StoreFile     = "file"

	KeyringService = "Thalassa Cloud CLI"

	ThalassaCredentialStoreEnvVar = "THALASSA_CREDENTIAL_STORE"
)

var ErrNotAvailable = errors.New("credential store not available")

// Secrets holds sensitive user credentials.
type Secrets struct {
	Token             string `json:"token,omitempty"`
	AccessToken       string `json:"accessToken,omitempty"`
	RefreshToken      string `json:"refreshToken,omitempty"`
	AccessTokenExpiry string `json:"accessTokenExpiry,omitempty"`
	ClientID          string `json:"clientID,omitempty"`
	ClientSecret      string `json:"clientSecret,omitempty"`
}

func (s Secrets) Empty() bool {
	return s.Token == "" && s.AccessToken == "" && s.RefreshToken == "" && s.ClientID == "" && s.ClientSecret == ""
}

// Store persists credentials outside the config file.
type Store interface {
	Name() string
	Available() bool
	Set(userName string, secrets Secrets) error
	Get(userName string) (Secrets, error)
	Delete(userName string) error
}

// PreferredStore reads THALASSA_CREDENTIAL_STORE and returns a normalized store preference.
func PreferredStore() string {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(ThalassaCredentialStoreEnvVar))) {
	case StoreKeychain:
		return StoreKeychain
	case StoreFile:
		return StoreFile
	default:
		return StoreAuto
	}
}

// ResolveStore returns the credential store to use for the given preference.
func ResolveStore(preference string) Store {
	switch preference {
	case StoreKeychain:
		return keyringStore{}
	case StoreFile:
		return fileStore{}
	default:
		kr := keyringStore{}
		if kr.Available() {
			return kr
		}
		return fileStore{}
	}
}

// UseKeychain reports whether secrets for the given user store should be kept in the keychain.
func UseKeychain(preference, userStore string) bool {
	switch preference {
	case StoreKeychain:
		return true
	case StoreFile:
		return false
	default:
		if userStore == StoreKeychain {
			return true
		}
		if userStore == StoreFile {
			return false
		}
		return keyringStore{}.Available()
	}
}

func marshalSecrets(secrets Secrets) (string, error) {
	data, err := json.Marshal(secrets)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshalSecrets(raw string) (Secrets, error) {
	var secrets Secrets
	if raw == "" {
		return secrets, nil
	}
	err := json.Unmarshal([]byte(raw), &secrets)
	return secrets, err
}

// IsNotFound reports whether err indicates missing credentials in the keychain.
func IsNotFound(err error) bool {
	return isKeyringNotFound(err)
}
