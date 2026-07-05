package credentials

// fileStore keeps credentials in the config file. It exists so callers can
// treat file and keychain storage through the same interface when needed.
type fileStore struct{}

func (fileStore) Name() string {
	return StoreFile
}

func (fileStore) Available() bool {
	return true
}

func (fileStore) Set(string, Secrets) error {
	return nil
}

func (fileStore) Get(string) (Secrets, error) {
	return Secrets{}, ErrNotAvailable
}

func (fileStore) Delete(string) error {
	return nil
}
