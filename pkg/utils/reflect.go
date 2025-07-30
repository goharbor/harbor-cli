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
	// Iterate through all fields in ConfigurationsResponse
	for i := 0; i < apiConfigurationsResponseObject.NumField(); i++ {
		responseObjField := apiConfigurationsResponseObject.Field(i)
		responseObjFieldName := apiConfigurationsResponseType.Field(i).Name
		// fmt.Println(responseObjFieldName, responseObjField, responseObjField.Elem().Type(), responseObjField.Type())

		targetConfigurationsField := targetConfigurationsObject.FieldByName(responseObjFieldName)
		if targetConfigurationsField.IsValid() && targetConfigurationsField.CanSet() {
			if responseObjFieldName == "OIDCClientSecret" || responseObjFieldName == "UaaClientSecret" || responseObjFieldName == "LdapSearchPassword" {
				fmt.Println("Converting secret field:", responseObjFieldName, targetConfigurationsField)
				convertSecretField(responseObjField, targetConfigurationsField)
			} else {
				convertAndSetField(responseObjField, targetConfigurationsField)
			}
		}
	}
	// DEBUG: Print the final converted configuration
	fmt.Println("\n=== FINAL CONVERTED CONFIGURATION ===")
	finalValue := reflect.ValueOf(targetConfigurationsPointer).Elem()
	finalType := finalValue.Type()

	for i := 0; i < finalValue.NumField(); i++ {
		field := finalValue.Field(i)
		fieldName := finalType.Field(i).Name

		if field.IsNil() {
			fmt.Printf("  %s: <nil>\n", fieldName)
		} else {
			fmt.Printf("  %s: %v (type: %T)\n", fieldName, field.Elem().Interface(), field.Elem().Interface())
		}
	}
	return targetConfigurationsPointer
}

// func convertAndSetField(source, target reflect.Value) {
// 	if !source.IsValid() || source.IsNil() {
// 		return
// 	}
// 	sourceObject := source.Elem()
// 	sourceObjectType := sourceObject.Type()
// 	if sourceObjectType.Kind() == reflect.Struct {

//			fmt.Println(sourceObject, "xxx", sourceObjectType)
//			if valueField := sourceObject.FieldByName("Value"); valueField.IsValid() {
//				valuePtr := reflect.New(valueField.Type())
//				valuePtr.Elem().Set(valueField)
//				target.Set(valuePtr)
//			}
//		}
//	}
func convertAndSetField(source, target reflect.Value) {
	if !source.IsValid() || source.IsNil() {
		fmt.Println("DEBUG: Source is invalid or nil")
		return
	}

	sourceObject := source.Elem()
	sourceObjectType := sourceObject.Type()

	if sourceObjectType.Kind() == reflect.Struct {
		fmt.Printf("DEBUG: Processing %v\n", sourceObjectType)

		// Get both fields for debugging
		if editableField := sourceObject.FieldByName("Editable"); editableField.IsValid() {
			fmt.Printf("DEBUG: Editable = %v\n", editableField.Interface())
		}

		if valueField := sourceObject.FieldByName("Value"); valueField.IsValid() {
			fmt.Printf("DEBUG: Value = %v (type: %v)\n", valueField.Interface(), valueField.Type())

			// Create pointer to the value
			valuePtr := reflect.New(valueField.Type())
			valuePtr.Elem().Set(valueField)

			// Debug the pointer we created
			fmt.Printf("DEBUG: Created pointer: %v (type: %v, points to: %v)\n",
				valuePtr.Interface(), valuePtr.Type(), valuePtr.Elem().Interface())

			// Debug target before setting
			fmt.Printf("DEBUG: Target before set - CanSet: %v, Type: %v, Kind: %v\n",
				target.CanSet(), target.Type(), target.Kind())

			// Set the target
			target.Set(valuePtr)

			// Debug target after setting
			fmt.Printf("DEBUG: Target after set - IsNil: %v, Type: %v, Kind: %v\n",
				target.IsNil(), target.Type(), target.Kind())

			if !target.IsNil() {
				fmt.Printf("DEBUG: Target points to: %v\n", target.Elem().Interface())
			} else {
				fmt.Printf("DEBUG: Target is nil!\n")
			}

			fmt.Printf("DEBUG: Set target to pointer of %v\n", valueField.Interface())
		} else {
			fmt.Println("DEBUG: No 'Value' field found")
		}
	} else {
		fmt.Printf("DEBUG: Source is not a struct, kind: %v\n", sourceObjectType.Kind())
	}
}

func convertSecretField(source, target reflect.Value) {
	if !source.IsValid() || source.IsNil() {
		// Set to empty string pointer instead of nil for secret fields
		emptyString := ""
		target.Set(reflect.ValueOf(&emptyString))
		fmt.Printf("DEBUG: Set secret field to empty string (source was nil)\n")
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

	// IMPORTANT: Don't encrypt for Harbor API updates - Harbor expects plain text
	// Set to empty string if no value, otherwise use the actual value
	if originalSecretValue == "" {
		emptyString := ""
		target.Set(reflect.ValueOf(&emptyString))
		fmt.Printf("DEBUG: Set secret field to empty string (value was empty)\n")
	} else {
		target.Set(reflect.ValueOf(&originalSecretValue))
		fmt.Printf("DEBUG: Set secret field to: '%s'\n", originalSecretValue)
	}
}

// func convertSecretField(source, target reflect.Value) {
// 	if !source.IsValid() || source.IsNil() {
// 		// Set to nil pointer for empty secrets
// 		target.Set(reflect.Zero(target.Type()))
// 		return
// 	}
// 	sourceElem := source.Elem()
// 	var originalSecretValue string
// 	// Extract the actual secret value from the config item
// 	if sourceElem.Type().Kind() == reflect.Struct {
// 		if valueField := sourceElem.FieldByName("Value"); valueField.IsValid() {
// 			if secretValue, ok := valueField.Interface().(string); ok {
// 				originalSecretValue = secretValue
// 			}
// 		}
// 	}
// 	// If the secret is empty, set appropriate null/empty value
// 	if originalSecretValue == "" {
// 		// Set to nil pointer (which becomes null in YAML)
// 		target.Set(reflect.Zero(target.Type()))
// 		return
// 	}
// 	// Only encrypt non-empty secrets
// 	key, err := GetEncryptionKey()
// 	if err != nil {
// 		fmt.Printf("Error getting encryption key for secret: %v\n", err)
// 		// Set to nil on encryption error
// 		target.Set(reflect.Zero(target.Type()))
// 		return
// 	}
// 	encryptedSecret, err := Encrypt(key, []byte(originalSecretValue))
// 	if err != nil {
// 		fmt.Printf("Error encrypting secret: %v\n", err)
// 		// Set to nil on encryption error
// 		target.Set(reflect.Zero(target.Type()))
// 		return
// 	}
// 	// Set the encrypted secret
// 	target.Set(reflect.ValueOf(&encryptedSecret))
// 	fmt.Printf("Secret field encrypted and stored successfully (encrypted length: %d)\n", len(encryptedSecret))
// }
