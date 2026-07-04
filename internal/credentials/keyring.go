package credentials

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

type keyringStore struct{}

func (keyringStore) Name() string {
	return StoreKeychain
}

func (keyringStore) Available() bool {
	probeUser := "__tcloud_probe__"
	if err := keyring.Set(KeyringService, probeUser, "probe"); err != nil {
		return false
	}
	_ = keyring.Delete(KeyringService, probeUser)
	return true
}

func (keyringStore) Set(userName string, secrets Secrets) error {
	if userName == "" {
		return fmt.Errorf("user name is required")
	}
	payload, err := marshalSecrets(secrets)
	if err != nil {
		return err
	}
	if secrets.Empty() {
		return keyringStore{}.Delete(userName)
	}
	return keyring.Set(KeyringService, userName, payload)
}

func (keyringStore) Get(userName string) (Secrets, error) {
	if userName == "" {
		return Secrets{}, fmt.Errorf("user name is required")
	}
	raw, err := keyring.Get(KeyringService, userName)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return Secrets{}, err
		}
		return Secrets{}, err
	}
	return unmarshalSecrets(raw)
}

func (keyringStore) Delete(userName string) error {
	if userName == "" {
		return nil
	}
	err := keyring.Delete(KeyringService, userName)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	return err
}

// ResetKeyring removes all CLI credentials from the keychain. Intended for tests.
func ResetKeyring() {
	_ = keyring.DeleteAll(KeyringService)
}
