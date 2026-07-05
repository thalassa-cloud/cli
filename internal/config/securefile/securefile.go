package securefile

import "os"

const privateFileMode = 0o600

// PrivateFileMode is the recommended permission mode for config files containing secrets.
const PrivateFileMode = privateFileMode

// Write writes data and always enforces owner-only permissions, including when
// overwriting an existing file whose mode was too permissive.
func Write(name string, data []byte) error {
	if err := os.WriteFile(name, data, privateFileMode); err != nil {
		return err
	}
	return os.Chmod(name, privateFileMode)
}

// IsPermissive reports whether the given permission bits are looser than owner-only.
func IsPermissive(mode os.FileMode) bool {
	return mode.Perm() != privateFileMode
}

// CheckPermissions returns the current file mode and whether it is too permissive.
func CheckPermissions(name string) (os.FileMode, bool, error) {
	info, err := os.Stat(name)
	if err != nil {
		return 0, false, err
	}
	mode := info.Mode().Perm()
	return mode, IsPermissive(mode), nil
}

// EnsurePermissions corrects permissions on an existing file.
func EnsurePermissions(name string) error {
	info, err := os.Stat(name)
	if err != nil {
		return err
	}
	if info.Mode().Perm() == privateFileMode {
		return nil
	}
	return os.Chmod(name, privateFileMode)
}
