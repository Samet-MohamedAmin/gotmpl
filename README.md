# Project Name

## Overview
Easy and fast go template

## Features
- Modular design for easy maintenance and scalability.
- Example configurations and templates for quick setup.
- Organized folder structure for better code management.

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

