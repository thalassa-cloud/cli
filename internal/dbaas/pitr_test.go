package dbaasutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/dbaas"
)

func TestParseRestoreTargetTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantUTC string
		wantErr bool
	}{
		{
			name:    "rfc3339 zulu",
			input:   "2023-12-25T10:00:00Z",
			wantUTC: "2023-12-25T10:00:00Z",
		},
		{
			name:    "rfc3339 offset",
			input:   "2023-08-11T11:14:21+02:00",
			wantUTC: "2023-08-11T09:14:21Z",
		},
		{
			name:    "barman positive offset",
			input:   "2023-08-11 11:14:21.00000+02",
			wantUTC: "2023-08-11T09:14:21Z",
		},
		{
			name:    "barman negative offset",
			input:   "2023-08-11 11:14:21.00000-05",
			wantUTC: "2023-08-11T16:14:21Z",
		},
		{
			name:    "trims whitespace",
			input:   "  2023-12-25T10:00:00Z  ",
			wantUTC: "2023-12-25T10:00:00Z",
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
			got, err := ParseRestoreTargetTime(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantUTC, FormatRestoreTargetTimeRFC3339(got))
		})
	}
}

func TestValidateTargetTime(t *testing.T) {
	from := time.Date(2023, 8, 11, 9, 0, 0, 0, time.UTC)
	to := time.Date(2023, 8, 11, 12, 0, 0, 0, time.UTC)
	window := PITRWindow{From: &from, To: &to}

	tests := []struct {
		name    string
		target  time.Time
		window  PITRWindow
		wantErr bool
	}{
		{name: "inside window", target: time.Date(2023, 8, 11, 10, 0, 0, 0, time.UTC), window: window},
		{name: "at lower bound", target: from, window: window},
		{name: "at upper bound", target: to, window: window},
		{name: "before window", target: from.Add(-time.Second), window: window, wantErr: true},
		{name: "after window", target: to.Add(time.Second), window: window, wantErr: true},
		{name: "open window", target: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC), window: PITRWindow{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTargetTime(tt.target, tt.window)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestChainStatusNeedsWarning(t *testing.T) {
	tests := []struct {
		status dbaas.DbObjectStoreWalChainStatus
		want   bool
	}{
		{status: dbaas.DbObjectStoreWalChainStatusContinuous, want: false},
		{status: dbaas.DbObjectStoreWalChainStatusUnknown, want: false},
		{status: dbaas.DbObjectStoreWalChainStatusNotApplicable, want: false},
		{status: dbaas.DbObjectStoreWalChainStatusGapDetected, want: true},
		{status: dbaas.DbObjectStoreWalChainStatusEndWalMissing, want: true},
		{status: dbaas.DbObjectStoreWalChainStatusNoWal, want: true},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.want, ChainStatusNeedsWarning(tt.status))
		})
	}
}

func TestChainStatusWarning(t *testing.T) {
	missing := "000000010000000000000001"
	assert.Empty(t, ChainStatusWarning(dbaas.DbObjectStoreWalChainStatusContinuous, nil))
	got := ChainStatusWarning(dbaas.DbObjectStoreWalChainStatusGapDetected, &missing)
	assert.Contains(t, got, "gap_detected")
	assert.Contains(t, got, missing)
}

func TestBackupCoverageIsRestorable(t *testing.T) {
	tests := []struct {
		status string
		want   bool
	}{
		{status: "ready", want: true},
		{status: "completed", want: true},
		{status: "failed", want: false},
		{status: "deleted", want: false},
		{status: "unavailable", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			assert.Equal(t, tt.want, BackupCoverageIsRestorable(dbaas.DbObjectStoreRecoveryBackupCoverage{Status: tt.status}))
		})
	}
}

func TestWindowFromBackupCoverage(t *testing.T) {
	stopped := time.Date(2023, 8, 11, 10, 0, 0, 0, time.UTC)
	through := time.Date(2023, 8, 11, 11, 0, 0, 0, time.UTC)
	backup := dbaas.DbObjectStoreRecoveryBackupCoverage{
		CreatedAt: time.Date(2023, 8, 11, 9, 0, 0, 0, time.UTC),
		StoppedAt: &stopped,
		Coverage: dbaas.DbObjectStoreRecoveryBackupCoverageDetail{
			PitrThroughApproxAt: &through,
		},
	}

	window := WindowFromBackupCoverage(backup)
	require.NotNil(t, window.From)
	require.NotNil(t, window.To)
	assert.True(t, window.From.Equal(stopped))
	assert.True(t, window.To.Equal(through))
}
