package main

import (
	"fmt"
	"os"

	"github.com/kumarabd/gokit/cli/commands"
	"github.com/spf13/cobra"
)

// Version information - will be set during build
var (
	version   = "dev"
	buildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "gokit",
	Short: "GoKit - Microservice development toolkit",
	Long: `GoKit is a comprehensive toolkit for building Go microservices with standardized patterns.

Features:
- Service scaffolding with best practices
- Configuration management
- Structured logging
- Error handling
- Caching
- Monitoring and tracing
- HTTP client utilities
- Server abstractions

Examples:
  gokit new service --name user-service --template http
  gokit add monitoring --service user-service
  gokit add tracing --service user-service`,
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(commands.NewServiceCmd)
	rootCmd.AddCommand(commands.AddFeatureCmd)
	rootCmd.AddCommand(commands.VersionCmd)

	// Set version information
	rootCmd.Version = version
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
