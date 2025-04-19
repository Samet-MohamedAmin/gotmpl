package template

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Samet-MohamedAmin/gotmpl/pkg/config"

	"gopkg.in/yaml.v3"
)

// TemplateConfig holds configuration options for template processing
type TemplateConfig struct {
	Extension string
	Separate  bool
}

// TemplateProcessor handles the processing of Go templates
type TemplateProcessor struct {
	defaultSeparate bool
	config          TemplateConfig
	configSet       bool
}

// NewProcessor creates a new template processor with default settings
func NewProcessor(defaultSeparate bool) *TemplateProcessor {
	return &TemplateProcessor{
		defaultSeparate: defaultSeparate,
		config: TemplateConfig{
			Extension: strings.TrimPrefix(config.OutputExtension, "."),
			Separate:  defaultSeparate,
		},
		configSet: false,
	}
}

// ProcessTemplate processes a single template file
func (p *TemplateProcessor) ProcessTemplate(templatePath string, multiple bool) error {
	// Reset configuration for each template
	p.resetConfig()

	// Load the template
	tmpl, err := p.loadTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Get data file path based on template path
	dataPath := p.getDataFilePath(templatePath)
	fmt.Printf("Using data file: %s\n", dataPath)

	// Load data from YAML file
	data, err := p.loadData(dataPath)
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Process the template
	return p.executeTemplate(tmpl, data, templatePath, multiple)
}

// resetConfig resets the configuration to default values
func (p *TemplateProcessor) resetConfig() {
	p.config.Extension = strings.TrimPrefix(config.OutputExtension, ".")
	p.config.Separate = p.defaultSeparate
	p.configSet = false
}

// getDataFilePath returns the path to the data file for a template
func (p *TemplateProcessor) getDataFilePath(templatePath string) string {
	dir := filepath.Dir(templatePath)
	return filepath.Join(dir, config.DataFile)
}

// loadTemplate loads and parses a template file
func (p *TemplateProcessor) loadTemplate(templatePath string) (*template.Template, error) {
	return template.ParseFiles(templatePath)
}

// loadData loads data from a YAML file
func (p *TemplateProcessor) loadData(dataPath string) (interface{}, error) {
	file, err := os.Open(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %w", err)
	}
	defer file.Close()

	var data interface{}
	if err := yaml.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode YAML data: %w", err)
	}

	return data, nil
}

// executeTemplate executes a template with provided data
func (p *TemplateProcessor) executeTemplate(tmpl *template.Template, data interface{}, templatePath string, multiple bool) error {
	// Determine output directory
	outputDir := p.determineOutputDir(templatePath, multiple)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Debug output
	fmt.Printf("Template path: %s\n", templatePath)
	fmt.Printf("Template data: %+v\n", data)

	// Execute template to buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Debug output
	fmt.Printf("Template output:\n%s\n", buf.String())

	// Parse the template output
	return p.processTemplateOutput(&buf, outputDir)
}

// determineOutputDir determines the output directory for a template
func (p *TemplateProcessor) determineOutputDir(templatePath string, multiple bool) string {
	var outputDir string

	if multiple {
		// Get immediate directory name (template name)
		dirName := filepath.Base(filepath.Dir(templatePath))
		outputDir = filepath.Join(config.OutputDir, dirName)
	} else {
		outputDir = config.OutputDir
	}

	fmt.Printf("Output directory: %s\n", outputDir)
	return outputDir
}

// processTemplateOutput processes the output of a template execution
func (p *TemplateProcessor) processTemplateOutput(output *bytes.Buffer, outputDir string) error {
	// Create a new reader from the buffer
	reader := bufio.NewReader(output)

	// Read all lines into a slice
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read template output: %w", err)
		}

		// Add line if not empty (even without newline)
		if line != "" {
			lines = append(lines, strings.TrimRight(line, "\n"))
		}

		if err == io.EOF {
			break
		}
	}

	// Process config line if present
	startIdx := 0
	if len(lines) > 0 && strings.HasPrefix(lines[0], "# config") {
		p.parseConfigLine(lines[0])
		startIdx = 1
	}

	// Join remaining lines
	content := strings.Join(lines[startIdx:], "\n")

	// Process based on separate setting
	if p.config.Separate {
		return p.processSeparatedOutput(bytes.NewBufferString(content), outputDir)
	}

	// Process as a single file
	return p.processSingleOutput(bytes.NewBufferString(content), outputDir)
}

// parseConfigLine parses a configuration line at the top of a template output
func (p *TemplateProcessor) parseConfigLine(line string) bool {
	// Check if this is a config line
	if !strings.HasPrefix(line, "# config") {
		return false
	}

	// Parse config directives
	parts := strings.Fields(line)
	for i := 2; i < len(parts); i++ {
		part := parts[i]
		if strings.HasPrefix(part, "ext=") {
			p.config.Extension = strings.TrimPrefix(part, "ext=")
			p.configSet = true
		} else if strings.HasPrefix(part, "separate=") {
			separateValue := strings.TrimPrefix(part, "separate=")
			if separateValue == "true" {
				p.config.Separate = true
				p.configSet = true
			} else if separateValue == "false" {
				p.config.Separate = false
				p.configSet = true
			}
		}
	}

	return true
}

// processSeparatedOutput processes output as multiple files split by YAML separators
func (p *TemplateProcessor) processSeparatedOutput(output *bytes.Buffer, outputDir string) error {
	scanner := bufio.NewScanner(output)
	var contentBuffer strings.Builder
	var currentFilePath string
	fileCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Handle YAML document separator
		if trimmedLine == "---" {
			// Process previous content block if any
			if contentBuffer.Len() > 0 {
				if err := p.writeContentToFile(&contentBuffer, outputDir, fileCount, currentFilePath); err != nil {
					return err
				}
				fileCount++
				contentBuffer.Reset()
				currentFilePath = ""
			}
			continue
		}

		// Handle file directive
		if strings.HasPrefix(trimmedLine, "# file:") {
			currentFilePath = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "# file:"))
			continue
		}

		// Add content to buffer
		contentBuffer.WriteString(line + "\n")
	}

	// Process the last content block
	if contentBuffer.Len() > 0 {
		if err := p.writeContentToFile(&contentBuffer, outputDir, fileCount, currentFilePath); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// processSingleOutput processes output as a single file
func (p *TemplateProcessor) processSingleOutput(output *bytes.Buffer, outputDir string) error {
	// Determine output file name
	outputName := getOutputFileName(p.config.Extension, "", 0)
	outputPath := filepath.Join(outputDir, outputName)

	// Write content to file
	content := output.String()
	fmt.Printf("Writing content to %s:\n%s\n", outputPath, content)
	return p.writeFile(outputPath, content)
}

// writeContentToFile writes content to a file based on file directive or default naming
func (p *TemplateProcessor) writeContentToFile(content *strings.Builder, outputDir string, fileCount int, filePath string) error {
	if content.Len() == 0 {
		return nil
	}

	var outputPath string
	if filePath != "" {
		// Use the specified file path
		outputPath = filepath.Join(outputDir, filePath)

		// Create any necessary directories
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	} else {
		// Use default naming scheme
		outputName := getOutputFileName(p.config.Extension, config.DefaultPrefix, fileCount)
		outputPath = filepath.Join(outputDir, outputName)
	}

	return p.writeFile(outputPath, content.String())
}

// writeFile writes content to a file
func (p *TemplateProcessor) writeFile(path string, content string) error {
	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Write file with sync to ensure it's written to disk
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	fmt.Printf("Successfully wrote file: %s\n", path)
	return nil
}

// getOutputFileName generates a filename based on extension and count
func getOutputFileName(extension, prefix string, count int) string {
	if prefix == "" {
		prefix = config.DefaultPrefix
	}

	if count > 0 {
		if extension != "" {
			return fmt.Sprintf("%s-%02d.%s", prefix, count, extension)
		}
		return fmt.Sprintf("%s-%02d", prefix, count)
	}

	if extension != "" {
		return fmt.Sprintf("%s.%s", prefix, extension)
	}
	return prefix
}
