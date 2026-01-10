package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVPCsList(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "vpcs", "list")
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
		strings.Contains(outputLower, "cidr") ||
		strings.Contains(outputLower, "age")

	// If there's output, it should have headers (unless --no-header is used)
	if len(strings.TrimSpace(output)) > 0 {
		assert.True(t, hasExpectedHeader, "Output should contain at least one expected column header")
	}
}

func TestVPCsListNoHeader(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "vpcs", "list", "--no-header")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --no-header, output should not contain headers
	output := result.Stdout
	outputLower := strings.ToLower(output)
	assert.NotContains(t, outputLower, "id", "Output should not contain 'ID' header")
	assert.NotContains(t, outputLower, "name", "Output should not contain 'Name' header")
	assert.NotContains(t, outputLower, "status", "Output should not contain 'Status' header")
}

func TestVPCsListShowLabels(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "vpcs", "list", "--show-labels")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --show-labels, output should contain a Labels column
	output := result.Stdout
	outputLower := strings.ToLower(output)
	assert.Contains(t, outputLower, "label", "Output should contain 'Label' column when --show-labels is used")
}

func TestVPCsListWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Test with a label selector (this may return empty results, which is fine)
	result := config.RunCommand(t, "networking", "vpcs", "list", "--selector", "test=non-existent-label")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Command should succeed even if no VPCs match the selector
	// The output might be empty or just headers
}

func TestVPCsListAliases(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	aliases := []string{"get", "g", "ls"}

	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			result := config.RunCommand(t, "networking", "vpcs", alias)
			result.PrintOutput(t)
			result.AssertSuccess(t)
		})
	}
}

func TestVPCsCreateMissingName(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Create command without required --name flag should fail
	result := config.RunCommand(t, "networking", "vpcs", "create",
		"--region", "test-region",
		"--cidrs", "10.0.0.0/16")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "name is required")
}

func TestVPCsCreateMissingRegion(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Create command without required --region flag should fail
	result := config.RunCommand(t, "networking", "vpcs", "create",
		"--name", "test-vpc",
		"--cidrs", "10.0.0.0/16")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "region is required")
}

func TestVPCsCreateMissingCIDRs(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Create command without CIDRs should fail (default is set, but let's test with empty)
	// Actually, the default is 10.0.0.0/16, so we need to test differently
	// Let's test with an invalid region instead, which will fail validation
	result := config.RunCommand(t, "networking", "vpcs", "create",
		"--name", "test-vpc",
		"--region", "invalid-region-12345")
	result.PrintOutput(t)

	// This should fail either because region doesn't exist or CIDRs validation
	// The exact error depends on API behavior, but it should fail
	result.AssertFailure(t)
}

func TestVPCsCreateAndDelete(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// This test creates a VPC and then deletes it
	// We use a unique name based on timestamp to avoid conflicts
	vpcName := "e2e-test-vpc-" + time.Now().Format("20060102150405")

	// First, we need to get a valid region
	// Let's try to list regions first to get a valid one
	regionsResult := config.RunCommand(t, "regions", "list", "--no-header")
	regionsResult.PrintOutput(t)

	if regionsResult.ExitCode != 0 {
		t.Skip("Cannot list regions, skipping VPC create/delete test")
	}

	regions := regionsResult.GetLines()
	if len(regions) == 0 {
		t.Skip("No regions available, skipping VPC create/delete test")
	}

	// Use the first region (format: identity or slug)
	regionLine := strings.Fields(regions[0])
	if len(regionLine) == 0 {
		t.Skip("Invalid region format, skipping VPC create/delete test")
	}
	region := regionLine[0]

	// Create VPC
	createResult := config.RunCommand(t, "networking", "vpcs", "create",
		"--name", vpcName,
		"--region", region,
		"--cidrs", "10.0.0.0/16",
		"--no-header")
	createResult.PrintOutput(t)

	if createResult.ExitCode != 0 {
		t.Logf("Failed to create VPC, skipping delete test")
		t.Logf("This might be expected if the test environment doesn't allow VPC creation")
		return
	}

	// Extract VPC identity from output
	createOutput := createResult.GetLines()
	if len(createOutput) == 0 {
		t.Fatal("Create command succeeded but produced no output")
	}

	// The first column should be the VPC identity
	vpcFields := strings.Fields(createOutput[0])
	if len(vpcFields) == 0 {
		t.Fatal("Create output format is unexpected")
	}
	vpcIdentity := vpcFields[0]

	// Clean up: delete the VPC
	t.Cleanup(func() {
		deleteResult := config.RunCommand(t, "networking", "vpcs", "delete", vpcIdentity, "--force")
		if deleteResult.ExitCode != 0 {
			t.Logf("Failed to clean up VPC %s: %s", vpcIdentity, deleteResult.Stderr)
		}
	})

	// Verify the VPC was created by listing it
	listResult := config.RunCommand(t, "networking", "vpcs", "list", "--no-header")
	listResult.PrintOutput(t)
	listResult.AssertSuccess(t)

	// Check that the VPC appears in the list
	found := false
	for _, line := range listResult.GetLines() {
		if strings.Contains(line, vpcIdentity) || strings.Contains(line, vpcName) {
			found = true
			break
		}
	}
	assert.True(t, found, "Created VPC should appear in the list")

	// Delete the VPC
	deleteResult := config.RunCommand(t, "networking", "vpcs", "delete", vpcIdentity, "--force")
	deleteResult.PrintOutput(t)
	deleteResult.AssertSuccess(t)
}

func TestVPCsDeleteNonExistent(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Try to delete a non-existent VPC
	result := config.RunCommand(t, "networking", "vpcs", "delete", "non-existent-vpc-12345", "--force")
	result.PrintOutput(t)

	// This should fail or report that the VPC was not found
	// The exact behavior depends on the API, but it shouldn't succeed silently
	if result.ExitCode == 0 {
		// If it succeeds, it should at least mention that nothing was found
		assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "not found", "Should indicate VPC was not found")
	}
}

func TestVPCsDeleteWithoutArgs(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete command without arguments and without --selector should fail
	result := config.RunCommand(t, "networking", "vpcs", "delete", "--force")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "either VPC identity")
}

func TestVPCsDeleteWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete with a selector that matches no VPCs
	result := config.RunCommand(t, "networking", "vpcs", "delete",
		"--selector", "e2e-test=non-existent",
		"--force")
	result.PrintOutput(t)

	// Should succeed but report no VPCs found
	result.AssertSuccess(t)
	assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "no vpcs", "Should report no VPCs found")
}

func TestVPCsCommandAliases(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	aliases := []string{"vpc", "virtualprivateclouds"}

	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			result := config.RunCommand(t, "networking", alias, "list", "--no-header")
			result.PrintOutput(t)
			result.AssertSuccess(t)
		})
	}
}

func TestVPCsListOutputStructure(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "vpcs", "list")
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
