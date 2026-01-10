package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotsList(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "storage", "snapshots", "list")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Check that output contains expected headers
	output := result.Stdout
	outputLower := strings.ToLower(output)

	// Should contain at least one expected column header
	hasExpectedHeader := strings.Contains(outputLower, "id") ||
		strings.Contains(outputLower, "name") ||
		strings.Contains(outputLower, "status") ||
		strings.Contains(outputLower, "region") ||
		strings.Contains(outputLower, "size") ||
		strings.Contains(outputLower, "age")

	// If there's output, it should have headers (unless --no-header is used)
	if len(strings.TrimSpace(output)) > 0 {
		assert.True(t, hasExpectedHeader, "Output should contain at least one expected column header")
	}
}

func TestSnapshotsListNoHeader(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "storage", "snapshots", "list", "--no-header")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --no-header, output should not contain header row
	// Check that the first line looks like data (starts with snapshot ID pattern "s-")
	output := result.Stdout
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Header row would start with "ID" (uppercase), data rows start with snapshot IDs
		if len(firstLine) > 0 {
			// Snapshot IDs typically start with "s-" or similar pattern
			assert.True(t, strings.HasPrefix(firstLine, "s-") || len(lines) == 0, "First line should be data (snapshot ID), not header")
		}
	}
}

func TestSnapshotsListShowLabels(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "storage", "snapshots", "list", "--show-labels")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --show-labels, output should contain a Labels column
	output := result.Stdout
	outputLower := strings.ToLower(output)
	assert.Contains(t, outputLower, "label", "Output should contain 'Label' column when --show-labels is used")
}

func TestSnapshotsListWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Test with a label selector (this may return empty results, which is fine)
	result := config.RunCommand(t, "storage", "snapshots", "list", "--selector", "test=non-existent-label")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Command should succeed even if no snapshots match the selector
	// The output might be empty or just headers
}

func TestSnapshotsListAliases(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	aliases := []string{"get", "g", "ls"}

	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			result := config.RunCommand(t, "storage", "snapshots", alias)
			result.PrintOutput(t)
			result.AssertSuccess(t)
		})
	}
}

func TestSnapshotsCreateMissingVolume(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Create command without required --volume flag should fail
	result := config.RunCommand(t, "storage", "snapshots", "create", "test-snapshot")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "volume")
}

func TestSnapshotsCreateAndDelete(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// This test creates a snapshot and then deletes it
	// We use a unique name based on timestamp to avoid conflicts
	snapshotName := "e2e-test-snapshot-" + time.Now().Format("20060102150405")

	// Get a valid volume
	volumeIdentity := config.GetVolume(t)

	// Create snapshot
	createResult := config.RunCommand(t, "storage", "snapshots", "create", snapshotName,
		"--volume", volumeIdentity)
	createResult.PrintOutput(t)

	if createResult.ExitCode != 0 {
		t.Logf("Failed to create snapshot, skipping delete test")
		t.Logf("This might be expected if the test environment doesn't allow snapshot creation")
		return
	}

	// Extract snapshot identity from output
	// The output format is: "Snapshot created successfully: <name> (<identity>)"
	output := createResult.Stdout
	snapshotIdentity := ""
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Snapshot created successfully") {
			// Extract identity from line like "Snapshot created successfully: test-snapshot (s-xxx)"
			parts := strings.Split(line, "(")
			if len(parts) > 1 {
				snapshotIdentity = strings.TrimSuffix(strings.TrimSpace(parts[1]), ")")
			}
			break
		}
	}

	if snapshotIdentity == "" {
		t.Fatal("Could not extract snapshot identity from create output")
	}

	// Clean up: delete the snapshot
	t.Cleanup(func() {
		deleteResult := config.RunCommand(t, "storage", "snapshots", "delete", snapshotIdentity, "--force")
		if deleteResult.ExitCode != 0 {
			t.Logf("Failed to clean up snapshot %s: %s", snapshotIdentity, deleteResult.Stderr)
		}
	})

	// Verify the snapshot was created by listing it
	listResult := config.RunCommand(t, "storage", "snapshots", "list", "--no-header")
	listResult.PrintOutput(t)
	listResult.AssertSuccess(t)

	// Check that the snapshot appears in the list
	found := false
	for _, line := range listResult.GetLines() {
		if strings.Contains(line, snapshotIdentity) || strings.Contains(line, snapshotName) {
			found = true
			break
		}
	}
	assert.True(t, found, "Created snapshot should appear in the list")

	// Delete the snapshot
	deleteResult := config.RunCommand(t, "storage", "snapshots", "delete", snapshotIdentity, "--force")
	deleteResult.PrintOutput(t)
	deleteResult.AssertSuccess(t)
}

func TestSnapshotsDeleteNonExistent(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Try to delete a non-existent snapshot
	result := config.RunCommand(t, "storage", "snapshots", "delete", "non-existent-snapshot-12345", "--force")
	result.PrintOutput(t)

	// This should fail or report that the snapshot was not found
	// The exact behavior depends on the API, but it shouldn't succeed silently
	if result.ExitCode == 0 {
		// If it succeeds, it should at least mention that nothing was found
		assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "not found", "Should indicate snapshot was not found")
	}
}

func TestSnapshotsDeleteWithoutArgs(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete command without arguments and without --selector should fail
	result := config.RunCommand(t, "storage", "snapshots", "delete", "--force")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "either snapshot identity")
}

func TestSnapshotsDeleteWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete with a selector that matches no snapshots
	result := config.RunCommand(t, "storage", "snapshots", "delete",
		"--selector", "e2e-test=non-existent",
		"--force")
	result.PrintOutput(t)

	// Should succeed but report no snapshots found
	result.AssertSuccess(t)
	assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "no snapshots", "Should report no snapshots found")
}

func TestSnapshotsListOutputStructure(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "storage", "snapshots", "list")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	lines := result.GetLines()
	require.Greater(t, len(lines), 0, "Should have at least one line of output")

	// If we have more than one line, the first should be headers
	if len(lines) > 1 {
		headerLine := strings.ToLower(lines[0])
		hasExpectedColumn := strings.Contains(headerLine, "id") ||
			strings.Contains(headerLine, "name") ||
			strings.Contains(headerLine, "status")
		assert.True(t, hasExpectedColumn, "Header line should contain expected column names")
	}
}

