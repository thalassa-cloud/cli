package credentials

import (
	"errors"

	"github.com/zalando/go-keyring"
)

func isKeyringNotFound(err error) bool {
	return errors.Is(err, keyring.ErrNotFound)
}
