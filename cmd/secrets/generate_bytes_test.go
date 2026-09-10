package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateGenerateBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		n       int
		wantErr string
	}{
		{name: "min", n: 16},
		{name: "default", n: 32},
		{name: "max", n: 4096},
		{name: "too small", n: 10, wantErr: "between 16 and 4096"},
		{name: "too large", n: 5000, wantErr: "between 16 and 4096"},
		{name: "zero", n: 0, wantErr: "between 16 and 4096"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateGenerateBytes(tt.n)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
