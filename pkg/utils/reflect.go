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
	targetConfigurationsPointer := &models.Configurations{}
	targetConfigurationsObject := reflect.ValueOf(targetConfigurationsPointer).Elem()

	apiConfigurationsResponseObject := reflect.ValueOf(resp).Elem()
	apiConfigurationsResponseType := apiConfigurationsResponseObject.Type()

	for i := 0; i < apiConfigurationsResponseObject.NumField(); i++ {
		responseObjField := apiConfigurationsResponseObject.Field(i)
		responseObjFieldName := apiConfigurationsResponseType.Field(i).Name
		targetConfigurationsField := targetConfigurationsObject.FieldByName(responseObjFieldName)

		if targetConfigurationsField.IsValid() && targetConfigurationsField.CanSet() {
			isSecretField := false
			if responseObjFieldName == "OIDCClientSecret" || responseObjFieldName == "UaaClientSecret" || responseObjFieldName == "LdapSearchPassword" {
				isSecretField = true
			}
			convertAndSetField(responseObjField, targetConfigurationsField, isSecretField)
		}
	}
	return targetConfigurationsPointer
}

func convertAndSetField(source, target reflect.Value, secret bool) {
	if !source.IsValid() || source.IsNil() {
		return
	}
	sourceObject := source.Elem()
	sourceObjectType := sourceObject.Type()
	if sourceObjectType.Kind() == reflect.Struct {
		if valueField := sourceObject.FieldByName("Value"); valueField.IsValid() {
			actualValue := valueField.Interface()
			displayValue := fmt.Sprintf("%v", actualValue)
			var finalValue any
			if displayValue != "" && secret {
				encryptedValue, err := encrypt(displayValue)
				if err != nil {
					fmt.Printf("Error encrypting field %s: %v\n", sourceObjectType.Name(), err)
					return
				}
				finalValue = encryptedValue
			} else {
				finalValue = actualValue
			}
			valuePtr := reflect.New(valueField.Type())
			valuePtr.Elem().Set(reflect.ValueOf(finalValue))
			target.Set(valuePtr)
		}
	}
}

func encrypt(originalSecretValue string) (string, error) {
	key, err := GetEncryptionKey()
	if err != nil {
		return "", fmt.Errorf("failed to get encryption key: %w", err)
	}
	encryptedSecret, err := Encrypt(key, []byte(originalSecretValue))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt secret: %w", err)
	}
	return string(encryptedSecret), nil
}
