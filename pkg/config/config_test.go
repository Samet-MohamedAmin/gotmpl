package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitialize(t *testing.T) {
	// Setup temporary test directory
	tempDir, err := os.MkdirTemp("", "gotmpl-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config file
	configPath := filepath.Join(tempDir, "test-config.yaml")
	configContent := `
OutputDir: "test-output"
OutputExtension: ".test"
TemplateFile: "test-template.go.tmpl"
DataFile: "test-data.yaml"
DefaultPrefix: "test-file"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Reset the configuration before test
	Reset()

	// Test initializing with the test config
	if err := Initialize(configPath); err != nil {
		t.Errorf("Initialize failed: %v", err)
	}

	// Check if values were loaded correctly
	if OutputDir != "test-output" {
		t.Errorf("Expected OutputDir to be 'test-output', got %s", OutputDir)
	}
	if OutputExtension != ".test" {
		t.Errorf("Expected OutputExtension to be '.test', got %s", OutputExtension)
	}
	if TemplateFile != "test-template.go.tmpl" {
		t.Errorf("Expected TemplateFile to be 'test-template.go.tmpl', got %s", TemplateFile)
	}
	if DataFile != "test-data.yaml" {
		t.Errorf("Expected DataFile to be 'test-data.yaml', got %s", DataFile)
	}
	if DefaultPrefix != "test-file" {
		t.Errorf("Expected DefaultPrefix to be 'test-file', got %s", DefaultPrefix)
	}

	// Test the singleton instance
	config := GetConfig()
	if config.OutputDir != "test-output" {
		t.Errorf("Expected config.OutputDir to be 'test-output', got %s", config.OutputDir)
	}
}

func TestInitializeWithNonExistentFile(t *testing.T) {
	// Reset the configuration before test
	Reset()

	// Test initializing with a non-existent config file (should use defaults)
	nonExistentPath := "non-existent-config.yaml"
	if err := Initialize(nonExistentPath); err != nil {
		t.Errorf("Initialize with non-existent file failed: %v", err)
	}

	// Check if default values were used
	if OutputDir != defaultConfig.OutputDir {
		t.Errorf("Expected default OutputDir, got %s", OutputDir)
	}
	if OutputExtension != defaultConfig.OutputExtension {
		t.Errorf("Expected default OutputExtension, got %s", OutputExtension)
	}
	if TemplateFile != defaultConfig.TemplateFile {
		t.Errorf("Expected default TemplateFile, got %s", TemplateFile)
	}
	if DataFile != defaultConfig.DataFile {
		t.Errorf("Expected default DataFile, got %s", DataFile)
	}
	if DefaultPrefix != defaultConfig.DefaultPrefix {
		t.Errorf("Expected default DefaultPrefix, got %s", DefaultPrefix)
	}
}

func TestGetOutputPath(t *testing.T) {
	// Reset configuration
	Reset()

	// Test with different configurations
	testCases := []struct {
		name           string
		outputDir      string
		outputExt      string
		baseName       string
		expectedOutput string
	}{
		{
			name:           "Default extension",
			outputDir:      "output",
			outputExt:      ".yaml",
			baseName:       "test",
			expectedOutput: filepath.Join("output", "test.yaml"),
		},
		{
			name:           "No extension",
			outputDir:      "output",
			outputExt:      "",
			baseName:       "test",
			expectedOutput: filepath.Join("output", "test"),
		},
		{
			name:           "Extension without dot",
			outputDir:      "output",
			outputExt:      "txt",
			baseName:       "test",
			expectedOutput: filepath.Join("output", "test.txt"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &AppConfig{
				OutputDir:       tc.outputDir,
				OutputExtension: tc.outputExt,
			}

			result := config.GetOutputPath(tc.baseName)
			if result != tc.expectedOutput {
				t.Errorf("Expected path: %s, got: %s", tc.expectedOutput, result)
			}
		})
	}
}
