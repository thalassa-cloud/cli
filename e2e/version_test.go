package e2e

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	config := LoadTestConfig(t)
	// Version command doesn't require authentication, so we don't skip if not configured

	result := config.RunCommand(t, "version")
	result.PrintOutput(t)

	// Version command should always succeed
	result.AssertSuccess(t)

	// Check that version information is printed
	assert.Contains(t, result.Stdout, "Version", "Version output should contain 'Version'")
	assert.Contains(t, result.Stdout, "Commit", "Version output should contain 'Commit'")
}

func TestVersionAlias(t *testing.T) {
	config := LoadTestConfig(t)

	result := config.RunCommand(t, "v")
	result.PrintOutput(t)

	result.AssertSuccess(t)
	assert.Contains(t, result.Stdout, "Version", "Version output should contain 'Version'")
}

func TestVersionOutputFormat(t *testing.T) {
	config := LoadTestConfig(t)

	result := config.RunCommand(t, "version")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Check that the output has the expected format
	lines := result.GetLines()
	require.Greater(t, len(lines), 0, "Version output should have at least one line")

	// Check for key fields in the output
	output := strings.ToLower(result.Stdout)
	assert.Contains(t, output, "version", "Output should contain version information")
	assert.Contains(t, output, "commit", "Output should contain commit information")
}
