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
	"strconv"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// UpdateConfigItemCommand creates the 'harbor config update' subcommand,
// allowing you to do: harbor config update <item> <value>.
func UpdateContextItemCommand() *cobra.Command {
	var credentialName string

	cmd := &cobra.Command{
		Use:   "update <item> <value>",
		Short: "Set/update a specific config item",
		Example: `
  # Set/update the current credential's password
  harbor context update credentials.password myNewSecret

  # Set/update a credential's password by specifying the credential name
  harbor config update credentials.password myNewSecret --name admin@http://demo.goharbor.io
`,
		Long: `Set/update the value of a specific CLI config item. 
Case-insensitive field lookup, but uses the canonical (Go) field name internally.
If you specify --name, that credential (rather than the "current" one) will be updated.`,
		Args: cobra.ExactArgs(2),

		// Switch from Run to RunE so we can propagate errors
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Load the current config
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				return fmt.Errorf("failed to load Harbor config: %w", err)
			}

			// 2. Parse the user-supplied item path (e.g., "credentials.password")
			itemPath := strings.Split(args[0], ".")
			newValue := args[1]

			// 3. Reflection-based set
			actualSegments := []string{}
			if err := setValueInConfig(config, itemPath, newValue, &actualSegments, credentialName); err != nil {
				return fmt.Errorf("failed to set value in config: %w", err)
			}

			// 4. Persist the updated config to disk
			if err := utils.UpdateConfigFile(config); err != nil {
				return fmt.Errorf("failed to save updated config: %w", err)
			}

			// 5. Confirm to the user (logrus.Info is fine here; no error)
			canonicalPath := strings.Join(actualSegments, ".")
			logrus.Infof("Successfully updated %s to '%s'", canonicalPath, newValue)

			return nil
		},
	}

	// Add a --name / -n flag to allow specifying a credential
	cmd.Flags().StringVarP(
		&credentialName,
		"name",
		"n",
		"",
		"Name of the credential to set fields on (default: the current credential)",
	)

	return cmd
}

// setValueInConfig checks whether the user is updating something
// under "credentials" (i.e., a credential) or a top-level field.
//
// If path[0] == "credentials", we decide which credential to modify:
//   - If credentialName is non-empty, use that
//   - Otherwise, fallback to config.CurrentCredentialName
func setValueInConfig(
	config *utils.HarborConfig,
	path []string,
	newValue string,
	actualSegments *[]string,
	credentialName string,
) error {
	if len(path) == 0 {
		return fmt.Errorf("no config item specified")
	}

	// If the first segment is "credentials", then we pivot to a specific credential.
	if strings.EqualFold(path[0], "credentials") {
		*actualSegments = append(*actualSegments, "Credentials")

		// Determine which credential name to use
		credName := config.CurrentCredentialName
		if credentialName != "" {
			credName = credentialName
		}

		// find the matching credential
		var matchingCred *utils.Credential
		for i := range config.Credentials {
			if strings.EqualFold(config.Credentials[i].Name, credName) {
				matchingCred = &config.Credentials[i]
				break
			}
		}
		if matchingCred == nil {
			return fmt.Errorf("no matching credential found for '%s'", credName)
		}

		// Remove "credentials" from the path, and set the value in that credential
		return setNestedValue(matchingCred, path[1:], newValue, actualSegments)
	}

	// Otherwise, we set a field in the main HarborConfig struct
	return setNestedValue(config, path, newValue, actualSegments)
}

// setNestedValue navigates a pointer to a struct, following the path segments
// in a case-insensitive manner, until the last segment, where it sets the value.
//
// If the last segment is Credentials.Password, it encrypts the user-supplied
// password before storing it.
func setNestedValue(obj interface{}, path []string, newValue string, actualSegments *[]string) error {
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
			// If the field is a pointer and nil, allocate a new instance
			if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
				newElem := reflect.New(fieldValue.Type().Elem())
				fieldValue.Set(newElem)
			}
			// Descend
			val = fieldValue
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			continue
		}

		// If this is the last segment, set the value
		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set field '%s'", field.Name)
		}

		switch fieldValue.Kind() {
		case reflect.String:
			// Special case: If we are setting Credentials.Password, encrypt it
			// We'll check the last two actual segments, e.g. ["Credentials", "Password"].
			if isCredentialsPassword(*actualSegments) {
				encrypted, err := encryptPassword(newValue)
				if err != nil {
					return err
				}
				fieldValue.SetString(encrypted)
			} else {
				fieldValue.SetString(newValue)
			}

		case reflect.Bool:
			boolVal, err := strconv.ParseBool(newValue)
			if err != nil {
				return fmt.Errorf("field '%s' expects a bool, but got '%s'", field.Name, newValue)
			}
			fieldValue.SetBool(boolVal)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(newValue, 10, 64)
			if err != nil {
				return fmt.Errorf("field '%s' expects an integer, but got '%s'", field.Name, newValue)
			}
			fieldValue.SetInt(intVal)

		// If you need to handle other types (e.g. float, slice), add them here.

		default:
			return fmt.Errorf(
				"unsupported field type '%s' for field '%s'",
				fieldValue.Kind().String(), field.Name,
			)
		}
	}

	return nil
}

// isCredentialsPassword checks if the actualSegments match ["Credentials", "Password"]
// (case-insensitive).
func isCredentialsPassword(actualSegments []string) bool {
	if len(actualSegments) < 2 {
		return false
	}
	// e.g. last two items might be Credentials, Password
	last := actualSegments[len(actualSegments)-1]
	secondLast := actualSegments[len(actualSegments)-2]
	return strings.EqualFold(secondLast, "Credentials") &&
		strings.EqualFold(last, "Password")
}

// encryptPassword uses your existing utility functions to generate/retrieve a key
// and return an encrypted version of the supplied password.
func encryptPassword(plaintext string) (string, error) {
	// Make sure a key exists
	if err := utils.GenerateEncryptionKey(); err != nil {
		// It's okay if the key already exists; that might not be a fatal error for you
		logrus.Debugf("Encryption key might already exist: %v", err)
	}

	key, err := utils.GetEncryptionKey()
	if err != nil {
		return "", fmt.Errorf("failed to get encryption key: %w", err)
	}

	encrypted, err := utils.Encrypt(key, []byte(plaintext))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password: %w", err)
	}
	return encrypted, nil
}
