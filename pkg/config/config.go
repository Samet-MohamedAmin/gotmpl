package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Version is the application version
const Version = "v0.1.2"

// AllTemplates represents a wildcard for template selection
const AllTemplates = "ALL"

// AppConfig manages the application configuration
type AppConfig struct {
	// Directories
	OutputDir string `yaml:"OutputDir"`

	// File patterns
	OutputExtension string `yaml:"OutputExtension"`
	TemplateFile    string `yaml:"TemplateFile"`
	DataFile        string `yaml:"DataFile"`
	DefaultPrefix   string `yaml:"DefaultPrefix"`
}

// Default configuration values
var defaultConfig = AppConfig{
	OutputDir:       "output",
	OutputExtension: "",
	TemplateFile:    "template.go.tmpl",
	DataFile:        "data.yaml",
	DefaultPrefix:   "file",
}

// Global instance
var (
	instance *AppConfig
	once     sync.Once
)

// Global accessor methods (for backward compatibility)
var (
	OutputDir       = defaultConfig.OutputDir
	OutputExtension = defaultConfig.OutputExtension
	TemplateFile    = defaultConfig.TemplateFile
	DataFile        = defaultConfig.DataFile
	DefaultPrefix   = defaultConfig.DefaultPrefix
)

// GetConfig returns the singleton config instance
func GetConfig() *AppConfig {
	once.Do(func() {
		if instance == nil {
			instance = &defaultConfig
		}
	})
	return instance
}

// Initialize loads the configuration from a file
func Initialize(configPath string) error {
	// Load default configuration first
	config := defaultConfig

	// If config file exists, try to load it
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		var fileConfig AppConfig
		if err := yaml.NewDecoder(file).Decode(&fileConfig); err != nil {
			return fmt.Errorf("failed to decode config file: %w", err)
		}

		// Override default values with values from config file
		if fileConfig.OutputDir != "" {
			config.OutputDir = fileConfig.OutputDir
		}
		if fileConfig.OutputExtension != "" {
			config.OutputExtension = fileConfig.OutputExtension
		}
		if fileConfig.TemplateFile != "" {
			config.TemplateFile = fileConfig.TemplateFile
		}
		if fileConfig.DataFile != "" {
			config.DataFile = fileConfig.DataFile
		}
		if fileConfig.DefaultPrefix != "" {
			config.DefaultPrefix = fileConfig.DefaultPrefix
		}
	} else if !os.IsNotExist(err) {
		// If there's an error other than "file not exists"
		return fmt.Errorf("error checking config file: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := ensureDirectory(config.OutputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Set the global instance
	once.Do(func() {
		instance = &config
	})

	// For backward compatibility with global vars
	OutputDir = config.OutputDir
	OutputExtension = config.OutputExtension
	TemplateFile = config.TemplateFile
	DataFile = config.DataFile
	DefaultPrefix = config.DefaultPrefix

	return nil
}

// ensureDirectory creates a directory if it doesn't exist
func ensureDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// GetOutputPath constructs the full output path
func (c *AppConfig) GetOutputPath(baseName string) string {
	ext := c.OutputExtension
	if len(ext) > 0 && ext[0] != '.' {
		ext = "." + ext
	}
	return filepath.Join(c.OutputDir, baseName+ext)
}

// Reset resets the configuration to default values (useful for testing)
func Reset() {
	OutputDir = defaultConfig.OutputDir
	OutputExtension = defaultConfig.OutputExtension
	TemplateFile = defaultConfig.TemplateFile
	DataFile = defaultConfig.DataFile
	DefaultPrefix = defaultConfig.DefaultPrefix
	instance = nil
}
