package browseroidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRealmURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr string
	}{
		{
			name: "allows https",
			raw:  "https://login.thalassa.cloud/realms/thalassa-cloud",
		},
		{
			name: "allows localhost http for testing",
			raw:  "http://127.0.0.1:8080/realms/test",
		},
		{
			name:    "rejects remote http",
			raw:     "http://evil.example/realms/test",
			wantErr: "realm URL must use https",
		},
		{
			name:    "rejects missing host",
			raw:     "https://",
			wantErr: "realm URL must use https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRealmURL(tt.raw)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestValidateLoopbackRedirectURL(t *testing.T) {
	require.NoError(t, validateLoopbackRedirectURL("http://127.0.0.1:8765/auth/callback"))
	require.NoError(t, validateLoopbackRedirectURL("http://localhost:8765/auth/callback"))

	err := validateLoopbackRedirectURL("http://evil.example/auth/callback")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "redirect URL must use http://127.0.0.1 or http://localhost")
}
