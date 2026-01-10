package e2e

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNatGatewaysList(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "natgateways", "list")
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
		strings.Contains(outputLower, "region") ||
		strings.Contains(outputLower, "ip") ||
		strings.Contains(outputLower, "age")

	// If there's output, it should have headers (unless --no-header is used)
	if len(strings.TrimSpace(output)) > 0 {
		assert.True(t, hasExpectedHeader, "Output should contain at least one expected column header")
	}
}

func TestNatGatewaysListNoHeader(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "natgateways", "list", "--no-header")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --no-header, output should not contain header row
	// Check that the first line looks like data (not starting with "ID")
	output := result.Stdout
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Header row would start with "ID" (uppercase), data rows start with NAT gateway IDs
		if len(firstLine) > 0 {
			assert.True(t, !strings.HasPrefix(strings.ToUpper(firstLine), "ID"), "First line should be data, not header")
		}
	}
}

func TestNatGatewaysListShowLabels(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "natgateways", "list", "--show-labels")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --show-labels, output should contain a Labels column
	output := result.Stdout
	outputLower := strings.ToLower(output)
	assert.Contains(t, outputLower, "label", "Output should contain 'Label' column when --show-labels is used")
}

func TestNatGatewaysListWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Test with a label selector (this may return empty results, which is fine)
	result := config.RunCommand(t, "networking", "natgateways", "list", "--selector", "test=non-existent-label")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Command should succeed even if no NAT gateways match the selector
	// The output might be empty or just headers
}

func TestNatGatewaysListWithVPCFilter(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Get a VPC to filter by
	vpc := config.GetVPC(t)

	// Test with a VPC filter
	result := config.RunCommand(t, "networking", "natgateways", "list", "--vpc", vpc)
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Command should succeed even if no NAT gateways match the VPC
}

func TestNatGatewaysListAliases(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	aliases := []string{"get", "g", "ls"}

	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			result := config.RunCommand(t, "networking", "natgateways", alias)
			result.PrintOutput(t)
			result.AssertSuccess(t)
		})
	}
}

func TestNatGatewaysDeleteNonExistent(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Try to delete a non-existent NAT gateway
	result := config.RunCommand(t, "networking", "natgateways", "delete", "non-existent-ngw-12345", "--force")
	result.PrintOutput(t)

	// This should fail or report that the NAT gateway was not found
	// The exact behavior depends on the API, but it shouldn't succeed silently
	if result.ExitCode == 0 {
		// If it succeeds, it should at least mention that nothing was found
		assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "not found", "Should indicate NAT gateway was not found")
	}
}

func TestNatGatewaysDeleteWithoutArgs(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete command without arguments and without --selector should fail
	result := config.RunCommand(t, "networking", "natgateways", "delete", "--force")
	result.PrintOutput(t)

	result.AssertFailure(t)
	result.AssertStderrContains(t, "either NAT gateway identity")
}

func TestNatGatewaysDeleteWithSelector(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	// Delete with a selector that matches no NAT gateways
	result := config.RunCommand(t, "networking", "natgateways", "delete",
		"--selector", "e2e-test=non-existent",
		"--force")
	result.PrintOutput(t)

	// Should succeed but report no NAT gateways found
	result.AssertSuccess(t)
	assert.Contains(t, strings.ToLower(result.Stdout+result.Stderr), "no nat gateways", "Should report no NAT gateways found")
}

func TestNatGatewaysListOutputStructure(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "networking", "natgateways", "list")
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

