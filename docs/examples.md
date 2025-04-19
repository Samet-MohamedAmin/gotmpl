# Examples

This document provides comprehensive examples of using gotmpl with different options and configurations.

## Basic Usage

### Generate All Templates

```bash
gotmpl gen ./templates
```

This will:
- Process all templates in the `./templates` directory
- Use default configuration from `config.yaml`
- Output files to the `output` directory
- Split output into multiple files at YAML document separators
- Clean the output directory before processing

### Generate Specific Template

```bash
gotmpl gen ./templates --template=example
```

This will:
- Process only the template named "example"
- Use default configuration from `config.yaml`
- Output files to the `output` directory
- Split output into multiple files at YAML document separators
- Clean the output directory before processing

## Output Options

### Custom Output Directory

```bash
gotmpl gen ./templates --output=generated
```

This will:
- Process all templates
- Output files to the `generated` directory instead of `output`
- Create the directory if it doesn't exist

### Disable Output Separation

```bash
gotmpl gen ./templates --separate=false
```

This will:
- Process all templates
- Combine all output into a single file
- Ignore YAML document separators (---)

### Disable Output Directory Cleaning

```bash
gotmpl gen ./templates --clean=false
```

This will:
- Process all templates
- Keep existing files in the output directory
- Append new files to the existing directory

## Template Processing Options

### Process Multiple Directories

```bash
gotmpl gen ./templates --multiple=true
```

This will:
- Process templates from all subdirectories
- Create separate output directories for each template
- Look for `template.go.tmpl` and `data.yaml` in each subdirectory

### Process Single Directory

```bash
gotmpl gen ./templates/example --multiple=false
```

This will:
- Process only the specified directory
- Look for `template.go.tmpl` and `data.yaml` in the current directory
- Output files to the main output directory

## Configuration Options

### Custom Configuration File

```bash
gotmpl gen ./templates --config=custom.yaml
```

This will:
- Use `custom.yaml` instead of the default `config.yaml`
- Override default configuration values
- Support the following configuration options:
  - `OutputDir`: Output directory path
  - `OutputExtension`: File extension for output files
  - `TemplateFile`: Name of template files
  - `DataFile`: Name of data files
  - `DefaultPrefix`: Default prefix for output files

## Template-Specific Configuration

You can also configure options within your template files:

```go
# config ext=txt separate=true
---
# file: {{ .Name | lower }}/hello.txt
Hello, {{ .Name }}!
```

This will:
- Override global configuration for this template
- Set output extension to `.txt`
- Enable output separation
- Apply only to this specific template

## Shell Completion

### Bash Completion

```bash
# Generate completion script
gotmpl completion bash > /etc/bash_completion.d/gotmpl

# Or load temporarily
source <(gotmpl completion bash)
```

### Zsh Completion

```bash
# Generate completion script
gotmpl completion zsh > "${fpath[1]}/_gotmpl"

# Or load temporarily
source <(gotmpl completion zsh)
```

## Version Information

```bash
gotmpl version
```

This will:
- Display the current version of gotmpl
- Exit immediately after showing version

## Help Information

```bash
# Show general help
gotmpl help

# Show help for specific command
gotmpl help gen
```

This will:
- Display detailed usage information
- Show available options and their descriptions
- Provide examples for the specified command 