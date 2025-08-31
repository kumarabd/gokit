package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	TemplateRepoURL = "https://github.com/kumarabd/service-template"
	GokitConfigFile = ".gokit.yml"
)

// GokitConfig represents the configuration stored in .gokit.yml
type GokitConfig struct {
	Initialized bool      `yaml:"initialized"`
	TemplateURL string    `yaml:"template_url"`
	CreatedAt   time.Time `yaml:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at"`
}

// cloneTemplateToProject clones the service-template repository directly into the project directory
func cloneTemplateToProject(projectDir string) error {
	templateDir := filepath.Join(projectDir, ".template")
	gokitConfigPath := filepath.Join(projectDir, GokitConfigFile)

	// Check if .gokit.yml already exists
	if _, err := os.Stat(gokitConfigPath); err == nil {
		return fmt.Errorf("project already initialized with GoKit. Remove the .gokit.yml file or use a different project location")
	}

	// Clone the template repository into the project
	cmd := exec.Command("git", "clone", TemplateRepoURL, templateDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone template repository: %w", err)
	}

	// Create .gokit.yml configuration file
	config := GokitConfig{
		Initialized: true,
		TemplateURL: TemplateRepoURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := writeGokitConfig(projectDir, config); err != nil {
		// Clean up template directory if config creation fails
		os.RemoveAll(templateDir)
		return fmt.Errorf("failed to create .gokit.yml: %w", err)
	}

	fmt.Printf("📥 Template repository cloned to %s\n", templateDir)
	fmt.Printf("📝 GoKit configuration created: %s\n", gokitConfigPath)
	return nil
}

// cloneTemplateTemporarily clones the template repository to a temporary location for feature operations
func cloneTemplateTemporarily(projectDir string) (string, error) {
	// Create a temporary directory for the template
	tempDir, err := os.MkdirTemp(projectDir, "gokit-template-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Clone the template repository into the temporary directory
	cmd := exec.Command("git", "clone", TemplateRepoURL, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to clone template repository: %w", err)
	}

	return tempDir, nil
}

// getProjectTemplatePath returns the path to the template within the current project
func getProjectTemplatePath(projectDir string) (string, error) {
	// First check if .gokit.yml exists
	gokitConfigPath := filepath.Join(projectDir, GokitConfigFile)
	if _, err := os.Stat(gokitConfigPath); os.IsNotExist(err) {
		return "", fmt.Errorf("project not initialized with GoKit. Run 'gokit new' first to initialize the project")
	}

	// Verify the configuration
	config, err := readGokitConfig(projectDir)
	if err != nil {
		return "", fmt.Errorf("failed to read GoKit configuration: %w", err)
	}

	if !config.Initialized {
		return "", fmt.Errorf("project not properly initialized with GoKit")
	}

	// For feature operations, we'll clone the template temporarily
	// This eliminates the need to keep the .template directory around
	return "", nil
}

// getFeatureTemplatePath returns the path to a specific feature template within the project
func getFeatureTemplatePath(projectDir, feature string) (string, error) {
	// Validate that the project is initialized
	_, err := getProjectTemplatePath(projectDir)
	if err != nil {
		return "", err
	}

	// Clone template temporarily for feature operations
	tempTemplateDir, err := cloneTemplateTemporarily(projectDir)
	if err != nil {
		return "", fmt.Errorf("failed to clone template for feature operation: %w", err)
	}

	featurePath := filepath.Join(tempTemplateDir, "internal", feature)

	// Check if feature path exists
	if _, err := os.Stat(featurePath); os.IsNotExist(err) {
		os.RemoveAll(tempTemplateDir)
		return "", fmt.Errorf("feature template not found: %s", featurePath)
	}

	// Return the temporary path - the caller is responsible for cleanup
	return tempTemplateDir, nil
}

// writeGokitConfig writes the GoKit configuration to .gokit.yml
func writeGokitConfig(projectDir string, config GokitConfig) error {
	configPath := filepath.Join(projectDir, GokitConfigFile)

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// readGokitConfig reads the GoKit configuration from .gokit.yml
func readGokitConfig(projectDir string) (GokitConfig, error) {
	configPath := filepath.Join(projectDir, GokitConfigFile)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return GokitConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config GokitConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return GokitConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// copyTemplateContents copies all contents from the template directory to the output directory
func copyTemplateContents(templatePath, outputDir string) error {
	// Debug: check if template path exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("template path does not exist: %s", templatePath)
	}

	return filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory and .git directory
		if path == templatePath || info.Name() == ".git" {
			return nil
		}

		// Calculate relative path from template root
		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}

		// Calculate destination path
		destPath := filepath.Join(outputDir, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, 0755)
		} else {
			// Create parent directories
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}

			// Copy file
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			return os.WriteFile(destPath, content, 0644)
		}
	})
}
