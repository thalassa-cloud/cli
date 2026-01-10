package e2e

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeOrganisations(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "me", "organisations")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Check that output contains expected headers (unless --no-header is used)
	output := result.Stdout
	assert.NotEmpty(t, output, "Organisations output should not be empty")

	// The output should contain at least one of the expected columns
	hasID := strings.Contains(output, "ID") || strings.Contains(output, "Identity")
	hasName := strings.Contains(output, "Name")
	hasSlug := strings.Contains(output, "Slug")

	// At least one header should be present (unless there are no organisations)
	if len(strings.TrimSpace(output)) > 0 {
		assert.True(t, hasID || hasName || hasSlug, "Output should contain at least one expected column header")
	}
}

func TestMeOrganisationsNoHeader(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "me", "organisations", "--no-header")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --no-header, output should not contain headers
	output := result.Stdout
	assert.NotContains(t, output, "ID", "Output should not contain 'ID' header when --no-header is used")
	assert.NotContains(t, output, "Name", "Output should not contain 'Name' header when --no-header is used")
	assert.NotContains(t, output, "Slug", "Output should not contain 'Slug' header when --no-header is used")
}

func TestMeOrganisationsSlugOnly(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "me", "organisations", "--slug-only")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// With --slug-only, output should only contain slugs
	lines := result.GetLines()

	// If there are organisations, each line should be a slug (no spaces, no headers)
	if len(lines) > 0 {
		// First line might be a header, but with --slug-only it shouldn't be
		// Check that output doesn't contain column separators that would indicate multiple columns
		for _, line := range lines {
			// Slug-only output should be simple - just the slug values
			// No need to check for specific format, just that it's not the full table format
			assert.NotContains(t, strings.ToLower(line), "id", "Slug-only output should not contain 'ID'")
			assert.NotContains(t, strings.ToLower(line), "name", "Slug-only output should not contain 'Name'")
		}
	}
}

func TestMeOrganisationsSlugOnlyNoHeader(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "me", "organisations", "--slug-only", "--no-header")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	// Should only contain slug values, no headers
	output := result.Stdout
	assert.NotContains(t, output, "ID", "Output should not contain 'ID'")
	assert.NotContains(t, output, "Name", "Output should not contain 'Name'")
	assert.NotContains(t, output, "Slug", "Output should not contain 'Slug' header")
}

func TestMeOrganisationsOutputStructure(t *testing.T) {
	config := LoadTestConfig(t)
	config.SkipIfNotConfigured(t)

	result := config.RunCommand(t, "me", "organisations")
	result.PrintOutput(t)

	result.AssertSuccess(t)

	lines := result.GetLines()
	require.Greater(t, len(lines), 0, "Should have at least one line of output (header or data)")

	// If we have more than one line, the first should be headers and the rest should be data
	if len(lines) > 1 {
		headerLine := lines[0]
		// Header should contain expected column names
		headerLower := strings.ToLower(headerLine)
		hasExpectedColumn := strings.Contains(headerLower, "id") ||
			strings.Contains(headerLower, "name") ||
			strings.Contains(headerLower, "slug")
		assert.True(t, hasExpectedColumn, "Header line should contain expected column names")
	}
}
