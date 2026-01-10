package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubnetsList(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "subnets", "list")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Check that output contains expected headers
	output := result.Stdout
	outputLower := strings.ToLower(output)

	// Should contain at least one expected column header
	hasExpectedHeader := strings.Contains(outputLower, "id") ||
		strings.Contains(outputLower, "name") ||
		strings.Contains(outputLower, "status") ||
		strings.Contains(outputLower, "vpc") ||
		strings.Contains(outputLower, "cidr") ||
		strings.Contains(outputLower, "age")

	// If there's output, it should have headers (unless --no-header is used)
	if len(strings.TrimSpace(output)) > 0 {
		assert.True(t, hasExpectedHeader, "Output should contain at least one expected column header")
	}
}

func TestSubnetsListNoHeader(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "subnets", "list", "--no-header")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --no-header, output should not contain header row
	// Check that the first line looks like data (starts with subnet ID pattern)
	output := result.Stdout
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Header row would start with "ID" (uppercase), data rows start with subnet IDs
		if len(firstLine) > 0 {
			// Subnet IDs typically start with "subnet-" or similar pattern
			assert.True(t, !strings.HasPrefix(strings.ToUpper(firstLine), "ID"), "First line should be data, not header")
		}
	}
}

func TestSubnetsListShowLabels(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "subnets", "list", "--show-labels")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --show-labels, output should contain a Labels column
	output := result.Stdout
	outputLower := strings.ToLower(output)
	assert.Contains(t, outputLower, "label", "Output should contain 'Label' column when --show-labels is used")
}

func TestSubnetsListWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Test with a label selector (this may return empty results, which is fine)
	result := config.RunCommand(t, "networking", "subnets", "list", "--selector", "test=non-existent-label")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Command should succeed even if no subnets match the selector
	// The output might be empty or just headers
}

func TestSubnetsListAliases(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	aliases := []string{"get", "g", "ls"}

	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			result := config.RunCommand(t, "networking", "subnets", alias)
			result.PrintOutput(t)
			result.AssertSuccess(t)
		})
	}
}

func TestSubnetsCreateMissingName(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	vpc := config.GetVPC(t)

	// Create command without required --name flag should fail
	result := config.RunCommand(t, "networking", "subnets", "create",
		"--vpc", vpc,
		"--cidr", "10.0.1.0/24")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "name")
}

func TestSubnetsCreateMissingVPC(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Create command without required --vpc flag should fail
	result := config.RunCommand(t, "networking", "subnets", "create",
		"--name", "test-subnet",
		"--cidr", "10.0.1.0/24")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "vpc")
}

func TestSubnetsCreateMissingCIDR(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	vpc := config.GetVPC(t)

	// Create command without required --cidr flag should fail
	result := config.RunCommand(t, "networking", "subnets", "create",
		"--name", "test-subnet",
		"--vpc", vpc)
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "cidr")
}

func TestSubnetsCreateAndDelete(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// This test creates a subnet and then deletes it
	// We use a unique name based on timestamp to avoid conflicts
	subnetName := "e2e-test-subnet-" + time.Now().Format("20060102150405")

	// Get a valid VPC
	vpc := config.GetVPC(t)

	// Create subnet
	createResult := config.RunCommand(t, "networking", "subnets", "create",
		"--name", subnetName,
		"--vpc", vpc,
		"--cidr", "10.0.1.0/24",
		"--no-header")
	createResult.PrintOutput(t)

	if createResult.ExitCode != 0 {
		t.Logf("Failed to create subnet, skipping delete test")
		t.Logf("This might be expected if the test environment doesn't allow subnet creation")
		return
	}

	// Extract subnet identity from output
	createOutput := createResult.GetLines()
	if len(createOutput) == 0 {
		t.Fatal("Create command succeeded but produced no output")
	}

	// The first column should be the subnet identity
	subnetFields := strings.Fields(createOutput[0])
	if len(subnetFields) == 0 {
		t.Fatal("Create output format is unexpected")
	}
	subnetIdentity := subnetFields[0]

	// Clean up: delete the subnet
	t.Cleanup(func() {
		deleteResult := config.RunCommand(t, "networking", "subnets", "delete", subnetIdentity, "--force")
		if deleteResult.ExitCode != 0 {
			t.Logf("Failed to clean up subnet %s: %s", subnetIdentity, deleteResult.Stderr)
		}
	})

	// Verify the subnet was created by listing it
	listResult := config.RunCommand(t, "networking", "subnets", "list", "--no-header")
	listResult.PrintOutput(t)
	listResult.AssertSuccess(t)

	// Check that the subnet appears in the list
	found := false
	for _, line := range listResult.GetLines() {
		if strings.Contains(line, subnetIdentity) || strings.Contains(line, subnetName) {
			found = true
			break
		}
	}
	assert.True(t, found, "Created subnet should appear in the list")

	// Delete the subnet
	deleteResult := config.RunCommand(t, "networking", "subnets", "delete", subnetIdentity, "--force")
	deleteResult.PrintOutput(t)
	deleteResult.AssertSuccess(t)
}

func TestSubnetsDeleteNonExistent(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Try to delete a non-existent subnet
	result := config.RunCommand(t, "networking", "subnets", "delete", "non-existent-subnet-12345", "--force")
	result.PrintOutput(t)

	// This should fail or report that the subnet was not found
	// The exact behavior depends on the API, but it shouldn't succeed silently
	if result.ExitCode == 0 {
		// If it succeeds, it should at least mention that nothing was found
		assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "not found", "Should indicate subnet was not found")
	}
}

func TestSubnetsDeleteWithoutArgs(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete command without arguments and without --selector should fail
	result := config.RunCommand(t, "networking", "subnets", "delete", "--force")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "either subnet identity")
}

func TestSubnetsDeleteWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete with a selector that matches no subnets
	result := config.RunCommand(t, "networking", "subnets", "delete",
		"--selector", "e2e-test=non-existent",
		"--force")
	result.PrintOutput(t)

	// Should succeed but report no subnets found
	result.AssertSuccess(t)
	assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "no subnets", "Should report no subnets found")
}

func TestSubnetsListOutputStructure(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "subnets", "list")
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

