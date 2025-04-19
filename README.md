# gotmpl

A powerful Go template processor for generating files from templates using YAML data. Perfect for code generation, configuration management, and documentation workflows.

## Quick Links

- [Installation Guide](docs/installation.md)
- [User Guide](docs/user-guide.md)
- [Template Reference](docs/template-reference.md)
- [CLI Reference](docs/cli-reference.md)
- [Configuration Guide](docs/configuration.md)
- [Examples](docs/examples.md)

## Features

- 🚀 **High performance template processing**
- 📦 **Multiple output files from a single template**
- 🧩 **Template-specific configurations**
- 📂 **Organized output with directory structure preservation**
- 🧰 **Rich command-line interface**
- 🔍 **Smart template and data discovery**
- 🔄 **YAML document separator support**
- ⚙️ **Flexible configuration system**

## Quick Start

1. Install gotmpl:
   ```bash
   go install github.com/Samet-MohamedAmin/gotmpl@latest
   ```

2. Create a template and data file:
   ```bash
   mkdir -p templates/example
   ```

3. Generate files:
   ```bash
   gotmpl gen ./templates
   ```

For detailed instructions, see the [User Guide](docs/user-guide.md).

## Project Structure

```
.
├── cmd/               # Command-line interface
├── pkg/              # Core packages
│   ├── config/      # Configuration handling
│   └── template/    # Template processing
├── docs/            # Documentation
└── example/         # Example configurations
```

## Contributing

Contributions are welcome! Please see our [Contributing Guide](docs/contributing.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

