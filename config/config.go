package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// New creates a new configuration instance without requiring any arguments
func New(configObject interface{}) (interface{}, error) {
	// Create a root command for handling flags
	cmd := &cobra.Command{
		Use:   "",
		Short: "",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	// Add config file flag
	var configFile string
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to config file")
	cmd.PersistentFlags().StringSlice("from-env", []string{}, "Set config values from environment variables in format 'config.path::ENV_VAR_NAME'")

	// Register all config flags
	registerFlags(cmd, configObject, "")

	// Parse the command line (using args from os.Args)
	cmd.SetArgs(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		return nil, err
	}

	// Get the value of the config file flag
	configFile, _ = cmd.PersistentFlags().GetString("config")

	// Load config file if specified
	if configFile != "" {
		// Read the config file
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		// Unmarshal config file data
		if err := yaml.Unmarshal(data, configObject); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
	}

	// Apply flag values that override config file
	applyFlagOverrides(cmd, configObject, "")
	return configObject, nil
}

// registerFlags recursively registers flags for all fields in the config structure
func registerFlags(cmd *cobra.Command, config interface{}, prefix string) {
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get YAML tag name or use field name
		yamlTag := field.Tag.Get("yaml")
		name := field.Name
		if yamlTag != "" {
			parts := strings.Split(yamlTag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
		}

		flagName := name
		if prefix != "" {
			flagName = prefix + "." + name
		}

		// Handle different field types
		switch fieldValue.Kind() {
		case reflect.Ptr:
			// If nil, initialize with new instance of the type
			if fieldValue.IsNil() && fieldValue.CanSet() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}
			if !fieldValue.IsNil() {
				registerFlags(cmd, fieldValue.Interface(), flagName)
			}
		case reflect.Struct:
			registerFlags(cmd, fieldValue.Addr().Interface(), flagName)
		case reflect.String:
			var value string
			if fieldValue.CanInterface() {
				value = fieldValue.String()
			}
			cmd.PersistentFlags().String(flagName, value, fmt.Sprintf("Set %s", flagName))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var value int64
			if fieldValue.CanInterface() {
				value = fieldValue.Int()
			}
			cmd.PersistentFlags().Int64(flagName, value, fmt.Sprintf("Set %s", flagName))
		case reflect.Bool:
			var value bool
			if fieldValue.CanInterface() {
				value = fieldValue.Bool()
			}
			cmd.PersistentFlags().Bool(flagName, value, fmt.Sprintf("Set %s", flagName))
		case reflect.Float32, reflect.Float64:
			var value float64
			if fieldValue.CanInterface() {
				value = fieldValue.Float()
			}
			cmd.PersistentFlags().Float64(flagName, value, fmt.Sprintf("Set %s", flagName))
		}
	}
}

// resolveEnvVar checks if the input string is an environment variable reference
// (starting with $ or ${}) and returns the environment variable value if it is.
// Otherwise, it returns the original string.
func resolveEnvVar(val string) string {
	if len(val) == 0 {
		return val
	}

	// Handle ${VAR} format
	if len(val) > 3 && val[0:2] == "${" && val[len(val)-1] == '}' {
		envVarName := val[2 : len(val)-1]
		if envValue := os.Getenv(envVarName); envValue != "" {
			return envValue
		}
		return val // Return original if env var not found
	}

	// Handle $VAR format
	if val[0] == '$' {
		envVarName := val[1:]
		if envValue := os.Getenv(envVarName); envValue != "" {
			return envValue
		}
		return val // Return original if env var not found
	}

	return val
}
