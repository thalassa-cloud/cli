package dbaas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRestoreTargetTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "rfc3339 zulu",
			input: "2023-12-25T10:00:00Z",
			want:  "2023-12-25T10:00:00Z",
		},
		{
			name:  "barman positive offset converts to rfc3339 utc",
			input: "2023-08-11 11:14:21.00000+02",
			want:  "2023-08-11T09:14:21Z",
		},
		{
			name:  "barman negative offset converts to rfc3339 utc",
			input: "2023-08-11 11:14:21.00000-05",
			want:  "2023-08-11T16:14:21Z",
		},
		{
			name:  "trims whitespace",
			input: "  2023-08-11 11:14:21.00000+02  ",
			want:  "2023-08-11T09:14:21Z",
		},
		{
			name:    "invalid format",
			input:   "not-a-timestamp",
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
			got, err := parseRestoreTargetTime(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
