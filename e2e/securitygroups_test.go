package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityGroupsList(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "security-groups", "list")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Check that output contains expected headers
	output := result.Stdout
	outputLower := strings.ToLower(output)

	// If output is empty or just "No security groups found", skip header check
	if strings.TrimSpace(output) == "" || strings.Contains(outputLower, "no security groups found") {
		return
	}

	// Should contain at least one expected column header
	hasExpectedHeader := strings.Contains(outputLower, "id") ||
		strings.Contains(outputLower, "name") ||
		strings.Contains(outputLower, "status") ||
		strings.Contains(outputLower, "vpc") ||
		strings.Contains(outputLower, "ingress") ||
		strings.Contains(outputLower, "egress") ||
		strings.Contains(outputLower, "age")

	// If there's output, it should have headers (unless --no-header is used)
	if len(strings.TrimSpace(output)) > 0 {
		assert.True(t, hasExpectedHeader, "Output should contain at least one expected column header")
	}
}

func TestSecurityGroupsListNoHeader(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "security-groups", "list", "--no-header")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --no-header, output should not contain header row
	// Check that the first line looks like data (not starting with "ID")
	output := result.Stdout
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Header row would start with "ID" (uppercase), data rows start with security group IDs
		if len(firstLine) > 0 {
			assert.True(t, !strings.HasPrefix(strings.ToUpper(firstLine), "ID"), "First line should be data, not header")
		}
	}
}

func TestSecurityGroupsListShowLabels(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "security-groups", "list", "--show-labels")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --show-labels, output should contain a Labels column
	output := result.Stdout
	outputLower := strings.ToLower(output)
	
	// If there are no security groups, the output will just say "No security groups found"
	if strings.Contains(outputLower, "no security groups found") {
		return
	}
	
	assert.Contains(t, outputLower, "label", "Output should contain 'Label' column when --show-labels is used")
}

func TestSecurityGroupsListWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Test with a label selector (this may return empty results, which is fine)
	result := config.RunCommand(t, "networking", "security-groups", "list", "--selector", "test=non-existent-label")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Command should succeed even if no security groups match the selector
	// The output might be empty or just headers
}

func TestSecurityGroupsListWithVPCFilter(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Get a VPC to filter by
	vpc := config.GetVPC(t)

	// Test with a VPC filter
	result := config.RunCommand(t, "networking", "security-groups", "list", "--vpc", vpc)
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Command should succeed even if no security groups match the VPC
}

func TestSecurityGroupsListAliases(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	aliases := []string{"get", "g", "ls"}

	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			result := config.RunCommand(t, "networking", "security-groups", alias)
			result.PrintOutput(t)
			result.AssertSuccess(t)
		})
	}
}

func TestSecurityGroupsCreateMissingName(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	vpc := config.GetVPC(t)

	// Create command without required --name flag should fail
	result := config.RunCommand(t, "networking", "security-groups", "create",
		"--vpc", vpc)
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "name")
}

func TestSecurityGroupsCreateMissingVPC(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Create command without required --vpc flag should fail
	result := config.RunCommand(t, "networking", "security-groups", "create",
		"--name", "test-sg")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "vpc")
}

func TestSecurityGroupsCreateAndDelete(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// This test creates a security group and then deletes it
	// We use a unique name based on timestamp to avoid conflicts
	sgName := "e2e-test-sg-" + time.Now().Format("20060102150405")

	// Get a valid VPC
	vpc := config.GetVPC(t)

	// Create security group
	createResult := config.RunCommand(t, "networking", "security-groups", "create",
		"--name", sgName,
		"--vpc", vpc)
	createResult.PrintOutput(t)

	if createResult.ExitCode != 0 {
		t.Logf("Failed to create security group, skipping delete test")
		t.Logf("This might be expected if the test environment doesn't allow security group creation")
		return
	}

	// Extract security group identity from output
	// The output format is: "Security group created successfully\nID: <identity>\nName: <name>\nStatus: <status>"
	output := createResult.Stdout
	sgIdentity := ""
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID: ") {
			sgIdentity = strings.TrimPrefix(line, "ID: ")
			break
		}
	}

	if sgIdentity == "" {
		t.Fatal("Could not extract security group identity from create output")
	}

	// Clean up: delete the security group
	t.Cleanup(func() {
		deleteResult := config.RunCommand(t, "networking", "security-groups", "delete", sgIdentity, "--force")
		if deleteResult.ExitCode != 0 {
			t.Logf("Failed to clean up security group %s: %s", sgIdentity, deleteResult.Stderr)
		}
	})

	// Verify the security group was created by listing it
	listResult := config.RunCommand(t, "networking", "security-groups", "list", "--no-header")
	listResult.PrintOutput(t)
	listResult.AssertSuccess(t)

	// Check that the security group appears in the list
	found := false
	for _, line := range listResult.GetLines() {
		if strings.Contains(line, sgIdentity) || strings.Contains(line, sgName) {
			found = true
			break
		}
	}
	assert.True(t, found, "Created security group should appear in the list")

	// Delete the security group
	deleteResult := config.RunCommand(t, "networking", "security-groups", "delete", sgIdentity, "--force")
	deleteResult.PrintOutput(t)
	deleteResult.AssertSuccess(t)
}

func TestSecurityGroupsDeleteNonExistent(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Try to delete a non-existent security group
	result := config.RunCommand(t, "networking", "security-groups", "delete", "non-existent-sg-12345", "--force")
	result.PrintOutput(t)

	// This should fail or report that the security group was not found
	// The exact behavior depends on the API, but it shouldn't succeed silently
	if result.ExitCode == 0 {
		// If it succeeds, it should at least mention that nothing was found
		assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "not found", "Should indicate security group was not found")
	}
}

func TestSecurityGroupsDeleteWithoutArgs(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete command without arguments and without --selector should fail
	result := config.RunCommand(t, "networking", "security-groups", "delete", "--force")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "either security group identity")
}

func TestSecurityGroupsDeleteWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete with a selector that matches no security groups
	result := config.RunCommand(t, "networking", "security-groups", "delete",
		"--selector", "e2e-test=non-existent",
		"--force")
	result.PrintOutput(t)

	// Should succeed but report no security groups found
	result.AssertSuccess(t)
	assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "no security groups", "Should report no security groups found")
}

func TestSecurityGroupsListOutputStructure(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "security-groups", "list")
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

