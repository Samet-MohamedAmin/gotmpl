package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Samet-MohamedAmin/gotmpl/pkg/config"
	"github.com/Samet-MohamedAmin/gotmpl/pkg/template"
)

// CLIOptions holds all command-line options
type CLIOptions struct {
	SourceDir    string
	OutputDir    string
	TemplateName string
	Separate     bool
	Clean        bool
	Multiple     bool
	ConfigPath   string
	ShowHelp     bool
	ShowVersion  bool
}

// ParseFlags parses command-line flags and returns structured options
func ParseFlags() (*CLIOptions, error) {
	opts := &CLIOptions{}

	flag.BoolVar(&opts.ShowHelp, "help", false, "Show help message")
	flag.StringVar(&opts.TemplateName, "template", "ALL", "Template to process (use ALL for all templates)")
	flag.BoolVar(&opts.Separate, "separate", true, "Split output into multiple files at YAML document separators (---)")
	flag.BoolVar(&opts.Clean, "clean", true, "Clean output directory before processing templates")
	flag.StringVar(&opts.ConfigPath, "config", "config.yaml", "Path to the configuration file")
	flag.StringVar(&opts.OutputDir, "output", "output", "Output directory for generated files")
	flag.BoolVar(&opts.Multiple, "multiple", false, "Process multiple template directories")
	flag.BoolVar(&opts.ShowVersion, "version", false, "Print the version and exit")

	flag.Usage = printUsage
	flag.Parse()

	if opts.ShowHelp {
		flag.Usage()
		return opts, nil
	}

	if opts.ShowVersion {
		fmt.Println("Version:", config.Version)
		return opts, nil
	}

	// Get the directory from the last argument
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return nil, fmt.Errorf("missing required argument: directory")
	}

	opts.SourceDir = args[len(args)-1]
	return opts, nil
}

// printUsage prints detailed usage information
func printUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", "gotmpl")
	fmt.Fprintf(flag.CommandLine.Output(), "  This program processes template files with YAML data.\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  Template files should be in the source directory with extension .go.tmpl\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  Data files should be in the source directory with extension .yaml\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  gotmpl [flags] <directory>\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Required Arguments:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  <directory>\n")
	fmt.Fprintf(flag.CommandLine.Output(), "        Source directory containing templates (required)\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Optional Flags:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nConfiguration File Options:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  The following options can be set in the config file (default: config.yaml):\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  OutputDir: string\n")
	fmt.Fprintf(flag.CommandLine.Output(), "        Output directory for generated files (default: output)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  OutputExtension: string\n")
	fmt.Fprintf(flag.CommandLine.Output(), "        Extension for output files (default: none)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  TemplateFile: string\n")
	fmt.Fprintf(flag.CommandLine.Output(), "        Name of template files (default: template.go.tmpl)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  DataFile: string\n")
	fmt.Fprintf(flag.CommandLine.Output(), "        Name of data files (default: data.yaml)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  DefaultPrefix: string\n")
	fmt.Fprintf(flag.CommandLine.Output(), "        Default prefix for output files (default: file)\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Examples:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  # Generate all templates from a directory\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  gotmpl ./templates\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  # Generate specific template\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  gotmpl -template=example ./templates\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  # Generate with custom config\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  gotmpl -config=custom.yaml ./templates\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  # Generate without separating files\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  gotmpl -separate=false ./templates\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  # Generate with custom output directory\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  gotmpl -output=generated ./templates\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  # Generate from single directory (template.go.tmpl and data.yaml)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  gotmpl -multiple=false ./templates/example\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  # Generate from multiple directories\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  gotmpl -multiple=true ./templates\n")
}

// cleanOutputDir removes and recreates the output directory
func cleanOutputDir(outputDir string) error {
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("failed to clean output directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	return nil
}

// Run executes the main program logic
func Run() error {
	// Parse flags
	opts, err := ParseFlags()
	if err != nil {
		return err
	}

	// Exit early if just showing help or version
	if opts.ShowHelp || opts.ShowVersion {
		return nil
	}

	// Initialize configuration with the provided config path
	if err := config.Initialize(opts.ConfigPath); err != nil {
		return fmt.Errorf("failed to initialize configuration: %v", err)
	}

	// Set the output directory in config
	config.OutputDir = opts.OutputDir

	var templateFiles []string

	if !opts.Multiple {
		// Single directory mode: look for template.go.tmpl and data.yaml in src dir
		templateFiles, err = validateSingleDirectoryMode(opts.SourceDir)
		if err != nil {
			return err
		}
	} else {
		// Multiple directory mode: use the finder to locate templates
		finder := template.NewFinder(opts.SourceDir)
		templateFiles, err = finder.FindTemplates(opts.TemplateName)
		if err != nil {
			return err
		}
	}

	// Clean the output directory after validation but before generation
	if opts.Clean {
		if err := cleanOutputDir(opts.OutputDir); err != nil {
			return err
		}
	}

	// Process all templates
	processor := template.NewProcessor(opts.Separate)
	for _, templatePath := range templateFiles {
		if err := processor.ProcessTemplate(templatePath, opts.Multiple); err != nil {
			return err
		}
	}

	return nil
}

// validateSingleDirectoryMode checks if the necessary files exist in single directory mode
func validateSingleDirectoryMode(srcDir string) ([]string, error) {
	templatePath := filepath.Join(srcDir, config.TemplateFile)
	dataPath := filepath.Join(srcDir, config.DataFile)

	// Check if both files exist
	if _, err := os.Stat(templatePath); err != nil {
		return nil, fmt.Errorf(`template file not found in %s: %v

If you want to process multiple templates from subdirectories, please use the -multiple flag:
  gotmpl -multiple=true %s

Otherwise, make sure you have the following files in your source directory:
  - %s
  - %s`, srcDir, err, srcDir, config.TemplateFile, config.DataFile)
	}

	if _, err := os.Stat(dataPath); err != nil {
		return nil, fmt.Errorf(`data file not found in %s: %v

If you want to process multiple templates from subdirectories, please use the -multiple flag:
  gotmpl -multiple=true %s

Otherwise, make sure you have the following files in your source directory:
  - %s
  - %s`, srcDir, err, srcDir, config.TemplateFile, config.DataFile)
	}

	return []string{templatePath}, nil
}
