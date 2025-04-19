# Installation Guide

## Prerequisites

- Go 1.16 or later
- Git (for source installation)

## Installation Methods

### Using Go Install (Recommended)

The simplest way to install gotmpl is using `go install`:

```bash
go install github.com/Samet-MohamedAmin/gotmpl@latest
```

This will install the latest version of gotmpl to your `$GOPATH/bin` directory.

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/Samet-MohamedAmin/gotmpl.git
   cd gotmpl
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Install the binary:
   ```bash
   go install
   ```

## Verifying Installation

After installation, verify that gotmpl is available:

```bash
gotmpl --version
```

You should see the current version number displayed.

## Updating

To update to the latest version:

```bash
go install github.com/Samet-MohamedAmin/gotmpl@latest
```

## Troubleshooting

### Common Issues

1. **Command not found**
   - Ensure `$GOPATH/bin` is in your `$PATH`
   - Try running with full path: `$GOPATH/bin/gotmpl`

2. **Permission denied**
   - Make sure you have write permissions to the installation directory
   - Try running with sudo: `sudo go install github.com/Samet-MohamedAmin/gotmpl@latest`

3. **Go version too old**
   - Update your Go installation to version 1.16 or later

## Next Steps

- [User Guide](user-guide.md) - Learn how to use gotmpl
- [Configuration Guide](configuration.md) - Configure gotmpl for your needs
- [Examples](examples.md) - See gotmpl in action 