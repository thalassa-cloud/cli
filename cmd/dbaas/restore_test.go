package dbaas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBarmanRestoreTargetTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "valid positive offset",
			input: "2023-08-11 11:14:21.00000+02",
			want:  "2023-08-11 11:14:21.00000+02",
		},
		{
			name:  "valid negative offset",
			input: "2023-08-11 11:14:21.00000-05",
			want:  "2023-08-11 11:14:21.00000-05",
		},
		{
			name:  "trims whitespace",
			input: "  2023-08-11 11:14:21.00000+02  ",
			want:  "2023-08-11 11:14:21.00000+02",
		},
		{
			name:    "invalid format",
			input:   "2026-07-15T10:00:00Z",
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBarmanRestoreTargetTime(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
