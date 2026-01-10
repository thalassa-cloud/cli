# E2E Tests

This directory contains end-to-end (E2E) tests for the Thalassa Cloud CLI. These tests execute the actual CLI binary against a real API endpoint to verify that commands work correctly.

## Prerequisites

1. Build the CLI binary:
   ```bash
   make build
   ```

2. Set up authentication credentials. You need at least one of:
   - Personal Access Token
   - Access Token
   - OIDC Client ID and Secret

## Configuration

E2E tests are configured using environment variables:

### Required Variables

- `TCLOUD_E2E_API_ENDPOINT`: The API endpoint URL (e.g., `https://api.thalassa.cloud`)

### Authentication (at least one required)

- `TCLOUD_E2E_PERSONAL_ACCESS_TOKEN`: Personal access token for authentication
- `TCLOUD_E2E_ACCESS_TOKEN`: Access token for authentication
- `TCLOUD_E2E_OIDC_CLIENT_ID` + `TCLOUD_E2E_OIDC_CLIENT_SECRET`: OIDC credentials

### Optional Variables

- `TCLOUD_E2E_ORGANISATION`: Organisation slug or identity to use for tests
- `TCLOUD_E2E_BINARY_PATH`: Path to the CLI binary (defaults to `./bin/tcloud`)

## Running Tests

### Run all E2E tests

```bash
go test ./e2e/... -v
```

### Run specific test

```bash
go test ./e2e/... -v -run TestVersion
```

### Run with environment variables

```bash
export TCLOUD_E2E_API_ENDPOINT="https://api.thalassa.cloud"
export TCLOUD_E2E_PERSONAL_ACCESS_TOKEN="your-token"
go test ./e2e/... -v
```

### Using a .env file

You can create a `.env` file in the project root and source it:

```bash
source .env
go test ./e2e/... -v
```
