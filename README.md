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

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
