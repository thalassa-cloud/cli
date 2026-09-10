package dbaasutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/dbaas"
)

func TestFormatByteCount(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		want string
	}{
		{name: "bytes", n: 512, want: "512 B"},
		{name: "kibibytes", n: 2048, want: "2.0 KiB"},
		{name: "mebibytes", n: 1048576, want: "1.0 MiB"},
		{name: "gibibytes", n: 1073741824, want: "1.0 GiB"},
		{name: "negative", n: -1, want: "-"},
		{name: "zero", n: 0, want: "0 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatByteCount(tt.n))
		})
	}
}

func TestFormatBytes(t *testing.T) {
	assert.Equal(t, "-", FormatBytes(nil))
	n := int64(1024)
	assert.Equal(t, "1.0 KiB", FormatBytes(&n))
}

func TestFormatRetentionMode(t *testing.T) {
	tests := []struct {
		name string
		mode dbaas.DbObjectStoreRetentionMode
		want string
	}{
		{name: "empty defaults to retain", mode: "", want: string(dbaas.DbObjectStoreRetentionModeRetainForPointInTime)},
		{name: "retain", mode: dbaas.DbObjectStoreRetentionModeRetainForPointInTime, want: "retainForPointInTime"},
		{name: "force cleanup", mode: dbaas.DbObjectStoreRetentionModeForceCleanupAfterExpiry, want: "forceCleanupAfterExpiry"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatRetentionMode(tt.mode))
		})
	}
}

func TestParseRetentionMode(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    dbaas.DbObjectStoreRetentionMode
		wantErr bool
	}{
		{name: "empty defaults", value: "", want: dbaas.DbObjectStoreRetentionModeRetainForPointInTime},
		{name: "retain", value: "retainForPointInTime", want: dbaas.DbObjectStoreRetentionModeRetainForPointInTime},
		{name: "force cleanup", value: "forceCleanupAfterExpiry", want: dbaas.DbObjectStoreRetentionModeForceCleanupAfterExpiry},
		{name: "invalid", value: "deleteEverything", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRetentionMode(tt.value)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatBoolYesNo(t *testing.T) {
	assert.Equal(t, "yes", FormatBoolYesNo(true))
	assert.Equal(t, "no", FormatBoolYesNo(false))
}

func TestIsGoodStatus(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{value: "ready", want: true},
		{value: "Ready", want: true},
		{value: "healthy", want: true},
		{value: "ok", want: true},
		{value: "good", want: true},
		{value: "available", want: true},
		{value: "continuous", want: true},
		{value: "complete", want: true},
		{value: "failed", want: false},
		{value: "yes", want: false},
		{value: "creating", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.want, IsGoodStatus(tt.value))
		})
	}
}

func TestFormatStatus(t *testing.T) {
	assert.Equal(t, "failed", FormatStatus("failed"))
	got := FormatStatus("ready")
	assert.Contains(t, got, "ready")
	assert.Contains(t, got, "✔")
}

func TestFormatGoodYesNo(t *testing.T) {
	assert.Equal(t, "no", FormatGoodYesNo(false))
	got := FormatGoodYesNo(true)
	assert.Contains(t, got, "yes")
	assert.Contains(t, got, "✔")
}
