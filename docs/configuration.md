# Advanced Configuration

## Configuration File (config.yaml)

The configuration file supports the following properties:

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `OutputDir` | string | `output` | Directory for generated files |
| `OutputExtension` | string | `""` | Extension for output files (e.g., `.yaml`, `.txt`) |
| `TemplateFile` | string | `template.go.tmpl` | Name of template files |
| `DataFile` | string | `data.yaml` | Name of data files |
| `DefaultPrefix` | string | `file` | Default prefix for output files |

Example:
```yaml
OutputDir: "generated"
OutputExtension: ".yaml"
TemplateFile: "template.go.tmpl"
DataFile: "data.yaml"
DefaultPrefix: "output"
```

## Template Configuration

### Global Configuration

Add at the beginning of your template file:
```go
# config ext=json separate=false
```

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `ext` | string | `""` | Output file extension |
| `separate` | bool | `true` | Split output into multiple files |

### File Configuration

After document separator (`---`), specify output file:
```go
# file: path/to/output.json
```

Example:
```go
# config ext=json separate=false
---
# file: api/response.json
{
    "status": "{{ .Status }}",
    "data": {{ .Data | toJson }}
}
``` 