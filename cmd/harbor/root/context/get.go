// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package context

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// GetConfigItemCommand creates the 'harbor config get' subcommand.
func GetContextItemCommand() *cobra.Command {
	var credentialName string

	cmd := &cobra.Command{
		Use:   "get <item>",
		Short: "Get a specific config item",
		Example: `
  # Get the current credential's username
  harbor context get credentials.username

  # Get a credential's username by specifying the credential name
  harbor config get credentials.username --name admin@http://demo.goharbor.io
`,
		Long: `Get the value of a specific CLI config item.
If you specify --name, that credential (rather than the "current" one) will be used.`,
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Load config
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				return fmt.Errorf("failed to get config: %w", err)
			}

			// 2. Parse the user-supplied item path (e.g., "credentials.username")
			itemPath := strings.Split(args[0], ".")

			// 3. Get the value from the config (and track actual field segments for output)
			actualSegments := []string{}
			result, err := getValueFromConfig(config, itemPath, &actualSegments, credentialName)
			if err != nil {
				return err
			}

			// 4. Prepare the final output as a map for JSON/YAML rendering.
			canonicalPath := strings.Join(actualSegments, ".")
			output := map[string]interface{}{
				canonicalPath: result,
			}

			// 5. Determine the output format (json, yaml, etc.) and print.
			formatFlag := viper.GetString("output-format")
			switch formatFlag {
			case "json":
				data, err := json.MarshalIndent(output, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal output to JSON: %w", err)
				}
				fmt.Println(string(data))

			case "yaml", "":
				data, err := yaml.Marshal(output)
				if err != nil {
					return fmt.Errorf("failed to marshal output to YAML: %w", err)
				}
				fmt.Println(string(data))

			default:
				return fmt.Errorf("unsupported output format: %s", formatFlag)
			}

			return nil
		},
	}

	// Add a --name / -n flag to allow specifying a credential
	cmd.Flags().StringVarP(
		&credentialName,
		"name",
		"n",
		"",
		"Name of the credential to get fields from (default: the current credential)",
	)

	return cmd
}

// getValueFromConfig decides if the user requested something under "credentials"
// and if so, filters down to the *requested credential*, otherwise
// it just searches in the top-level config object.
func getValueFromConfig(
	config *utils.HarborConfig,
	path []string,
	actualSegments *[]string,
	credentialName string,
) (interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("no config item specified")
	}

	// If the first segment is "credentials", we pivot to a credential.
	if strings.EqualFold(path[0], "credentials") {
		*actualSegments = append(*actualSegments, "Credentials")

		// Determine which credential name to use
		credName := config.CurrentCredentialName
		if credentialName != "" {
			credName = credentialName
		}

		// Find the matching credential
		var targetCred *utils.Credential
		for i := range config.Credentials {
			if strings.EqualFold(config.Credentials[i].Name, credName) {
				targetCred = &config.Credentials[i]
				break
			}
		}
		if targetCred == nil {
			return nil, fmt.Errorf("no matching credential found for '%s'", credName)
		}

		// Remove "credentials" from the path, keep the rest
		return getNestedValue(*targetCred, path[1:], actualSegments)
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
		// If it's a pointer, dereference
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
