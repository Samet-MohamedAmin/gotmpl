# Template Reference

## Template Structure

A template file consists of two parts:
1. Configuration section (optional)
2. Template content

The configuration section is separated from the template content by a line containing three dashes (`---`).

### Configuration Section

The configuration section is a YAML document that specifies how the template should be processed. It can include:

```yaml
output: "output/path/file.ext"  # Output file path
data: "data.yaml"              # Data file path (optional)
```

### Template Content

The template content uses Go's text/template syntax. Here's a basic example:

```go
Hello, {{.Name}}!

You have {{.Count}} items:
{{range .Items}}
- {{.}}
{{end}}
```

## Template Syntax

### Variables

Access data using dot notation:

```go
{{.VariableName}}
```

### Control Structures

#### If-Else

```go
{{if .Condition}}
    // Content when condition is true
{{else}}
    // Content when condition is false
{{end}}
```

#### Range

```go
{{range .Items}}
    // Content for each item
    // Use {{.}} to access the current item
{{end}}
```

#### With

```go
{{with .Object}}
    // Content with .Object as the current context
{{end}}
```

### Functions

#### Built-in Functions

- `len`: Get length of array/slice/string
- `index`: Get element at index
- `printf`: Format string
- `html`: Escape HTML
- `js`: Escape JavaScript
- `urlquery`: Escape URL query

Example:
```go
{{printf "Count: %d" .Count}}
```

#### Custom Functions

You can define custom functions in your template:

```go
{{define "greet"}}Hello, {{.}}!{{end}}
{{template "greet" .Name}}
```

## Data Structure

Data is provided in YAML format. Example:

```yaml
Name: "World"
Count: 3
Items:
  - "Item 1"
  - "Item 2"
  - "Item 3"
```

## Multiple Documents

You can have multiple YAML documents in your data file, separated by `---`. Each document will be processed separately:

```yaml
Name: "Document 1"
Count: 1
---
Name: "Document 2"
Count: 2
```

## Best Practices

1. **Keep Templates Simple**
   - Avoid complex logic in templates
   - Move complex operations to data preparation

2. **Use Meaningful Names**
   - Use descriptive variable names
   - Group related data in structures

3. **Handle Missing Data**
   - Use `if` to check for existence
   - Provide default values where appropriate

4. **Format Output**
   - Use proper indentation
   - Add comments for complex sections

## Common Patterns

### Conditional Output

```go
{{if .Enabled}}
Enabled: true
{{else}}
Enabled: false
{{end}}
```

### Looping with Index

```go
{{range $index, $item := .Items}}
Item {{$index}}: {{$item}}
{{end}}
```

### Nested Data

```go
{{with .User}}
Name: {{.Name}}
Email: {{.Email}}
{{end}}
```

## Related Documentation

- [User Guide](user-guide.md) - How to use templates
- [CLI Reference](cli-reference.md) - Command-line options
- [Configuration Guide](configuration.md) - Advanced configuration
- [Examples](examples.md) - Template examples 