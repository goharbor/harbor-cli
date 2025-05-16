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
	"fmt"
	"reflect"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteCurrent bool

// DeleteContextItemCommand creates the 'harbor context delete' subcommand,
// allowing you to do: harbor context delete <item>
func DeleteContextItemCommand() *cobra.Command {
	var credentialName string

	cmd := &cobra.Command{
		Use:   "delete <item>",
		Short: "Delete (clear) a specific config item",
		Example: `
  # Clear the current credential's password
  harbor context delete credentials.password

  # Clear a specific credential's password using --name
  harbor context delete credentials.password --name admin@http://demo.goharbor.io

  # Clear the current credential
  harbor context delete --current
`,
		Long: `Clear the value of a specific CLI config item by setting it to its zero value.
Case-insensitive field lookup, but uses the canonical (Go) field name internally.
If you specify --name, that credential (rather than the "current" one) will be used.`,
		Args: cobra.RangeArgs(0, 1),

		// Use RunE so we can propagate errors
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Load the current config
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				return fmt.Errorf("failed to load Harbor config: %w", err)
			}

			// 1a. If --current is set, remove only the credential matching CurrentCredentialName
			if deleteCurrent {
				if len(args) > 0 {
					return fmt.Errorf("cannot specify both <item> and --current")
				}

				currentName := config.CurrentCredentialName
				found := false
				if currentName != "" {
					for i, cred := range config.Credentials {
						if strings.EqualFold(cred.Name, currentName) {
							// Remove just this credential
							config.Credentials = append(config.Credentials[:i], config.Credentials[i+1:]...)
							found = true
							break
						}
					}
					config.CurrentCredentialName = ""
				}

				if err := utils.UpdateConfigFile(config); err != nil {
					return fmt.Errorf("failed to save updated config: %w", err)
				}

				if found {
					logrus.Infof("Removed credential '%s' and cleared CurrentCredentialName", currentName)
				} else {
					logrus.Infof("No credential named '%s' found; cleared CurrentCredentialName anyway", currentName)
				}

				return nil
			}

			// If --all is NOT set, we'll perform the normal item-based delete.
			// Check we actually received an item (since now it's optional).
			if len(args) == 0 {
				return fmt.Errorf("please specify an <item> or use --current")
			}

			// 2. Parse the user-supplied item path (e.g., "credentials.password")
			itemPath := strings.Split(args[0], ".")

			// 3. Reflection-based delete (zero out)
			actualSegments := []string{}
			if err := deleteValueInConfig(config, itemPath, &actualSegments, credentialName); err != nil {
				return fmt.Errorf("failed to delete value in config: %w", err)
			}

			// 4. Persist the updated config to disk
			if err := utils.UpdateConfigFile(config); err != nil {
				return fmt.Errorf("failed to save updated config: %w", err)
			}

			// 5. Confirm to the user (no error here)
			canonicalPath := strings.Join(actualSegments, ".")
			logrus.Infof("Successfully cleared %s", canonicalPath)

			return nil
		},
	}

	// Add --name / -n to let the user pick a specific credential
	cmd.Flags().StringVarP(
		&credentialName,
		"name",
		"n",
		"",
		"Name of the credential to delete fields from (default: the current credential)",
	)

	cmd.Flags().BoolVar(
		&deleteCurrent,
		"current",
		false,
		"Remove current credentials from the config",
	)

	return cmd
}

// deleteValueInConfig checks whether the user is deleting something
// under "credentials" (i.e., *a* credential) or a top-level field.
//
// If the user says "credentials.*" AND provides --name, we'll look
// up that specific credential by name. Otherwise, we use CurrentCredentialName.
func deleteValueInConfig(
	config *utils.HarborConfig,
	path []string,
	actualSegments *[]string,
	credentialName string,
) error {
	if len(path) == 0 {
		return fmt.Errorf("no config item specified")
	}

	// If the first segment is "credentials", pivot to the chosen credential.
	if strings.EqualFold(path[0], "credentials") {
		*actualSegments = append(*actualSegments, "Credentials")

		// Figure out which credential name to use
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
			return fmt.Errorf("no matching credential found for '%s'", credName)
		}

		// Remove "credentials" from path, delete the value in that credential
		return deleteNestedValue(targetCred, path[1:], actualSegments)
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
			return fmt.Errorf("cannot set field '%s' to zero value", field.Name)
		}

		// The zero value for that field can be obtained with reflect.Zero().
		zeroVal := reflect.Zero(fieldValue.Type())
		fieldValue.Set(zeroVal)
	}

	return nil
}
