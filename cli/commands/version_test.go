package commands

import (
	"testing"
)

func TestVersionCmd(t *testing.T) {
	// Test that version command exists
	if VersionCmd == nil {
		t.Error("VersionCmd should not be nil")
	}

	// Test that version command has correct use
	if VersionCmd.Use != "version" {
		t.Errorf("Expected use to be 'version', got '%s'", VersionCmd.Use)
	}

	// Test that version command has correct short description
	if VersionCmd.Short != "Show version information" {
		t.Errorf("Expected short to be 'Show version information', got '%s'", VersionCmd.Short)
	}
}

func TestVersionVariables(t *testing.T) {
	// Test that version variables are set (even if default values)
	if Version == "" {
		t.Error("Version should not be empty")
	}

	if BuildTime == "" {
		t.Error("BuildTime should not be empty")
	}

	if GitCommit == "" {
		t.Error("GitCommit should not be empty")
	}
}
