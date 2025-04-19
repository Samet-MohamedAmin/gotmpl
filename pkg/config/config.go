package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const Version = "v0.1.2"

// Config holds all configuration values
type Config struct {
	// Directories
	OutputDir   string `yaml:"OutputDir"`
	TemplateDir string `yaml:"TemplateDir"`

	// File patterns
	OutputExtension string `yaml:"OutputExtension"`
	TemplatePattern string `yaml:"TemplatePattern"`
	DataPattern     string `yaml:"DataPattern"`

	// Special values
	AllTemplates string `yaml:"AllTemplates"`
}

// Default configuration values
var defaultConfig = Config{
	OutputDir:   "output",
	TemplateDir: "templates",

	OutputExtension: "", // ".yaml"
	TemplatePattern: "template.go.tmpl",
	DataPattern:     "data.yaml",

	AllTemplates: "ALL",
}

// Current configuration
var (
	OutputDir   = defaultConfig.OutputDir
	TemplateDir = defaultConfig.TemplateDir

	OutputExtension = defaultConfig.OutputExtension
	TemplatePattern = defaultConfig.TemplatePattern
	DataPattern     = defaultConfig.DataPattern

	AllTemplates = defaultConfig.AllTemplates
)

// Initialize loads the configuration
func Initialize(configPath string) error {
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		var localConfig Config
		if err := yaml.NewDecoder(file).Decode(&localConfig); err != nil {
			return fmt.Errorf("failed to decode config file: %w", err)
		}

		// Override default values with local config values
		if localConfig.OutputDir != "" {
			OutputDir = localConfig.OutputDir
		}
		if localConfig.TemplateDir != "" {
			TemplateDir = localConfig.TemplateDir
		}
		if localConfig.OutputExtension != "" {
			OutputExtension = localConfig.OutputExtension
		}
		if localConfig.TemplatePattern != "" {
			TemplatePattern = localConfig.TemplatePattern
		}
		if localConfig.DataPattern != "" {
			DataPattern = localConfig.DataPattern
		}
		if localConfig.AllTemplates != "" {
			AllTemplates = localConfig.AllTemplates
		}
	}

	return nil
}
