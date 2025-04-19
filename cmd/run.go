package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/Samet-MohamedAmin/gotmpl/pkg/config"
	"github.com/Samet-MohamedAmin/gotmpl/pkg/template"
)

func parseFlags() (string, bool, bool, string) {
	help := flag.Bool("help", false, "Show help message")
	templateName := flag.String("template", "ALL", "Template to process (use ALL for all templates)")
	separate := flag.Bool("separate", true, "Split output into multiple files at YAML document separators (---)")
	clean := flag.Bool("clean", true, "Clean output directory before processing templates")
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	version := flag.Bool("version", false, "Print the version and exit")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", "gotmpl")
		fmt.Fprintf(flag.CommandLine.Output(), "  This program processes template files with YAML data.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  Template files should be in the templates directory with extension .go.tmpl\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  Data files should be in the templates directory with extension .yaml\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return "", false, false, ""
	}

	if *version {
		fmt.Println("Version:", config.Version)
		return "", false, false, ""
	}

	return *templateName, *separate, *clean, *configPath
}

func cleanOutputDir() error {
	outputDir := config.OutputDir
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("failed to clean output directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	return nil
}

func Run() error {
	// Parse flags
	templateName, separate, clean, configPath := parseFlags()
	if templateName == "" {
		return nil
	}

	// Initialize configuration with the provided config path
	if err := config.Initialize(configPath); err != nil {
		return fmt.Errorf("failed to initialize configuration: %v", err)
	}

	if clean {
		if err := cleanOutputDir(); err != nil {
			return err
		}
	}

	finder := template.NewFinder()
	templateFiles, err := finder.FindTemplates(templateName)
	if err != nil {
		return err
	}

	processor := template.NewProcessor(separate)
	for _, templatePath := range templateFiles {
		if err := processor.ProcessTemplate(templatePath); err != nil {
			return err
		}
	}

	return nil
}
