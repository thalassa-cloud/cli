package dbaasutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/thalassa-cloud/client-go/dbaas"
)

const (
	barmanTargetTimeExample = "2023-08-11 11:14:21.00000+02"
	barmanTargetTimeLayout  = "2006-01-02 15:04:05.00000-07"
)

// RestoreTargetTimeDescription describes accepted restore target time formats.
const RestoreTargetTimeDescription = "Timestamp to restore to. Accepts RFC3339 (e.g. 2023-12-25T10:00:00Z) or barman format (YYYY-MM-DD HH:MM:SS.00000±TZ, e.g. '" + barmanTargetTimeExample + "')"

// PITRWindow is an inclusive recoverability window.
type PITRWindow struct {
	From *time.Time
	To   *time.Time
}

// ParseRestoreTargetTime accepts RFC3339, RFC3339Nano, or barman timestamps.
func ParseRestoreTargetTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("restore target time is empty")
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		barmanTargetTimeLayout,
	}
	var firstErr error
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}

	return time.Time{}, fmt.Errorf("parse restore target time: expected RFC3339 or barman format (%s): %w", barmanTargetTimeExample, firstErr)
}

// FormatRestoreTargetTimeRFC3339 formats t as UTC RFC3339 for the restore API.
func FormatRestoreTargetTimeRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// ValidateTargetTime reports whether target falls inside window (inclusive).
func ValidateTargetTime(target time.Time, window PITRWindow) error {
	if window.From != nil && target.Before(*window.From) {
		return fmt.Errorf("target time %s is before earliest recoverability point %s",
			FormatRestoreTargetTimeRFC3339(target),
			FormatRestoreTargetTimeRFC3339(*window.From),
		)
	}
	if window.To != nil && target.After(*window.To) {
		return fmt.Errorf("target time %s is after latest approximate PITR time %s",
			FormatRestoreTargetTimeRFC3339(target),
			FormatRestoreTargetTimeRFC3339(*window.To),
		)
	}
	return nil
}

// WindowFromStoreSummary builds a PITR window from recovery overview summary fields.
func WindowFromStoreSummary(summary dbaas.DbObjectStoreRecoveryStoreSummary) PITRWindow {
	return PITRWindow{
		From: summary.EarliestPitrFrom,
		To:   summary.LatestPitrThroughApprox,
	}
}

// WindowFromBackupCoverage builds a PITR window for a single backup's WAL coverage.
// The lower bound is the backup stop/start/create time; the upper bound is the
// approximate PITR-through timestamp from WAL analysis.
func WindowFromBackupCoverage(backup dbaas.DbObjectStoreRecoveryBackupCoverage) PITRWindow {
	from := backup.StoppedAt
	if from == nil {
		from = backup.StartedAt
	}
	if from == nil {
		created := backup.CreatedAt
		from = &created
	}
	return PITRWindow{
		From: from,
		To:   backup.Coverage.PitrThroughApproxAt,
	}
}

// ChainStatusNeedsWarning reports whether WAL coverage should be confirmed before PITR.
func ChainStatusNeedsWarning(status dbaas.DbObjectStoreWalChainStatus) bool {
	switch status {
	case dbaas.DbObjectStoreWalChainStatusGapDetected,
		dbaas.DbObjectStoreWalChainStatusEndWalMissing,
		dbaas.DbObjectStoreWalChainStatusNoWal:
		return true
	default:
		return false
	}
}

// ChainStatusWarning returns a human-readable warning for incomplete WAL coverage.
func ChainStatusWarning(status dbaas.DbObjectStoreWalChainStatus, firstMissingWAL *string) string {
	if !ChainStatusNeedsWarning(status) {
		return ""
	}
	msg := fmt.Sprintf("WAL chain status is %s; point-in-time recovery may fail or stop early", status)
	if firstMissingWAL != nil && *firstMissingWAL != "" {
		msg = fmt.Sprintf("%s (first missing WAL: %s)", msg, *firstMissingWAL)
	}
	return msg
}

// BackupCoverageIsRestorable reports whether a backup looks usable as a restore base.
func BackupCoverageIsRestorable(backup dbaas.DbObjectStoreRecoveryBackupCoverage) bool {
	status := strings.ToLower(strings.TrimSpace(backup.Status))
	switch status {
	case string(dbaas.ObjectStatusFailed), string(dbaas.ObjectStatusDeleted), "unavailable":
		return false
	default:
		return true
	}
}
