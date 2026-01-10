package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestConfig holds configuration for E2E tests
type TestConfig struct {
	APIEndpoint         string
	AccessToken         string
	PersonalAccessToken string
	OIDCClientID        string
	OIDCClientSecret    string
	Organisation        string
	BinaryPath          string
}

// LoadTestConfig loads test configuration from environment variables
func LoadTestConfig(t *testing.T) *TestConfig {
	config := &TestConfig{
		APIEndpoint:         os.Getenv("TCLOUD_E2E_API_ENDPOINT"),
		AccessToken:         os.Getenv("TCLOUD_E2E_ACCESS_TOKEN"),
		PersonalAccessToken: os.Getenv("TCLOUD_E2E_PERSONAL_ACCESS_TOKEN"),
		OIDCClientID:        os.Getenv("TCLOUD_E2E_OIDC_CLIENT_ID"),
		OIDCClientSecret:    os.Getenv("TCLOUD_E2E_OIDC_CLIENT_SECRET"),
		Organisation:        os.Getenv("TCLOUD_E2E_ORGANISATION"),
		BinaryPath:          os.Getenv("TCLOUD_E2E_BINARY_PATH"),
	}

	// Default binary path to ./bin/tcloud if not set
	if config.BinaryPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}
		// Go up from e2e directory to project root
		projectRoot := filepath.Dir(wd)
		config.BinaryPath = filepath.Join(projectRoot, "bin", "tcloud")
	}

	return config
}

// SkipIfNotConfigured skips the test if required configuration is missing
func (c *TestConfig) SkipIfNotConfigured(t *testing.T) {
	hasAuth := c.AccessToken != "" || c.PersonalAccessToken != "" || (c.OIDCClientID != "" && c.OIDCClientSecret != "")
	if !hasAuth {
		t.Skip("Skipping E2E test: no authentication configured (set TCLOUD_E2E_ACCESS_TOKEN, TCLOUD_E2E_PERSONAL_ACCESS_TOKEN, or TCLOUD_E2E_OIDC_CLIENT_ID/SECRET)")
	}

	if c.APIEndpoint == "" {
		t.Skip("Skipping E2E test: TCLOUD_E2E_API_ENDPOINT not set")
	}

	// Check if binary exists
	if _, err := os.Stat(c.BinaryPath); os.IsNotExist(err) {
		t.Skipf("Skipping E2E test: binary not found at %s (build it first with 'make build')", c.BinaryPath)
	}
}

// CommandResult holds the result of a command execution
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

// RunCommand executes a CLI command and returns the result
func (c *TestConfig) RunCommand(t *testing.T, args ...string) *CommandResult {
	t.Helper()

	cmd := exec.Command(c.BinaryPath, args...)

	// Set up authentication flags
	if c.APIEndpoint != "" {
		cmd.Args = append(cmd.Args, "--api", c.APIEndpoint)
	}
	if c.AccessToken != "" {
		cmd.Args = append(cmd.Args, "--access-token", c.AccessToken)
	} else if c.PersonalAccessToken != "" {
		cmd.Args = append(cmd.Args, "--token", c.PersonalAccessToken)
	} else if c.OIDCClientID != "" && c.OIDCClientSecret != "" {
		cmd.Args = append(cmd.Args, "--client-id", c.OIDCClientID, "--client-secret", c.OIDCClientSecret)
	}
	if c.Organisation != "" {
		cmd.Args = append(cmd.Args, "--organisation", c.Organisation)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Print the command being executed (including all flags)
	// Print the command but avoid leaking secrets in logs
	maskedArgs := make([]string, len(cmd.Args))
	copy(maskedArgs, cmd.Args)
	secretFlags := map[string]bool{
		"--access-token":  true,
		"--token":         true,
		"--client-id":     true,
		"--client-secret": true,
	}
	for i := 0; i < len(maskedArgs); i++ {
		if secretFlags[maskedArgs[i]] {
			// Mask the argument that follows the secret flag
			if i+1 < len(maskedArgs) {
				maskedArgs[i+1] = "****"
			}
		}
	}
	t.Logf("Executing command: %s", strings.Join(maskedArgs, " "))

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	return &CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Error:    err,
	}
}

// AssertSuccess asserts that the command succeeded
func (r *CommandResult) AssertSuccess(t *testing.T) {
	t.Helper()
	if r.ExitCode != 0 {
		t.Errorf("Command failed with exit code %d\nStdout: %s\nStderr: %s", r.ExitCode, r.Stdout, r.Stderr)
	}
}

// AssertFailure asserts that the command failed
func (r *CommandResult) AssertFailure(t *testing.T) {
	t.Helper()
	if r.ExitCode == 0 {
		t.Errorf("Command succeeded but was expected to fail\nStdout: %s\nStderr: %s", r.Stdout, r.Stderr)
	}
}

// AssertStdoutContains asserts that stdout contains the given string
func (r *CommandResult) AssertStdoutContains(t *testing.T, substr string) {
	t.Helper()
	if !strings.Contains(r.Stdout, substr) {
		t.Errorf("Expected stdout to contain %q, but got:\n%s", substr, r.Stdout)
	}
}

// AssertStderrContains asserts that stderr contains the given string
func (r *CommandResult) AssertStderrContains(t *testing.T, substr string) {
	t.Helper()
	if !strings.Contains(r.Stderr, substr) {
		t.Errorf("Expected stderr to contain %q, but got:\n%s", substr, r.Stderr)
	}
}

// AssertStdoutNotContains asserts that stdout does not contain the given string
func (r *CommandResult) AssertStdoutNotContains(t *testing.T, substr string) {
	t.Helper()
	if strings.Contains(r.Stdout, substr) {
		t.Errorf("Expected stdout not to contain %q, but got:\n%s", substr, r.Stdout)
	}
}

// GetLines returns stdout split into lines
func (r *CommandResult) GetLines() []string {
	lines := strings.Split(strings.TrimSpace(r.Stdout), "\n")
	// Filter out empty lines
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}

// PrintOutput prints the command output for debugging
func (r *CommandResult) PrintOutput(t *testing.T) {
	t.Helper()
	t.Logf("Exit Code: %d", r.ExitCode)
	if r.Stdout != "" {
		t.Logf("Stdout:\n%s", r.Stdout)
	}
	if r.Stderr != "" {
		t.Logf("Stderr:\n%s", r.Stderr)
	}
	if r.Error != nil {
		t.Logf("Error: %v", r.Error)
	}
}

// GetRegion gets a valid region from the API for use in tests
func (c *TestConfig) GetRegion(t *testing.T) string {
	t.Helper()

	regionsResult := c.RunCommand(t, "regions", "list", "--no-header")
	regionsResult.PrintOutput(t)

	if regionsResult.ExitCode != 0 {
		t.Fatalf("Cannot list regions: %s", regionsResult.Stderr)
	}

	regions := regionsResult.GetLines()
	if len(regions) == 0 {
		t.Fatal("No regions available")
	}

	// Use the first region (format: identity or slug)
	regionLine := strings.Fields(regions[0])
	if len(regionLine) == 0 {
		t.Fatal("Invalid region format")
	}

	return regionLine[0]
}

// GetVPC gets a valid VPC from the API for use in tests
// If no VPC exists, it creates one and registers cleanup
func (c *TestConfig) GetVPC(t *testing.T) string {
	t.Helper()

	vpcsResult := c.RunCommand(t, "networking", "vpcs", "list", "--no-header")
	vpcsResult.PrintOutput(t)

	if vpcsResult.ExitCode != 0 {
		t.Fatalf("Cannot list VPCs: %s", vpcsResult.Stderr)
	}

	vpcs := vpcsResult.GetLines()
	if len(vpcs) > 0 {
		// Use the first existing VPC (format: identity is first column)
		vpcLine := strings.Fields(vpcs[0])
		if len(vpcLine) > 0 {
			return vpcLine[0]
		}
	}

	// No VPCs exist, create one
	return c.CreateVPC(t)
}

// CreateVPC creates a VPC for testing and registers cleanup
func (c *TestConfig) CreateVPC(t *testing.T) string {
	t.Helper()

	region := c.GetRegion(t)
	vpcName := "e2e-test-vpc-" + time.Now().Format("20060102150405")

	// Create VPC
	createResult := c.RunCommand(t, "networking", "vpcs", "create",
		"--name", vpcName,
		"--region", region,
		"--cidrs", "10.0.0.0/16",
		"--no-header")
	createResult.PrintOutput(t)

	if createResult.ExitCode != 0 {
		t.Fatalf("Failed to create VPC: %s", createResult.Stderr)
	}

	// Extract VPC identity from output
	createOutput := createResult.GetLines()
	if len(createOutput) == 0 {
		t.Fatal("Create command succeeded but produced no output")
	}

	vpcFields := strings.Fields(createOutput[0])
	if len(vpcFields) == 0 {
		t.Fatal("Create output format is unexpected")
	}
	vpcIdentity := vpcFields[0]

	// Register cleanup to delete the VPC
	t.Cleanup(func() {
		deleteResult := c.RunCommand(t, "networking", "vpcs", "delete", vpcIdentity, "--force")
		if deleteResult.ExitCode != 0 {
			t.Logf("Failed to clean up VPC %s: %s", vpcIdentity, deleteResult.Stderr)
		}
	})

	return vpcIdentity
}

// CreateSubnet creates a subnet in the given VPC for testing and registers cleanup
// It uses a unique CIDR based on timestamp to avoid conflicts
func (c *TestConfig) CreateSubnet(t *testing.T, vpcIdentity string) string {
	t.Helper()

	subnetName := "e2e-test-subnet-" + time.Now().Format("20060102150405")

	// Use timestamp to create a unique CIDR (10.0.X.0/24 where X is derived from timestamp)
	// This helps avoid conflicts when multiple tests run
	timestamp := time.Now().Unix()
	cidrOctet := int(timestamp%254) + 1 // Use 1-254 range
	cidr := fmt.Sprintf("10.0.%d.0/24", cidrOctet)

	// Create subnet
	createResult := c.RunCommand(t, "networking", "subnets", "create",
		"--name", subnetName,
		"--vpc", vpcIdentity,
		"--cidr", cidr,
		"--no-header")
	createResult.PrintOutput(t)

	if createResult.ExitCode != 0 {
		t.Fatalf("Failed to create subnet: %s", createResult.Stderr)
	}

	// Extract subnet identity from output
	createOutput := createResult.GetLines()
	if len(createOutput) == 0 {
		t.Fatal("Create command succeeded but produced no output")
	}

	subnetFields := strings.Fields(createOutput[0])
	if len(subnetFields) == 0 {
		t.Fatal("Create output format is unexpected")
	}
	subnetIdentity := subnetFields[0]

	// Register cleanup to delete the subnet
	t.Cleanup(func() {
		deleteResult := c.RunCommand(t, "networking", "subnets", "delete", subnetIdentity, "--force")
		if deleteResult.ExitCode != 0 {
			t.Logf("Failed to clean up subnet %s: %s", subnetIdentity, deleteResult.Stderr)
		}
	})

	return subnetIdentity
}

// GetVolume gets a valid volume from the API for use in tests
func (c *TestConfig) GetVolume(t *testing.T) string {
	t.Helper()

	volumesResult := c.RunCommand(t, "storage", "volumes", "list", "--no-header")
	volumesResult.PrintOutput(t)

	if volumesResult.ExitCode != 0 {
		t.Fatalf("Cannot list volumes: %s", volumesResult.Stderr)
	}

	volumes := volumesResult.GetLines()
	if len(volumes) == 0 {
		t.Fatal("No volumes available")
	}

	// Use the first volume (format: identity is first column)
	volumeLine := strings.Fields(volumes[0])
	if len(volumeLine) == 0 {
		t.Fatal("Invalid volume format")
	}

	return volumeLine[0]
}
