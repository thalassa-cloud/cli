package contextstate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnsureFreshAccessTokenWithoutContext(t *testing.T) {
	path := t.TempDir() + "/.tcloud"
	globalConfigManager = NewConfigFileContextManager(path)

	err := EnsureFreshAccessToken(context.Background())
	require.NoError(t, err)
}
