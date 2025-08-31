package commands

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// These variables will be set during build
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoKit CLI v%s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildTime)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}
