package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// applyFlagOverrides recursively applies flag values to the config object
func applyFlagOverrides(cmd *cobra.Command, config interface{}, prefix string) {
	// Process the --from-env flag first if it exists (only for the root config object)
	if prefix == "" && cmd.PersistentFlags().Changed("from-env") {
		fromEnvPairs, _ := cmd.PersistentFlags().GetStringSlice("from-env")
		for _, pair := range fromEnvPairs {
			parts := strings.SplitN(pair, "::", 2)
			if len(parts) != 2 {
				fmt.Fprintf(os.Stderr, "Warning: Invalid format for --from-env flag: %s (expected 'config.path::ENV_VAR_NAME')\n", pair)
				continue
			}

			configPath := strings.TrimSpace(parts[0])
			envVarName := strings.TrimSpace(parts[1])

			// Get the environment variable value
			envValue := os.Getenv(envVarName)
			if envValue == "" {
				fmt.Fprintf(os.Stderr, "Warning: Environment variable %s is not set or empty\n", envVarName)
				continue
			}

			// Set the value in the config using dot notation path
			setValueByPath(cmd, config, configPath, envValue)
		}
	}

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

		// Process based on the type
		switch fieldValue.Kind() {
		case reflect.Ptr:
			if !fieldValue.IsNil() {
				applyFlagOverrides(cmd, fieldValue.Interface(), flagName)
			}
		case reflect.Struct:
			applyFlagOverrides(cmd, fieldValue.Addr().Interface(), flagName)
		case reflect.String:
			if cmd.PersistentFlags().Changed(flagName) {
				val, _ := cmd.PersistentFlags().GetString(flagName)
				// Check if it's an environment variable reference
				resolvedVal := resolveEnvVar(val)
				fieldValue.SetString(resolvedVal)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if cmd.PersistentFlags().Changed(flagName) {
				val, _ := cmd.PersistentFlags().GetInt64(flagName)
				fieldValue.SetInt(val)
			}
		case reflect.Bool:
			if cmd.PersistentFlags().Changed(flagName) {
				val, _ := cmd.PersistentFlags().GetBool(flagName)

				fieldValue.SetBool(val)
			}
		case reflect.Float32, reflect.Float64:
			if cmd.PersistentFlags().Changed(flagName) {
				val, _ := cmd.PersistentFlags().GetFloat64(flagName)
				fieldValue.SetFloat(val)
			}
		}
	}
}

// setValueByPath sets a configuration value using a dot notation path
func setValueByPath(_ *cobra.Command, config interface{}, path string, value string) {
	// Split the path into segments
	segments := strings.Split(path, ".")
	if len(segments) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Empty path provided\n")
		return
	}

	// Navigate to the target struct field
	current := config

	for i := 0; i < len(segments)-1; i++ {
		v := reflect.ValueOf(current)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			fmt.Fprintf(os.Stderr, "Error: Cannot navigate path %s, %s is not a struct (it's %s)\n",
				path, segments[i], v.Kind())
			return
		}

		// Find the field by name or YAML tag
		fieldName := findFieldByNameOrTag(v.Type(), segments[i])
		if fieldName == "" {
			fmt.Fprintf(os.Stderr, "Error: Field %s not found in path %s (in type %s)\n",
				segments[i], path, v.Type().Name())
			return
		}

		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			fmt.Fprintf(os.Stderr, "Error: Invalid field %s in path %s\n", segments[i], path)
			return
		}

		// Handle pointers
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				// Initialize if nil

				field.Set(reflect.New(field.Type().Elem()))
			}
			current = field.Interface()
		} else {
			// For structs, we need a pointer
			if field.Kind() == reflect.Struct {
				current = field.Addr().Interface()
			} else {
				fmt.Fprintf(os.Stderr, "Error: Field %s in path %s is not a struct or pointer (it's %s)\n",
					segments[i], path, field.Kind())
				return
			}
		}
	}

	// Now we have the parent struct, set the target field
	lastSegment := segments[len(segments)-1]
	v := reflect.ValueOf(current)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fieldName := findFieldByNameOrTag(v.Type(), lastSegment)
	if fieldName == "" {
		fmt.Fprintf(os.Stderr, "Error: Field %s not found in path %s (in type %s)\n",
			lastSegment, path, v.Type().Name())
		return
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		fmt.Fprintf(os.Stderr, "Error: Invalid field %s in path %s\n", lastSegment, path)
		return
	}

	if !field.CanSet() {
		fmt.Fprintf(os.Stderr, "Error: Cannot set field %s in path %s (unexported)\n", lastSegment, path)
		return
	}

	// Set the field value based on its type
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot convert %s to int for field %s: %v\n", value, path, err)
			return
		}
		field.SetInt(intVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot convert %s to bool for field %s: %v\n", value, path, err)
			return
		}
		field.SetBool(boolVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot convert %s to float for field %s: %v\n", value, path, err)
			return
		}
		field.SetFloat(floatVal)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported type for field %s: %v\n", path, field.Kind())
	}
}

// findFieldByNameOrTag finds a struct field by name or YAML tag, with case-insensitive matching
func findFieldByNameOrTag(t reflect.Type, name string) string {
	nameLower := strings.ToLower(name)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Check direct name match (case insensitive)
		if strings.ToLower(field.Name) == nameLower {
			return field.Name
		}

		// Check YAML tag match (case insensitive)
		yamlTag := field.Tag.Get("yaml")
		if yamlTag != "" {
			parts := strings.Split(yamlTag, ",")
			if strings.ToLower(parts[0]) == nameLower {
				return field.Name
			}
		}

		// Check JSON tag match (case insensitive)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if strings.ToLower(parts[0]) == nameLower {
				return field.Name
			}
		}
	}

	return ""
}
