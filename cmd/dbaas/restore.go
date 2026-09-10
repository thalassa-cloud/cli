package dbaas

import (
	dbaasutil "github.com/thalassa-cloud/cli/internal/dbaas"
)

const restoreTargetTimeDescription = dbaasutil.RestoreTargetTimeDescription

func parseRestoreTargetTime(value string) (string, error) {
	targetTime, err := dbaasutil.ParseRestoreTargetTime(value)
	if err != nil {
		return "", err
	}
	return dbaasutil.FormatRestoreTargetTimeRFC3339(targetTime), nil
}
