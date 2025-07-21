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
package utils

import (
	"fmt"
	"reflect"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

func ConvertToConfigurations(resp *models.ConfigurationsResponse) *models.Configurations {
	config := &models.Configurations{}
	configValue := reflect.ValueOf(config).Elem()
	respValue := reflect.ValueOf(resp).Elem()
	respType := respValue.Type()
	// Iterate through all fields in ConfigurationsResponse
	for i := 0; i < respValue.NumField(); i++ {
		respField := respValue.Field(i)
		respFieldName := respType.Field(i).Name
		// Find corresponding field in Configurations
		if configField := configValue.FieldByName(respFieldName); configField.IsValid() && configField.CanSet() {
			if respFieldName == "OIDCClientSecret" || respFieldName == "UaaClientSecret" {
				convertSecretField(respField, configField)
			} else {
				convertAndSetField(respField, configField)
			}
		}
	}
	return config
}

func convertAndSetField(source, target reflect.Value) {
	if !source.IsValid() || source.IsNil() {
		return
	}
	sourceElem := source.Elem()
	sourceType := sourceElem.Type()
	if sourceType.Kind() == reflect.Struct {
		if valueField := sourceElem.FieldByName("Value"); valueField.IsValid() {
			valuePtr := reflect.New(valueField.Type())
			valuePtr.Elem().Set(valueField)
			target.Set(valuePtr)
		}
	}
}

func convertSecretField(source, target reflect.Value) {
	if !source.IsValid() || source.IsNil() {
		// Set to nil pointer for empty secrets
		target.Set(reflect.Zero(target.Type()))
		return
	}
	sourceElem := source.Elem()
	var originalSecretValue string
	// Extract the actual secret value from the config item
	if sourceElem.Type().Kind() == reflect.Struct {
		if valueField := sourceElem.FieldByName("Value"); valueField.IsValid() {
			if secretValue, ok := valueField.Interface().(string); ok {
				originalSecretValue = secretValue
			}
		}
	}
	// If the secret is empty, set appropriate null/empty value
	if originalSecretValue == "" {
		// Set to nil pointer (which becomes null in YAML)
		target.Set(reflect.Zero(target.Type()))
		return
	}
	// Only encrypt non-empty secrets
	key, err := GetEncryptionKey()
	if err != nil {
		fmt.Printf("Error getting encryption key for secret: %v\n", err)
		// Set to nil on encryption error
		target.Set(reflect.Zero(target.Type()))
		return
	}
	encryptedSecret, err := Encrypt(key, []byte(originalSecretValue))
	if err != nil {
		fmt.Printf("Error encrypting secret: %v\n", err)
		// Set to nil on encryption error
		target.Set(reflect.Zero(target.Type()))
		return
	}
	// Set the encrypted secret
	target.Set(reflect.ValueOf(&encryptedSecret))
	fmt.Printf("Secret field encrypted and stored successfully (encrypted length: %d)\n", len(encryptedSecret))
}
