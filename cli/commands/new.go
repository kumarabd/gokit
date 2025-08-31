package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	serviceName string
	template    string
	outputDir   string
	force       bool
)

var NewServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Create a new microservice",
	Long: `Create a new microservice with standardized structure and best practices.

Supported templates:
- http: HTTP API service
- grpc: gRPC service
- event: Event-driven service
- worker: Background worker service

Examples:
  gokit new service --name user-service --template http
  gokit new service --name payment-service --template grpc --output ./services
  gokit new service --name notification-worker --template worker --force`,
	RunE: runNewService,
}

func init() {
	NewServiceCmd.Flags().StringVarP(&serviceName, "name", "n", "", "Service name (required)")
	NewServiceCmd.Flags().StringVarP(&template, "template", "t", "http", "Service template (http, grpc, event, worker)")
	NewServiceCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	NewServiceCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing directory")

	NewServiceCmd.MarkFlagRequired("name")
}

func runNewService(cmd *cobra.Command, args []string) error {
	// Validate service name
	if err := validateServiceName(serviceName); err != nil {
		return fmt.Errorf("invalid service name: %w", err)
	}

	// Validate template
	if err := validateTemplate(template); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// Create output directory
	serviceDir := filepath.Join(outputDir, serviceName)
	if !force {
		if _, err := os.Stat(serviceDir); err == nil {
			return fmt.Errorf("directory %s already exists. Use --force to overwrite", serviceDir)
		}
	}

	// Create the service
	if err := createService(serviceName, template, serviceDir); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	fmt.Printf("✅ Service '%s' created successfully in '%s'\n", serviceName, serviceDir)
	fmt.Printf("📁 Project structure:\n")
	fmt.Printf("   %s/\n", serviceName)
	fmt.Printf("   ├── cmd/\n")
	fmt.Printf("   ├── internal/\n")
	fmt.Printf("   ├── pkg/\n")
	fmt.Printf("   ├── go.mod\n")
	fmt.Printf("   └── README.md\n")
	fmt.Printf("\n🚀 Next steps:\n")
	fmt.Printf("   cd %s\n", serviceName)
	fmt.Printf("   go mod tidy\n")
	fmt.Printf("   go run cmd/main.go\n")

	return nil
}

func validateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	// Check for valid characters
	if strings.ContainsAny(name, " \t\n\r") {
		return fmt.Errorf("service name cannot contain whitespace")
	}

	// Check for valid Go package name
	if strings.ContainsAny(name, "!@#$%^&*()+={}[]|\\:;\"'<>?,./") {
		return fmt.Errorf("service name contains invalid characters")
	}

	return nil
}

func validateTemplate(template string) error {
	validTemplates := []string{"http", "grpc", "event", "worker"}
	for _, valid := range validTemplates {
		if template == valid {
			return nil
		}
	}
	return fmt.Errorf("template must be one of: %s", strings.Join(validTemplates, ", "))
}

func createService(name, template, outputDir string) error {
	// Create base directory structure
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Clone the template repository directly into the project
	if err := cloneTemplateToProject(outputDir); err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}

	// Copy template contents to output directory
	// Use the .template directory that was just cloned
	templatePath := filepath.Join(outputDir, ".template")
	if err := copyTemplateContents(templatePath, outputDir); err != nil {
		return err
	}

	// Customize the template based on service type
	if err := customizeTemplate(name, template, outputDir); err != nil {
		return err
	}

	// Remove the .template directory after all operations are complete
	if err := os.RemoveAll(filepath.Join(outputDir, ".template")); err != nil {
		return fmt.Errorf("failed to clean up template directory: %w", err)
	}

	return nil
}

func customizeTemplate(name, template, outputDir string) error {
	// Update go.mod with the new module name
	if err := updateServiceGoMod(name, outputDir); err != nil {
		return err
	}

	// Copy template-specific files
	if err := copyTemplateFiles(name, template, outputDir); err != nil {
		return err
	}

	// Update README.md with service-specific information
	if err := updateREADME(name, template, outputDir); err != nil {
		return err
	}

	// Update Makefile with service-specific targets
	if err := updateMakefile(name, outputDir); err != nil {
		return err
	}

	// Create template-specific configuration
	if err := createConfig(name, template, outputDir); err != nil {
		return err
	}

	return nil
}

func updateServiceGoMod(name, outputDir string) error {
	goModPath := filepath.Join(outputDir, "go.mod")

	// Read existing go.mod
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	// Replace module name
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "module ") {
			lines[i] = fmt.Sprintf("module %s", name)
			break
		}
	}

	// Add GoKit dependency if not present
	hasGoKit := false
	for _, line := range lines {
		if strings.Contains(line, "github.com/kumarabd/gokit") {
			hasGoKit = true
			break
		}
	}

	if !hasGoKit {
		// Find the require block and add GoKit
		for i, line := range lines {
			if strings.TrimSpace(line) == "require (" {
				lines = append(lines[:i+1], append([]string{"\tgithub.com/kumarabd/gokit v0.0.0"}, lines[i+1:]...)...)
				break
			}
		}
	}

	return os.WriteFile(goModPath, []byte(strings.Join(lines, "\n")), 0644)
}

func copyTemplateFiles(name, template, outputDir string) error {
	// Copy template-specific files from the service template
	// Use the .template directory that was just cloned
	templatePath := filepath.Join(outputDir, ".template", "templates", template)

	// Check if template path exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// Template doesn't exist, skip
		return nil
	}

	// Copy all files from template to service
	return filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == templatePath {
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

			// Read and process file content
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Replace placeholders in content
			processedContent := processTemplateContent(string(content), name, template)

			// Write processed content
			return os.WriteFile(destPath, []byte(processedContent), 0644)
		}
	})
}

func processTemplateContent(content, name, template string) string {
	// Replace common placeholders
	content = strings.ReplaceAll(content, "{{.ServiceName}}", name)
	content = strings.ReplaceAll(content, "{{.Template}}", template)

	// Replace service name in various formats
	content = strings.ReplaceAll(content, "{{.ServiceNameCamel}}", toCamelCase(name))
	content = strings.ReplaceAll(content, "{{.ServiceNameLower}}", strings.ToLower(name))
	content = strings.ReplaceAll(content, "{{.ServiceNameUpper}}", strings.ToUpper(name))

	return content
}

func toCamelCase(s string) string {
	// Simple camel case conversion
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = strings.Title(strings.ToLower(part))
		}
	}
	return strings.Join(parts, "")
}

func updateREADME(name, template, outputDir string) error {
	readmePath := filepath.Join(outputDir, "README.md")

	// Read existing README
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	// Replace service name and add GoKit-specific information
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "Any Service Template", name)
	contentStr = strings.ReplaceAll(contentStr, "A template for creating Golang based microservice from scratch.",
		fmt.Sprintf("A %s microservice built with GoKit.", template))

	// Add GoKit-specific sections
	goKitInfo := `

## GoKit Integration

This service is built using [GoKit](https://github.com/kumarabd/gokit), a comprehensive toolkit for building Go microservices.

### Features Included

- **Configuration Management**: YAML, environment variables, and command-line flags
- **Structured Logging**: JSON logging with zerolog integration
- **Error Handling**: Standardized error types with severity levels
- **Health Checks**: Built-in health check endpoints
- **Graceful Shutdown**: Proper service shutdown handling

### Quick Start

1. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

2. Run the service:
   ` + "```bash" + `
   go run cmd/main.go
   ` + "```" + `

3. Test the service:
   ` + "```bash" + `
   curl http://localhost:8080/health
   ` + "```" + `

### Configuration

The service can be configured using:
- YAML configuration file
- Environment variables
- Command-line flags

Example configuration:
` + "```yaml" + `
server:
  host: "0.0.0.0"
  port: 8080

log:
  format: "json"
  debug_level: "info"
` + "```" + `
`

	contentStr += goKitInfo

	return os.WriteFile(readmePath, []byte(contentStr), 0644)
}

func updateMakefile(name, outputDir string) error {
	makefilePath := filepath.Join(outputDir, "Makefile")

	// Read existing Makefile
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		return err
	}

	// Update service name in build targets
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "service", name)

	// Add GoKit-specific targets
	goKitTargets := `

# GoKit specific targets
gokit-add-monitoring:
	gokit add monitoring --service .

gokit-add-tracing:
	gokit add tracing --service .

gokit-add-caching:
	gokit add caching --service .

gokit-add-client:
	gokit add client --service .

gokit-add-middleware:
	gokit add middleware --service .
`

	contentStr += goKitTargets

	return os.WriteFile(makefilePath, []byte(contentStr), 0644)
}

func createConfig(name, template, outputDir string) error {
	// Create config directory
	configDir := filepath.Join(outputDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Create default configuration file
	configContent := fmt.Sprintf(`# Configuration for %s service

server:
  host: "0.0.0.0"
  port: 8080

log:
  format: "json"
  debug_level: "info"

# Add template-specific configuration
`, name)

	switch template {
	case "http":
		configContent += `
# HTTP service specific config
http:
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 60s
`
	case "grpc":
		configContent += `
# gRPC service specific config
grpc:
  max_concurrent_streams: 100
  max_connection_idle: 30s
  max_connection_age: 60s
`
	case "event":
		configContent += `
# Event service specific config
event:
  consumer_group: "default"
  topic: "events"
  batch_size: 100
`
	case "worker":
		configContent += `
# Worker service specific config
worker:
  concurrency: 5
  poll_interval: 30s
  max_retries: 3
`
	}

	configPath := filepath.Join(configDir, "config.yaml")
	return os.WriteFile(configPath, []byte(configContent), 0644)
}
