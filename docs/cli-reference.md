# CLI Reference

## Command Overview

```bash
gotmpl [command] [flags] [directory]
```

## Commands

### gen

Generate files from templates.

```bash
gotmpl gen [flags] [directory]
```

#### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--template` | `-t` | `ALL` | Template to process (use ALL for all templates) |
| `--separate` | `-s` | `true` | Split output into multiple files at YAML document separators (---) |
| `--clean` | `-c` | `true` | Clean output directory before processing templates |
| `--config` | `-f` | `config.yaml` | Path to the configuration file |
| `--output` | `-o` | `output` | Output directory for generated files |
| `--multiple` | `-m` | `false` | Process multiple template directories |

#### Examples

```bash
# Generate all templates
gotmpl gen ./templates

# Generate specific template
gotmpl gen ./templates --template=example

# Generate with custom output
gotmpl gen ./templates --output=generated

# Generate without separating files
gotmpl gen ./templates --separate=false

# Generate from multiple directories
gotmpl gen ./templates --multiple=true
```

### completion

Generate shell completion scripts.

```bash
gotmpl completion [bash|zsh]
```

#### Examples

```bash
# Generate bash completion
gotmpl completion bash > /etc/bash_completion.d/gotmpl

# Generate zsh completion
gotmpl completion zsh > "${fpath[1]}/_gotmpl"
```

### version

Print version information.

```bash
gotmpl version
```

### help

Show help for any command.

```bash
gotmpl help [command]
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--help` | `-h` | Show help message |
| `--version` | `-v` | Print version and exit |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GOTMPL_CONFIG` | Path to configuration file (overrides --config) |
| `GOTMPL_OUTPUT` | Output directory (overrides --output) |

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid command or flag |
| 3 | Template processing error |
| 4 | File system error |

## Configuration File

The default configuration file is `config.yaml` in the current directory. You can specify a different file with the `--config` flag.

Example configuration:

```yaml
OutputDir: "output"
OutputExtension: ".yaml"
TemplateFile: "template.go.tmpl"
DataFile: "data.yaml"
DefaultPrefix: "file"
```

## Related Documentation

- [User Guide](user-guide.md) - How to use gotmpl
- [Template Reference](template-reference.md) - Template syntax and features
- [Configuration Guide](configuration.md) - Advanced configuration options
- [Examples](examples.md) - Usage examples 