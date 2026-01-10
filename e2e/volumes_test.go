package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVolumesList(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "storage", "volumes", "list")
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

func TestVolumesListNoHeader(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "storage", "volumes", "list", "--no-header")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --no-header, output should not contain header row
	// Check that the first line looks like data (starts with volume ID pattern "v-")
	// rather than a header row
	output := result.Stdout
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Header row would start with "ID" (uppercase), data rows start with volume IDs like "v-"
		assert.True(t, strings.HasPrefix(firstLine, "v-"), "First line should be data (volume ID starting with 'v-'), not header")
		// Also verify it doesn't contain header words as column headers
		firstLineLower := strings.ToLower(firstLine)
		assert.NotContains(t, firstLineLower, "\tid\t", "First line should not contain 'ID' as a column header")
		assert.NotContains(t, firstLineLower, "\tname\t", "First line should not contain 'Name' as a column header")
	}
}

func TestVolumesListShowLabels(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "storage", "volumes", "list", "--show-labels")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --show-labels, output should contain a Labels column
	output := result.Stdout
	outputLower := strings.ToLower(output)
	assert.Contains(t, outputLower, "label", "Output should contain 'Label' column when --show-labels is used")
}

func TestVolumesListWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Test with a label selector (this may return empty results, which is fine)
	result := config.RunCommand(t, "storage", "volumes", "list", "--selector", "test=non-existent-label")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Command should succeed even if no volumes match the selector
	// The output might be empty or just headers
}

func TestVolumesListAliases(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	aliases := []string{"get", "g", "ls"}

	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			result := config.RunCommand(t, "storage", "volumes", alias)
			result.PrintOutput(t)
			result.AssertSuccess(t)
		})
	}
}

func TestVolumesCreateMissingName(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	region := config.GetRegion(t)

	// Create command without required --name flag should fail
	result := config.RunCommand(t, "storage", "volumes", "create",
		"--region", region,
		"--size", "10",
		"--type", "block")
	result.PrintOutput(t)

	result.AssertFailure(t)
	// Cobra returns "required flag(s) "name" not set" or custom error "name is required"
	result.AssertStderrContains(t, "name")
}

func TestVolumesCreateMissingRegion(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Create command without required --region flag should fail
	result := config.RunCommand(t, "storage", "volumes", "create",
		"--name", "test-volume",
		"--size", "10")
	result.PrintOutput(t)

	result.AssertFailure(t)
	// Cobra returns "required flag(s) "region" not set" or custom error "region is required"
	result.AssertStderrContains(t, "region")
}

func TestVolumesCreateMissingSize(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	region := config.GetRegion(t)

	// Create command without required --size flag should fail
	result := config.RunCommand(t, "storage", "volumes", "create",
		"--name", "test-volume",
		"--region", region,
		"--type", "block")
	result.PrintOutput(t)

	result.AssertFailure(t)
	// Cobra returns "required flag(s) "size" not set" or custom error "size is required"
	result.AssertStderrContains(t, "size")
}

func TestVolumesCreateInvalidSize(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	region := config.GetRegion(t)

	// Create command with invalid size should fail
	result := config.RunCommand(t, "storage", "volumes", "create",
		"--name", "test-volume",
		"--region", region,
		"--size", "0",
		"--type", "block")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "size must be greater than 0")
}

func TestVolumesCreateAndDelete(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// This test creates a volume and then deletes it
	// We use a unique name based on timestamp to avoid conflicts
	volumeName := "e2e-test-volume-" + time.Now().Format("20060102150405")

	// Get a valid region
	region := config.GetRegion(t)

	// Create volume with 'block' volume type
	createResult := config.RunCommand(t, "storage", "volumes", "create",
		"--name", volumeName,
		"--region", region,
		"--size", "10",
		"--type", "block",
		"--no-header")
	createResult.PrintOutput(t)

	if createResult.ExitCode != 0 {
		t.Logf("Failed to create volume, skipping delete test")
		t.Logf("This might be expected if the test environment doesn't allow volume creation")
		return
	}

	// Extract volume identity from output
	createOutput := createResult.GetLines()
	if len(createOutput) == 0 {
		t.Fatal("Create command succeeded but produced no output")
	}

	// The first column should be the volume identity
	volumeFields := strings.Fields(createOutput[0])
	if len(volumeFields) == 0 {
		t.Fatal("Create output format is unexpected")
	}
	volumeIdentity := volumeFields[0]

	// Clean up: delete the volume
	t.Cleanup(func() {
		deleteResult := config.RunCommand(t, "storage", "volumes", "delete", volumeIdentity, "--force")
		if deleteResult.ExitCode != 0 {
			t.Logf("Failed to clean up volume %s: %s", volumeIdentity, deleteResult.Stderr)
		}
	})

	// Verify the volume was created by listing it
	listResult := config.RunCommand(t, "storage", "volumes", "list", "--no-header")
	listResult.PrintOutput(t)
	listResult.AssertSuccess(t)

	// Check that the volume appears in the list
	found := false
	for _, line := range listResult.GetLines() {
		if strings.Contains(line, volumeIdentity) || strings.Contains(line, volumeName) {
			found = true
			break
		}
	}
	assert.True(t, found, "Created volume should appear in the list")

	// Delete the volume
	deleteResult := config.RunCommand(t, "storage", "volumes", "delete", volumeIdentity, "--force")
	deleteResult.PrintOutput(t)
	deleteResult.AssertSuccess(t)
}

func TestVolumesDeleteNonExistent(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Try to delete a non-existent volume
	result := config.RunCommand(t, "storage", "volumes", "delete", "non-existent-volume-12345", "--force")
	result.PrintOutput(t)

	// This should fail or report that the volume was not found
	// The exact behavior depends on the API, but it shouldn't succeed silently
	if result.ExitCode == 0 {
		// If it succeeds, it should at least mention that nothing was found
		assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "not found", "Should indicate volume was not found")
	}
}

func TestVolumesDeleteWithoutArgs(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete command without arguments and without --selector should fail
	result := config.RunCommand(t, "storage", "volumes", "delete", "--force")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "either volume identity")
}

func TestVolumesDeleteWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete with a selector that matches no volumes
	result := config.RunCommand(t, "storage", "volumes", "delete",
		"--selector", "e2e-test=non-existent",
		"--force")
	result.PrintOutput(t)

	// Should succeed but report no volumes found
	result.AssertSuccess(t)
	assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "no volumes", "Should report no volumes found")
}

func TestVolumesCommandAliases(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	aliases := []string{"volume", "vol"}

	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			result := config.RunCommand(t, "storage", alias, "list", "--no-header")
			result.PrintOutput(t)
			result.AssertSuccess(t)
		})
	}
}

func TestVolumesListOutputStructure(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "storage", "volumes", "list")
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
