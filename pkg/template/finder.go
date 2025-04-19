package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Samet-MohamedAmin/gotmpl/pkg/config"
)

type Finder struct{}

func NewFinder() *Finder {
	return &Finder{}
}

func (f *Finder) FindTemplates(templateName string) ([]string, error) {
	var templateFiles []string

	fmt.Printf("Searching for templates in %s\n", config.TemplateDir)
	fmt.Printf("Looking for template: %s\n", templateName)

	if templateName == config.AllTemplates {
		// Search all directories
		err := filepath.Walk(config.TemplateDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Check if it's a template file
			if strings.HasSuffix(info.Name(), config.TemplatePattern) {
				fmt.Printf("Found template: %s\n", path)
				templateFiles = append(templateFiles, path)
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("failed to walk template directory: %w", err)
		}
	} else {
		// Look for specific template in the specified directory
		templatePath := filepath.Join(config.TemplateDir, templateName, config.TemplatePattern)
		fmt.Printf("Looking for template at: %s\n", templatePath)

		info, err := os.Stat(templatePath)
		if err != nil {
			return nil, fmt.Errorf("template '%s' not found: %w", templateName, err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("template '%s' is a directory", templateName)
		}

		templateFiles = append(templateFiles, templatePath)
	}

	if len(templateFiles) == 0 {
		if templateName == config.AllTemplates {
			return nil, fmt.Errorf("no template files found in %s directory", config.TemplateDir)
		}
		return nil, fmt.Errorf("template '%s' not found in %s directory", templateName, config.TemplateDir)
	}

	return templateFiles, nil
}
