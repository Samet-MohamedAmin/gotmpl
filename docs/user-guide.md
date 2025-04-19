# User Guide

## Basic Usage

### Command Structure

```bash
gotmpl [command] [flags] [directory]
```

### Available Commands

- `gen` - Generate files from templates (default command)
- `completion` - Generate shell completion scripts
- `version` - Print version information
- `help` - Show help for any command

### Basic Example

1. Create a template directory:
   ```bash
   mkdir -p templates/example
   ```

2. Create a template file (`templates/example/template.go.tmpl`):
   ```go
   # config ext=txt separate=true
   ---
   # file: {{ .Name | lower }}/hello.txt
   Hello, {{ .Name }}!
   You have {{ .Count }} new messages.
   
   {{- range $index, $item := .Items}}
   - {{ printf "%02d" $index }}-{{ $item }}
   {{- end}}
   ```

3. Create a data file (`templates/example/data.yaml`):
   ```yaml
   Name: Alice
   Count: 3
   Items:
     - First item
     - Second item
     - Third item
   ```

4. Generate files:
   ```bash
   gotmpl gen ./templates
   ```

## Template Structure

### Configuration Directives

Configuration directives are specified at the top of your template:

```
# config ext=go separate=true
```

Supported directives:
- `ext`: Output file extension
- `separate`: Whether to separate output into multiple files (true/false)

### File Directives

Specify output filenames with file directives:

```
# file: path/to/output.txt
```

### Document Separators

Use YAML document separators (`---`) to split your template into multiple files.

## Advanced Usage

### Processing Multiple Templates

Process all templates in subdirectories:

```bash
gotmpl gen ./templates --multiple=true
```

### Generating a Single File

Disable separating the output:

```bash
gotmpl gen ./templates/example --separate=false
```

### Custom Output Directory

Specify a different output directory:

```bash
gotmpl gen ./templates --output=generated
```

### Processing Specific Templates

Process only specific templates:

```bash
gotmpl gen ./templates --template=example
```

### Shell Completion

Generate shell completion scripts:

```bash
# For bash
gotmpl completion bash > /etc/bash_completion.d/gotmpl

# For zsh
gotmpl completion zsh > "${fpath[1]}/_gotmpl"
```

## Best Practices

1. **Organize Templates**
   - Keep related templates in the same directory
   - Use meaningful names for template directories
   - Separate different types of templates

2. **Data Organization**
   - Keep data files close to their templates
   - Use consistent naming conventions
   - Document your data structure

3. **Output Management**
   - Use the `--clean` flag to ensure fresh output
   - Specify output directories for different purposes
   - Use meaningful file names in file directives

4. **Template Design**
   - Use configuration directives for flexibility
   - Document template requirements
   - Test templates with different data

## Common Patterns

### Multiple Output Files

```go
# config ext=txt separate=true
---
# file: {{ .Name | lower }}/hello.txt
Hello, {{ .Name }}!
---
# file: {{ .Name | lower }}/info.txt
This is information for {{ .Name }}.
```

### Conditional Output

```go
{{- if .Enabled }}
# file: enabled.txt
Feature is enabled
{{- else }}
# file: disabled.txt
Feature is disabled
{{- end }}
```

### Looping Over Data

```go
{{- range $index, $item := .Items }}
# file: item{{ $index }}.txt
Item {{ $index }}: {{ $item }}
{{- end }}
```

## Next Steps

- [Template Reference](template-reference.md) - Detailed template syntax
- [CLI Reference](cli-reference.md) - Complete command reference
- [Configuration Guide](configuration.md) - Advanced configuration options
- [Examples](examples.md) - More usage examples 