# Contributing to Thalassa Cloud CLI

Thank you for your interest in contributing to the Thalassa Cloud CLI! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Code Style and Conventions](#code-style-and-conventions)
- [Testing](#testing)
- [Commit Messages](#commit-messages)

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/cli.git
   cd cli
   ```
3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/thalassa-cloud/cli.git
   ```

## Development Setup

### Prerequisites

- **Go 1.24 or later** - [Install Go](https://golang.org/doc/install)
- **Make** - Usually pre-installed on Unix systems
- **Git** - For version control

### Building the Project

```bash
# Build the binary
make build

# The binary will be available at ./bin/tcloud
./bin/tcloud --help
```

### Running Tests

```bash
# Run unit tests
make test

# Run E2E tests (requires configuration)
# NOTE: Creates real resources, which may cost money
make test-e2e
```

See the [Testing](#testing) section for more details.

## Project Structure

```
cli/
├── cmd/              # CLI commands organized by feature
│   ├── context/      # Context management commands
│   ├── iaas/         # IaaS-related commands
│   ├── kubernetes/   # Kubernetes commands
│   └── ...
├── internal/         # Internal packages (not exported)
│   ├── config/       # Configuration management
│   ├── table/        # Table formatting utilities
│   └── ...
├── e2e/              # End-to-end tests
├── docs/             # Documentation
└── main.go          # Application entry point
```

### Command Organization

Commands are organized in the `cmd/` directory following this structure:
- Each top-level command has its own package (e.g., `cmd/kubernetes/`)
- Subcommands are in the same package or subdirectories
- Command files follow naming conventions:
  - `{command}.go` - Main command definition
  - `{subcommand}.go` - Subcommand implementation

## Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the code style guidelines

3. **Write tests** for your changes (see [Testing](#testing))

4. **Run tests** to ensure everything passes:
   ```bash
   make test
   ```

5. **Commit your changes** using semantic commit messages (see [Commit Messages](#commit-messages))

6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request** on GitHub

## Code Style and Conventions

### Go Code Style

- Follow standard Go formatting: `go fmt ./...`
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `golangci-lint` or similar tools for linting (if configured)

## Testing

### Unit Tests

Unit tests are located alongside the code they test (e.g., `labels_test.go` in the `labels` package).

- Use table-driven tests for multiple test cases
- Use the `testify` package for assertions
- Keep tests simple and focused on one behavior

**Example:**
```go
func TestParseLabelSelector(t *testing.T) {
    tests := []struct {
        name     string
        selector string
        expected map[string]string
    }{
        {
            name:     "single label",
            selector: "key1=value1",
            expected: map[string]string{"key1": "value1"},
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ParseLabelSelector(tt.selector)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

**Running Unit Tests:**
```bash
# Run all unit tests
make test

# Run tests for a specific package
go test ./internal/labels/...

# Run tests with verbose output
go test -v ./internal/labels/...
```

### E2E Tests

End-to-end tests are located in the `e2e/` directory and test the actual CLI binary against a real API.

**Setup:**
1. Build the binary: `make build`
2. Configure environment variables (see `e2e/README.md` for details):
   ```bash
   export TCLOUD_E2E_API_ENDPOINT="https://api.thalassa.cloud"
   export TCLOUD_E2E_PERSONAL_ACCESS_TOKEN="your-token"
   ```

**Running E2E Tests:**
```bash
# Run all E2E tests
make test-e2e

# Run specific E2E test
go test ./e2e/... -v -run TestVersion

# Run with custom binary path
TCLOUD_E2E_BINARY_PATH=/path/to/tcloud go test ./e2e/... -v
```

For more details, see `e2e/README.md`.

## Commit Messages

We follow **semantic commit message** conventions. Commit messages should be clear, descriptive, and follow this format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring (no behavior change)
- `style`: Code style changes (formatting, etc.)
- `chore`: Maintenance tasks
- `perf`: Performance improvements

### Examples

```
feat(vpcs): add support for label selectors in list command

This allows users to filter VPCs by labels when listing them,
improving usability for managing large numbers of resources.

Closes #123
```

```
fix(auth): handle expired tokens gracefully

Previously, expired tokens would cause a panic. Now they are
handled with a clear error message prompting the user to re-authenticate.
```

```
test(e2e): add E2E tests for me organisations command

Adds comprehensive E2E tests covering all flags and output formats
for the me organisations command.
```

### Guidelines

- Use imperative mood ("add" not "added" or "adds")
- Keep the subject line under 50 characters
- Capitalize the first letter of the subject
- Don't end the subject with a period
- Use the body to explain *what* and *why* (not *how*)
- Reference issues and PRs in the footer
