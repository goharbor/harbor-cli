package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DeleteConfigItemCommand creates the 'harbor config delete' subcommand,
// allowing you to do: harbor config delete <item>
func DeleteConfigItemCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <item>",
		Short:   "Delete (clear) a specific config item",
		Example: "  harbor config delete credentials.password",
		Long: `Clear the value of a specific CLI config item by setting it to its zero value.
Case-insensitive field lookup, but uses the canonical (Go) field name internally.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// 1. Load the current config
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				logrus.Errorf("Failed to load Harbor config: %v", err)
				return
			}

			// 2. Parse the user-supplied item path (e.g., "credentials.password")
			itemPath := strings.Split(args[0], ".")

			// 3. Reflection-based delete (zero out)
			actualSegments := []string{}
			if err := deleteValueInConfig(config, itemPath, &actualSegments); err != nil {
				logrus.Error(err)
				return
			}

			// 4. Persist the updated config to disk
			if err := utils.UpdateConfigFile(config); err != nil {
				logrus.Errorf("Failed to save updated config: %v", err)
				return
			}

			// 5. Confirm to the user
			canonicalPath := strings.Join(actualSegments, ".")
			logrus.Infof("Successfully cleared %s", canonicalPath)
		},
	}

	return cmd
}

// deleteValueInConfig checks whether the user is deleting something
// under "credentials" (i.e., the current credential) or a top-level field.
func deleteValueInConfig(config *utils.HarborConfig, path []string, actualSegments *[]string) error {
	if len(path) == 0 {
		return fmt.Errorf("no config item specified")
	}

	// If the first segment is "credentials", then we pivot to the current credential.
	if strings.EqualFold(path[0], "credentials") {
		*actualSegments = append(*actualSegments, "Credentials")

		// find the current credential
		currentCredName := config.CurrentCredentialName
		var currentCred *utils.Credential
		for i := range config.Credentials {
			if strings.EqualFold(config.Credentials[i].Name, currentCredName) {
				currentCred = &config.Credentials[i]
				break
			}
		}
		if currentCred == nil {
			return fmt.Errorf("no matching credential found for '%s'", currentCredName)
		}

		// Remove "credentials" from the path, and delete (zero) the value in that credential
		return deleteNestedValue(currentCred, path[1:], actualSegments)
	}

	// Otherwise, we delete a field in the main HarborConfig struct
	return deleteNestedValue(config, path, actualSegments)
}

// deleteNestedValue navigates a pointer to a struct, following the path segments
// in a case-insensitive manner, until the last segment, where it sets the field
// to its zero value.
func deleteNestedValue(obj interface{}, path []string, actualSegments *[]string) error {
	// We require obj to be a pointer to a struct so we can modify it.
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("object must be a pointer to a struct, got %s", val.Kind())
	}
	val = val.Elem() // dereference pointer

	for i, segment := range path {
		if val.Kind() != reflect.Struct {
			return fmt.Errorf("cannot traverse non-struct for segment '%s'", segment)
		}
		t := val.Type()

		// Case-insensitive field lookup
		fieldIndex := -1
		for j := 0; j < val.NumField(); j++ {
			if strings.EqualFold(t.Field(j).Name, segment) {
				fieldIndex = j
				break
			}
		}
		if fieldIndex < 0 {
			return fmt.Errorf("config item '%s' does not exist", segment)
		}

		field := t.Field(fieldIndex)
		fieldValue := val.Field(fieldIndex)

		// Record the actual field name
		*actualSegments = append(*actualSegments, field.Name)

		// If this is NOT the last path segment, move deeper
		if i < len(path)-1 {
			// If the field is a pointer and nil, we can't go deeper
			if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
				return fmt.Errorf("field '%s' is nil and cannot be traversed", field.Name)
			}
			// Descend
			val = fieldValue
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			continue
		}

		// If this is the last segment, set the field to zero value
		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot delete (set zero value) for field '%s'", field.Name)
		}

		// The "zero" value for that field can be obtained with reflect.Zero().
		zeroVal := reflect.Zero(fieldValue.Type())
		fieldValue.Set(zeroVal)
	}

	return nil
}
