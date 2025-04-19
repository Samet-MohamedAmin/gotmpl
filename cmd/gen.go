package cmd

import (
	"fmt"

	"github.com/Samet-MohamedAmin/gotmpl/pkg/config"
	"github.com/Samet-MohamedAmin/gotmpl/pkg/template"
	"github.com/spf13/cobra"
)

var (
	sourceDir    string
	outputDir    string
	templateName string
	separate     bool
	clean        bool
	multiple     bool
	configPath   string
)

var genCmd = &cobra.Command{
	Use:   "gen [flags] [directory]",
	Short: "Generate output from templates",
	Long: `Generate files from templates in the specified directory.
Template files should be in the source directory with extension .go.tmpl
Data files should be in the source directory with extension .yaml`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceDir = args[0]
		return RunGen(sourceDir, outputDir, templateName, separate, clean, multiple, configPath)
	},
}

// RunGen wraps the existing Run logic for use with cobra
func RunGen(srcDir, outDir, tmplName string, sep, cln, mult bool, cfgPath string) error {
	// Create options struct to match what the existing Run function expects
	opts := &CLIOptions{
		SourceDir:    srcDir,
		OutputDir:    outDir,
		TemplateName: tmplName,
		Separate:     sep,
		Clean:        cln,
		Multiple:     mult,
		ConfigPath:   cfgPath,
	}

	// Initialize and run using existing functionality
	return runWithOptions(opts)
}

// Extract the core logic from Run() function to be reusable with cobra
func runWithOptions(opts *CLIOptions) error {
	// This function contains the logic from Run() but accepts the options directly
	// Initialize configuration with the provided config path
	if err := config.Initialize(opts.ConfigPath); err != nil {
		return fmt.Errorf("failed to initialize configuration: %v", err)
	}

	// Set the output directory in config
	config.OutputDir = opts.OutputDir

	var templateFiles []string
	var err error

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

func init() {
	rootCmd.AddCommand(genCmd)

	// Add flags to match the existing CLI options
	genCmd.Flags().StringVarP(&templateName, "template", "t", "ALL", "Template to process (use ALL for all templates)")
	genCmd.Flags().BoolVarP(&separate, "separate", "s", true, "Split output into multiple files at YAML document separators (---)")
	genCmd.Flags().BoolVarP(&clean, "clean", "c", true, "Clean output directory before processing templates")
	genCmd.Flags().StringVarP(&configPath, "config", "f", "config.yaml", "Path to the configuration file")
	genCmd.Flags().StringVarP(&outputDir, "output", "o", "output", "Output directory for generated files")
	genCmd.Flags().BoolVarP(&multiple, "multiple", "m", false, "Process multiple template directories")
}
