package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	featureName string
	servicePath string
)

var AddFeatureCmd = &cobra.Command{
	Use:   "add",
	Short: "Add features to an existing service",
	Long: `Add features to an existing microservice.

Supported features:
- monitoring: Add Prometheus metrics and monitoring
- tracing: Add OpenTelemetry tracing
- caching: Add in-memory caching
- client: Add HTTP client utilities
- middleware: Add common HTTP middleware

Examples:
  gokit add monitoring --service ./user-service
  gokit add tracing --service ./payment-service
  gokit add caching --service ./notification-service`,
	RunE: runAddFeature,
}

func init() {
	AddFeatureCmd.Flags().StringVarP(&featureName, "feature", "f", "", "Feature to add (monitoring, tracing, caching, client, middleware)")
	AddFeatureCmd.Flags().StringVarP(&servicePath, "service", "s", ".", "Path to the service directory")

	AddFeatureCmd.MarkFlagRequired("feature")
}

func runAddFeature(cmd *cobra.Command, args []string) error {
	// Validate feature name
	if err := validateFeature(featureName); err != nil {
		return fmt.Errorf("invalid feature: %w", err)
	}

	// Validate service path
	if err := validateServicePath(servicePath); err != nil {
		return fmt.Errorf("invalid service path: %w", err)
	}

	// Add the feature
	if err := addFeatureToService(featureName, servicePath); err != nil {
		return fmt.Errorf("failed to add feature: %w", err)
	}

	fmt.Printf("✅ Feature '%s' added successfully to service in '%s'\n", featureName, servicePath)
	fmt.Printf("📝 Next steps:\n")
	fmt.Printf("   - Update your configuration if needed\n")
	fmt.Printf("   - Restart your service to apply changes\n")

	return nil
}

func validateFeature(feature string) error {
	validFeatures := []string{"monitoring", "tracing", "caching", "client", "middleware"}
	for _, valid := range validFeatures {
		if feature == valid {
			return nil
		}
	}
	return fmt.Errorf("feature must be one of: %s", strings.Join(validFeatures, ", "))
}

func validateServicePath(path string) error {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("service path does not exist: %s", path)
	}

	// Check if it's a directory
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return fmt.Errorf("service path is not a directory: %s", path)
	}

	// Check if it contains go.mod
	goModPath := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("service path does not contain go.mod: %s", path)
	}

	return nil
}

func addFeatureToService(feature, servicePath string) error {
	switch feature {
	case "monitoring":
		return addMonitoring(servicePath)
	case "tracing":
		return addTracing(servicePath)
	case "caching":
		return addCaching(servicePath)
	case "client":
		return addClient(servicePath)
	case "middleware":
		return addMiddleware(servicePath)
	default:
		return fmt.Errorf("unsupported feature: %s", feature)
	}
}

func addMonitoring(servicePath string) error {
	// Add monitoring dependencies to go.mod
	if err := updateGoMod(servicePath, []string{
		"github.com/prometheus/client_golang v1.17.0",
	}); err != nil {
		return err
	}

	// Copy monitoring templates from the service template
	templatePath, err := getFeatureTemplatePath(servicePath, "monitoring")
	if err != nil {
		return fmt.Errorf("failed to get monitoring template path: %w", err)
	}

	// Clean up temporary template directory after copying
	defer os.RemoveAll(templatePath)

	if err := copyFeatureTemplates(templatePath, servicePath); err != nil {
		return fmt.Errorf("failed to copy monitoring templates: %w", err)
	}

	return nil
}

func addTracing(servicePath string) error {
	// Add tracing dependencies to go.mod
	if err := updateGoMod(servicePath, []string{
		"go.opentelemetry.io/otel v1.21.0",
		"go.opentelemetry.io/otel/trace v1.21.0",
		"go.opentelemetry.io/otel/exporters/jaeger v1.21.0",
	}); err != nil {
		return err
	}

	// Copy tracing templates from the service template
	templatePath, err := getFeatureTemplatePath(servicePath, "tracing")
	if err != nil {
		return fmt.Errorf("failed to get tracing template path: %w", err)
	}

	// Clean up temporary template directory after copying
	defer os.RemoveAll(templatePath)

	if err := copyFeatureTemplates(templatePath, servicePath); err != nil {
		return fmt.Errorf("failed to copy tracing templates: %w", err)
	}

	return nil
}

func addCaching(servicePath string) error {
	// Add caching dependencies to go.mod
	if err := updateGoMod(servicePath, []string{
		"github.com/patrickmn/go-cache v2.1.0+incompatible",
	}); err != nil {
		return err
	}

	// Copy caching templates from the service template
	templatePath, err := getFeatureTemplatePath(servicePath, "caching")
	if err != nil {
		return fmt.Errorf("failed to get caching template path: %w", err)
	}

	// Clean up temporary template directory after copying
	defer os.RemoveAll(templatePath)

	if err := copyFeatureTemplates(templatePath, servicePath); err != nil {
		return fmt.Errorf("failed to copy caching templates: %w", err)
	}

	return nil
}

func addClient(servicePath string) error {
	// Copy client templates from the service template
	templatePath, err := getFeatureTemplatePath(servicePath, "client")
	if err != nil {
		return fmt.Errorf("failed to get client template path: %w", err)
	}

	// Clean up temporary template directory after copying
	defer os.RemoveAll(templatePath)

	if err := copyFeatureTemplates(templatePath, servicePath); err != nil {
		return fmt.Errorf("failed to copy client templates: %w", err)
	}

	return nil
}

func addMiddleware(servicePath string) error {
	// Copy middleware templates from the service template
	templatePath, err := getFeatureTemplatePath(servicePath, "middleware")
	if err != nil {
		return fmt.Errorf("failed to get middleware template path: %w", err)
	}

	// Clean up temporary template directory after copying
	defer os.RemoveAll(templatePath)

	if err := copyFeatureTemplates(templatePath, servicePath); err != nil {
		return fmt.Errorf("failed to copy middleware templates: %w", err)
	}

	return nil
}

func copyFeatureTemplates(templatePath, servicePath string) error {
	// Check if template path exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("feature template not found: %s", templatePath)
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
		destPath := filepath.Join(servicePath, relPath)

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

func updateGoMod(servicePath string, dependencies []string) error {
	// This is a simplified version - in a real implementation,
	// you'd parse the go.mod file and add dependencies properly

	goModPath := filepath.Join(servicePath, "go.mod")

	// Read existing go.mod
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	// Add dependencies (simplified)
	for _, dep := range dependencies {
		content = append(content, []byte("\n\t"+dep)...)
	}

	// Write back
	return os.WriteFile(goModPath, content, 0644)
}
