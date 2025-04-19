# gotmpl
Fast way to generate go-templates

## Example

### Quick run
Need to have template and data files under `templates` dir.


```bash
go install github.com/Samet-MohamedAmin/gotmpl@v0.1.2
```

Run:
``` bash
$ gotmpl
Searching for templates in templates
Looking for template: ALL
Found template: templates/example/template.go.tmpl
Processing template: templates/example/template.go.tmpl
Data path constructed: templates/example/data.yaml
Looking for data file: templates/example/data.yaml
Output directory path: example
Template name: example
Generated output file: output/example/welcome/hello.txt
Generated output file: output/example/example-01.txt
Generated output file: output/example/example-02.txt
```

Output:
```bash
$ tree
.
├── output
│   └── example
│       ├── example-01.txt
│       ├── example-02.txt
│       └── welcome
│           └── hello.txt
└── templates
    └── example
        ├── data.yaml
        └── template.go.tmpl

6 directories, 5 files
```

If you want to create test files, follow the next step.


### Dummy template and data file
You can run `gotmpl` under `example` dir.
Or you can create new files with the following commands:

```bash
dir=templates/example

# create template dir
mkdir -p "$dir"

# create template file
cat <<EOF > "$dir/template.go.tmpl"
# config ext=txt
---
# file: welcome/hello.txt
Hello, {{ .Name }}!
You have {{ .Count }} new messages.

{{- range \$index, \$item := .Items}}
- {{ printf "%02d" \$index }}-{{ \$item }}
{{- end}}

---
file 2
---
file 3
EOF

# create data file
cat <<EOF > "$dir/data.yaml"
Name: Charlie
Count: 5
Items:
  - Item 1
  - Item 2
EOF
```


## Overview
Easy and fast go template

## Features
- Organized folder structure for better code management.
- Modular design for easy maintenance and scalability.

## Folder Structure
```
config.yaml       # Main configuration file
go.mod            # Go module file
go.sum            # Go dependencies checksum
main.go           # Entry point of the application
README.md         # Project documentation
cmd/              # Command-line related code
example/          # Example configurations and outputs
pkg/              # Core packages for configuration and template processing
```

## Getting Started

### Prerequisites
- Go 1.20 or later
- Linux OS (recommended)

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/Samet-MohamedAmin/gotmpl.git
   ```
2. Navigate to the project directory:
   ```bash
   cd <repository-name>
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

```
Usage of gotmpl:
  This program processes template files with YAML data.
  Template files should be in the templates directory with extension .go.tmpl
  Data files should be in the templates directory with extension .yaml

  -clean
        Clean output directory before processing templates (default true)
  -config string
        Path to the configuration file (default "config.yaml")
  -help
        Show help message
  -separate
        Split output into multiple files at YAML document separators (---) (default true)
  -template string
        Template to process (use ALL for all templates) (default "ALL")
  -version
        Print the version and exit
```

1. Run the application:
   ```bash
   go run main.go
   ```
2. Use the example configurations in the `example/` folder to test the application.

