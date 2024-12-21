package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// GetConfigItemCommand creates the 'harbor config get' subcommand.
func GetConfigItemCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <item>",
		Short:   "Get a specific config item",
		Example: `  harbor config get credentials.username`,
		Long:    `Get the value of a specific CLI config item`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// 1. Load config
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				logrus.Errorf("Failed to get config: %v", err)
				return
			}

			// 2. Parse the user-supplied item path (e.g. "credentials.username")
			itemPath := strings.Split(args[0], ".")

			// 3. Get the value from the config (filtering credentials if needed)
			//    also capture the actual field name segments for pretty output
			actualSegments := []string{}
			result, err := getValueFromConfig(config, itemPath, &actualSegments)
			if err != nil {
				logrus.Error(err)
				return
			}

			// 4. Prepare the final output as a map so we can render easily in JSON/YAML
			//    We join the actual field names with "." for the "canonical" path.
			canonicalPath := strings.Join(actualSegments, ".")
			output := map[string]interface{}{
				canonicalPath: result,
			}

			// 5. Determine output format
			formatFlag := viper.GetString("output-format")
			switch formatFlag {
			case "json":
				data, err := json.MarshalIndent(output, "", "  ")
				if err != nil {
					logrus.Errorf("Failed to marshal output to JSON: %v", err)
					return
				}
				fmt.Println(string(data))

			case "yaml", "":
				data, err := yaml.Marshal(output)
				if err != nil {
					logrus.Errorf("Failed to marshal output to YAML: %v", err)
					return
				}
				fmt.Println(string(data))

			default:
				logrus.Errorf("Unsupported output format: %s", formatFlag)
			}
		},
	}

	return cmd
}

// getValueFromConfig decides if the user requested something under "credentials"
// and if so, filters down to the current credential; otherwise, it just
// searches in the top-level config object.
//
// We also accept a pointer to 'actualSegments', so that if the user typed
// "credentials.Username", we can store the correct name for each field. E.g. "Credentials" -> "Username".
func getValueFromConfig(config *utils.HarborConfig, path []string, actualSegments *[]string) (interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("no config item specified")
	}

	// If the first segment is "credentials", we pivot to the "current credential"
	// and append the actual field name "Credentials" to 'actualSegments'.
	if strings.EqualFold(path[0], "credentials") {
		*actualSegments = append(*actualSegments, "Credentials")

		// Find the current credential
		currentCredName := config.CurrentCredentialName
		var currentCred *utils.Credential
		for _, cred := range config.Credentials {
			if strings.EqualFold(cred.Name, currentCredName) {
				currentCred = &cred
				break
			}
		}
		if currentCred == nil {
			return nil, fmt.Errorf("no matching credential found for '%s'", currentCredName)
		}

		// Remove "credentials" from the path, keep the rest
		return getNestedValue(*currentCred, path[1:], actualSegments)
	}

	// Otherwise, search in the overall config struct
	return getNestedValue(*config, path, actualSegments)
}

// getNestedValue uses reflection to walk through struct fields
// (case-insensitive) according to the provided path.
//
// 'actualSegments' is updated with the actual field names as we go.
func getNestedValue(obj interface{}, path []string, actualSegments *[]string) (interface{}, error) {
	current := reflect.ValueOf(obj)

	for _, key := range path {
		if current.Kind() == reflect.Ptr {
			current = current.Elem()
		}
		if current.Kind() != reflect.Struct {
			return nil, fmt.Errorf("cannot traverse non-struct for key '%s'", key)
		}

		// Find the actual field by name, ignoring case
		var foundField reflect.StructField
		var fieldValue reflect.Value
		fieldFound := false

		t := current.Type()
		for i := 0; i < current.NumField(); i++ {
			field := t.Field(i)
			if strings.EqualFold(field.Name, key) {
				foundField = field
				fieldValue = current.Field(i)
				fieldFound = true
				break
			}
		}
		if !fieldFound {
			return nil, fmt.Errorf("config item '%s' does not exist", key)
		}

		// Record the *actual* field name in our slice
		*actualSegments = append(*actualSegments, foundField.Name)

		// Descend for the next iteration
		current = fieldValue
	}

	// Finally, if we ended on a pointer, dereference it
	if current.Kind() == reflect.Ptr {
		current = current.Elem()
	}
	return current.Interface(), nil
}
