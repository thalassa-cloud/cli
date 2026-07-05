# Thalassa Cloud CLI (tcloud)

A command-line interface for managing your Thalassa Cloud Installation.

> This project is still in beta. Commands and UX may change while the project is in initial development.

## Installation

### macOS (using Homebrew)

```bash
brew install thalassa-cloud/tap/tcloud
```

### Manual Installation

Download the latest release for your platform from the [GitHub releases page](https://github.com/thalassa-cloud/cli/releases).

## Quick Start

Authenticate with Thalassa Cloud:
```bash
tcloud context create --api=https://api.thalassa.cloud --token=<PAT>
```

## Configuration file

### With personal access token

```yaml
configVersion: v1
contexts:
    - name: default
      context:
        api: api.thalassa.cloud
        user: default
        organisation: <ORG_SLUG_OR_IDENTITY>
current-context: default
servers:
    - name: api.thalassa.cloud
      api:
        server: https://api.thalassa.cloud
users:
    - name: default
      user:
        token: <PAT>
```

## Config file permissions

The CLI stores settings in `~/.tcloud`. On load, the CLI warns if that file is world-readable. Fix permissions with:

```bash
tcloud context fix
```

This restricts `~/.tcloud` to owner-only access (`0600`).

### Migrating credentials to the keychain

New logins store secrets in the OS credential store (macOS Keychain, Windows Credential Manager, or Linux Secret Service) when available. Existing configs with plaintext tokens in `~/.tcloud` are left unchanged until manually running the `fix` command.

To move existing plaintext credentials into the keychain:

```bash
tcloud context fix --migrate-credentials
```

Use `tcloud context view` to inspect your config with secrets redacted.

Control credential storage with `THALASSA_CREDENTIAL_STORE` (`auto`, `keychain`, or `file`).

## Development

### Prerequisites

- Go 1.24 or later
- Make

### Building from Source

```bash
# Clone the repository
git clone https://github.com/thalassa-cloud/cli.git
cd cli

# Build the binary
make build

# Run tests
make test
```

### Run E2E tests

> Note: Running E2E tests creates real resources and you may be charged for these!

```bash
# Build the binary first
make build

# Set environment variables
export TCLOUD_E2E_API_ENDPOINT="https://api.thalassa.cloud"
export TCLOUD_E2E_PERSONAL_ACCESS_TOKEN="your-token"
export TCLOUD_E2E_ORGANISATION="org-id"

# Run tests
make test-e2e
# or
go test ./e2e/... -v
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
