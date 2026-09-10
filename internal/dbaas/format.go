package dbaasutil

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/client-go/dbaas"
)

const (
	goodCheckMark = "✔"
	kibibyte      = 1024
	mebibyte      = 1024 * kibibyte
	gibibyte      = 1024 * mebibyte
	tebibyte      = 1024 * gibibyte
)

// FormatBytes returns a human-readable binary size, or "-" when n is nil.
func FormatBytes(n *int64) string {
	if n == nil {
		return "-"
	}
	return FormatByteCount(*n)
}

// FormatByteCount returns a human-readable binary size for a non-pointer value.
func FormatByteCount(n int64) string {
	if n < 0 {
		return "-"
	}
	switch {
	case n >= tebibyte:
		return fmt.Sprintf("%.1f TiB", float64(n)/float64(tebibyte))
	case n >= gibibyte:
		return fmt.Sprintf("%.1f GiB", float64(n)/float64(gibibyte))
	case n >= mebibyte:
		return fmt.Sprintf("%.1f MiB", float64(n)/float64(mebibyte))
	case n >= kibibyte:
		return fmt.Sprintf("%.1f KiB", float64(n)/float64(kibibyte))
	default:
		return fmt.Sprintf("%d B", n)
	}
}

// FormatOptionalTime formats t, or returns "-" when t is nil.
func FormatOptionalTime(t *time.Time, showExactTime bool) string {
	if t == nil {
		return "-"
	}
	return formattime.FormatTime(t.Local(), showExactTime)
}

// FormatRetentionMode returns a stable display value for a backup store retention mode.
// Empty values default to retainForPointInTime.
func FormatRetentionMode(mode dbaas.DbObjectStoreRetentionMode) string {
	if mode == "" {
		return string(dbaas.DbObjectStoreRetentionModeRetainForPointInTime)
	}
	return string(mode)
}

// ParseRetentionMode validates and returns a known backup store retention mode.
func ParseRetentionMode(value string) (dbaas.DbObjectStoreRetentionMode, error) {
	mode := dbaas.DbObjectStoreRetentionMode(value)
	if value == "" {
		return dbaas.DbObjectStoreRetentionModeRetainForPointInTime, nil
	}
	if !mode.IsValid() {
		return "", fmt.Errorf("invalid retention mode %q: must be %s or %s",
			value,
			dbaas.DbObjectStoreRetentionModeRetainForPointInTime,
			dbaas.DbObjectStoreRetentionModeForceCleanupAfterExpiry,
		)
	}
	return mode, nil
}

// FormatBoolYesNo returns "yes" or "no".
func FormatBoolYesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

// IsGoodStatus reports whether value is a healthy/ready/good status label.
func IsGoodStatus(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "ready", "healthy", "ok", "okay", "good", "available", "continuous", "complete":
		return true
	default:
		return false
	}
}

// FormatStatus prefixes ready/healthy/good values with a green check mark.
func FormatStatus(value string) string {
	if IsGoodStatus(value) {
		return prefixGoodCheck(value)
	}
	return value
}

// FormatGoodYesNo returns a green-checked "yes" when ok, otherwise "no".
func FormatGoodYesNo(ok bool) string {
	if ok {
		return prefixGoodCheck("yes")
	}
	return "no"
}

func prefixGoodCheck(value string) string {
	check := goodCheckMark
	if isatty.IsTerminal(os.Stdout.Fd()) {
		check = table.Green(goodCheckMark)
	}
	return check + " " + value
}
