package template

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/Samet-MohamedAmin/gotmpl/pkg/config"

	"gopkg.in/yaml.v3"
)

type Processor struct {
	separate         bool
	separateSet      bool
	originalSeparate bool
}

func NewProcessor(separate bool) *Processor {
	return &Processor{
		separate:         separate,
		separateSet:      false,
		originalSeparate: separate,
	}
}

// parseConfigLine parses a config line and updates extension and processor state
// Returns true if this was a config line, false otherwise
func (p *Processor) parseConfigLine(line string) (string, bool) {
	// Start with default extension
	ext := strings.TrimPrefix(config.OutputExtension, ".")

	// Not a config line
	if !strings.HasPrefix(line, "# config") {
		return ext, false
	}

	// Parse config directive
	parts := strings.Fields(line)
	for i := 2; i < len(parts); i++ {
		part := parts[i]
		if strings.HasPrefix(part, "ext=") {
			ext = strings.TrimPrefix(part, "ext=")
		} else if strings.HasPrefix(part, "separate=") {
			separateValue := strings.TrimPrefix(part, "separate=")
			if separateValue == "true" {
				p.separate = true
				p.separateSet = true
			} else if separateValue == "false" {
				p.separate = false
				p.separateSet = true
			}
		}
	}

	return ext, true
}

// getOutputPath determines the output path for a file based on the template name and optional file directive
func (p *Processor) getOutputPath(outputDir, templateName string, fileCount int, fileDirective string, ext string) (string, error) {
	var outputPath string
	if fileDirective != "" {
		// Use the specified file path
		outputPath = filepath.Join(outputDir, fileDirective)
		// Create any necessary subdirectories
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return "", fmt.Errorf("failed to create output directory: %w", err)
		}
	} else {
		// Use default naming with configured extension
		var outputName string
		if ext != "" {
			outputName = fmt.Sprintf("%s-%02d.%s", templateName, fileCount, ext)
		} else {
			outputName = fmt.Sprintf("%s-%02d", templateName, fileCount)
		}
		outputPath = filepath.Join(outputDir, outputName)
	}
	return outputPath, nil
}

// writeContent writes the content to the specified output path
func (p *Processor) writeContent(outputPath string, content string) error {
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	fmt.Printf("Generated output file: %s\n", outputPath)
	return nil
}

// processSeparatedContent handles the case when -separate is true
func (p *Processor) processSeparatedContent(content string, outputDir, templateName string) error {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentContent strings.Builder
	var currentFilePath string
	fileCount := 0
	processingContent := false
	ext := strings.TrimPrefix(config.OutputExtension, ".") // Default extension

	// Check for config directive at the start
	if scanner.Scan() {
		firstLine := strings.TrimSpace(scanner.Text())
		var isConfig bool
		ext, isConfig = p.parseConfigLine(firstLine)

		// If not a config line, process it as content
		if !isConfig {
			line := scanner.Text()
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, "# file:") {
				currentFilePath = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "# file:"))
			} else {
				currentContent.WriteString(line + "\n")
			}
		} else if !scanner.Scan() {
			// If it was a config line and there's no more content
			return nil
		}
	}

	// Process content before the first separator
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "---" {
			// Save any content before the first separator
			if currentContent.Len() > 0 {
				if err := p.processCurrentContent(&currentContent, outputDir, templateName, fileCount, currentFilePath, ext); err != nil {
					return err
				}
				fileCount++
				currentContent.Reset()
				currentFilePath = ""
			}
			processingContent = true
			break
		}

		if strings.HasPrefix(trimmedLine, "# file:") {
			currentFilePath = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "# file:"))
		} else {
			currentContent.WriteString(line + "\n")
		}
	}

	// Process content between separators
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "---" {
			if processingContent {
				// We've reached the end of a content block
				if err := p.processCurrentContent(&currentContent, outputDir, templateName, fileCount, currentFilePath, ext); err != nil {
					return err
				}
				fileCount++
				currentContent.Reset()
				currentFilePath = ""
			}
			processingContent = true
		} else if processingContent {
			if strings.HasPrefix(trimmedLine, "# file:") {
				currentFilePath = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "# file:"))
			} else {
				currentContent.WriteString(line + "\n")
			}
		}
	}

	// Process the last content block if there's any remaining
	if currentContent.Len() > 0 {
		if err := p.processCurrentContent(&currentContent, outputDir, templateName, fileCount, currentFilePath, ext); err != nil {
			return err
		}
	}

	return nil
}

// processCurrentContent handles writing a single content block to a file
func (p *Processor) processCurrentContent(content *strings.Builder, outputDir, templateName string, fileCount int, fileDirective string, ext string) error {
	if content.Len() == 0 {
		return nil
	}

	outputPath, err := p.getOutputPath(outputDir, templateName, fileCount, fileDirective, ext)
	if err != nil {
		return err
	}

	return p.writeContent(outputPath, content.String())
}

func (p *Processor) ProcessTemplate(templatePath string) error {
	fmt.Printf("Processing template: %s\n", templatePath)

	// Reset separate value to original command line flag value
	p.separate = p.originalSeparate
	p.separateSet = false

	dataPath := p.getDataPath(templatePath)
	fmt.Printf("Looking for data file: %s\n", dataPath)

	tmpl, err := p.loadTemplate(templatePath)
	if err != nil {
		return err
	}

	data, err := p.loadData(dataPath)
	if err != nil {
		return err
	}

	return p.executeTemplate(tmpl, data, templatePath)
}

func (p *Processor) getDataPath(templatePath string) string {
	// Get the relative path from templates directory
	relPath, err := filepath.Rel(config.TemplateDir, templatePath)
	if err != nil {
		fmt.Printf("Error getting relative path: %v\n", err)
		return ""
	}

	// Get the directory path
	dirPath := filepath.Dir(relPath)

	// Construct data path using data.yaml in the same directory
	dataPath := filepath.Join(config.TemplateDir, dirPath, config.DataPattern)
	fmt.Printf("Data path constructed: %s\n", dataPath)
	return dataPath
}

func (p *Processor) loadTemplate(templatePath string) (*template.Template, error) {
	return template.ParseFiles(templatePath)
}

func (p *Processor) loadData(dataPath string) (interface{}, error) {
	file, err := os.Open(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %w", err)
	}
	defer file.Close()

	var data interface{}
	if err := yaml.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}
	return data, nil
}

func (p *Processor) executeTemplate(tmpl *template.Template, data interface{}, templatePath string) error {
	// Get the relative path from templates directory
	relPath, err := filepath.Rel(config.TemplateDir, templatePath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	// Get the directory path and template name
	dirPath := filepath.Dir(relPath)
	templateName := filepath.Base(dirPath)

	fmt.Printf("Output directory path: %s\n", dirPath)
	fmt.Printf("Template name: %s\n", templateName)

	// Create output directory structure
	outputDir := filepath.Join(config.OutputDir, dirPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// First, execute the template to a temporary buffer
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Read the first line to check for config directive
	scanner := bufio.NewScanner(strings.NewReader(buf.String()))
	ext := strings.TrimPrefix(config.OutputExtension, ".")
	if scanner.Scan() {
		firstLine := strings.TrimSpace(scanner.Text())
		ext, _ = p.parseConfigLine(firstLine)
	}

	// If separate was set by config, use that value, otherwise use the command line flag
	shouldSeparate := p.separate
	if p.separateSet {
		shouldSeparate = p.separate
	}

	if shouldSeparate {
		return p.processSeparatedContent(buf.String(), outputDir, templateName)
	} else {
		// Write the entire content to a single file
		var outputName string
		if ext != "" {
			outputName = fmt.Sprintf("%s.%s", templateName, ext)
		} else {
			outputName = templateName
		}
		outputPath := filepath.Join(outputDir, outputName)
		return p.writeContent(outputPath, buf.String())
	}
}
