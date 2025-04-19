package template

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Samet-MohamedAmin/gotmpl/pkg/config"
)

// TemplateFinder handles locating template files in the filesystem
type TemplateFinder struct {
	rootDir string
}

// NewFinder creates a new template finder
func NewFinder(rootDir string) *TemplateFinder {
	return &TemplateFinder{
		rootDir: rootDir,
	}
}

// FindTemplates locates template files based on the provided template name
// If templateName is ALL, it finds all templates in all subdirectories
// Otherwise, it looks for a specific template in a subdirectory
func (f *TemplateFinder) FindTemplates(templateName string) ([]string, error) {
	if templateName == "" {
		return nil, fmt.Errorf("template name cannot be empty")
	}

	// Normalize the root directory path
	absRootDir, err := filepath.Abs(f.rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for root directory: %w", err)
	}

	// Check if the root directory exists
	rootInfo, err := os.Stat(absRootDir)
	if err != nil {
		return nil, fmt.Errorf("root directory not found: %w", err)
	}
	if !rootInfo.IsDir() {
		return nil, fmt.Errorf("root path is not a directory: %s", absRootDir)
	}

	fmt.Printf("Searching for templates in %s\n", absRootDir)
	fmt.Printf("Looking for template: %s\n", templateName)

	var templateFiles []string

	if templateName == config.AllTemplates {
		// Search for all templates in all subdirectories
		err := filepath.Walk(absRootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Log the error but continue walking
				fmt.Printf("Warning: skipping path %s due to error: %v\n", path, err)
				return nil
			}

			// Skip directories and non-template files
			if info.IsDir() || !isTemplateFile(info.Name()) {
				return nil
			}

			templateFiles = append(templateFiles, path)
			fmt.Printf("Found template: %s\n", path)
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("failed to search for templates: %w", err)
		}
	} else {
		// Look for a specific template
		templateDir := filepath.Join(absRootDir, templateName)
		templatePath := filepath.Join(templateDir, config.TemplateFile)

		// Check if the specific template directory exists
		if _, err := os.Stat(templateDir); err != nil {
			return nil, fmt.Errorf("template directory '%s' not found: %w", templateName, err)
		}

		// Check if the template file exists in the directory
		if _, err := os.Stat(templatePath); err != nil {
			return nil, fmt.Errorf("template file '%s' not found in directory '%s': %w",
				config.TemplateFile, templateName, err)
		}

		templateFiles = []string{templatePath}
		fmt.Printf("Found template: %s\n", templatePath)
	}

	if len(templateFiles) == 0 {
		if templateName == config.AllTemplates {
			return nil, fmt.Errorf("no template files found in %s directory", absRootDir)
		}
		return nil, fmt.Errorf("template '%s' not found in %s directory", templateName, absRootDir)
	}

	return templateFiles, nil
}

// isTemplateFile checks if a filename matches the configured template filename
func isTemplateFile(filename string) bool {
	return filename == config.TemplateFile
}

// GetDataFileForTemplate returns the path to the data file for a template
func (f *TemplateFinder) GetDataFileForTemplate(templatePath string) (string, error) {
	templateDir := filepath.Dir(templatePath)
	dataPath := filepath.Join(templateDir, config.DataFile)

	// Check if the data file exists
	if _, err := os.Stat(dataPath); err != nil {
		return "", fmt.Errorf("data file '%s' not found for template '%s': %w",
			config.DataFile, templatePath, err)
	}

	return dataPath, nil
}

// ListAllTemplates returns a list of all available templates with their paths
func (f *TemplateFinder) ListAllTemplates() (map[string]string, error) {
	templates := make(map[string]string)

	err := filepath.Walk(f.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors but continue walking
		}

		if !info.IsDir() && isTemplateFile(info.Name()) {
			// Extract the template name (directory name containing the template)
			templateDir := filepath.Dir(path)
			templateName := filepath.Base(templateDir)

			// Store with relative path from root directory
			relPath, err := filepath.Rel(f.rootDir, path)
			if err == nil {
				templates[templateName] = relPath
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	return templates, nil
}
