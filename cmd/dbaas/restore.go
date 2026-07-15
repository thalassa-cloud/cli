package dbaas

import (
	"fmt"
	"strings"
	"time"
)

const (
	barmanTargetTimeExample     = "2023-08-11 11:14:21.00000+02"
	barmanTargetTimeDescription = "Timestamp to restore to in barman format (YYYY-MM-DD HH:MM:SS.00000±TZ). Example: '" + barmanTargetTimeExample + "'"
	barmanTargetTimeLayout      = "2006-01-02 15:04:05.00000-07"
)

func parseBarmanRestoreTargetTime(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("restore target time is empty")
	}

	targetTime, err := time.Parse(barmanTargetTimeLayout, value)
	if err != nil {
		return "", fmt.Errorf("parse barman target time: %w", err)
	}

	return targetTime.Format(barmanTargetTimeLayout), nil
}
