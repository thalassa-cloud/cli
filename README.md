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

Authenticate with Thalassa Cloud (opens a browser login by default):

```bash
tcloud context create --api=https://api.thalassa.cloud
```

Or with a personal access token:

```bash
tcloud context create --api=https://api.thalassa.cloud --token=<PAT>
```

## Configuration file

The CLI stores non-secret settings in `~/.tcloud`. Sensitive credentials are kept out of this file when a system credential store is available (see [Credential storage](#credential-storage) below).

### With personal access token

When using the file-backed credential store, tokens may appear in the config. With the default keychain storage, the config only records where credentials are stored:

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
      credentialStore: keychain
      user: {}
```

## Credential storage

By default, the CLI stores secrets in the OS credential store when one is available:

- **macOS:** Keychain (service name `Thalassa Cloud CLI`)
- **Windows:** Credential Manager
- **Linux:** Secret Service (for example GNOME Keyring or KWallet)

Stored secrets include personal access tokens, OIDC client credentials, and browser login access/refresh tokens.

Control the store with `THALASSA_CREDENTIAL_STORE`:

| Value | Behaviour |
|-------|-----------|
| `auto` (default) | Use the OS credential store when available; otherwise store secrets in the config file |
| `keychain` | Always use the OS credential store; fail if unavailable |
| `file` | Store secrets in `~/.tcloud` (legacy behaviour) |

Existing plaintext credentials in `~/.tcloud` are migrated to the keychain automatically on load when keychain storage is enabled.

To inspect or remove keychain entries on macOS, open **Keychain Access** and search for `Thalassa Cloud CLI`.

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
