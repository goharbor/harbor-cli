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
package api

import (
	"fmt"
	"reflect"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/configure"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func GetConfigurations() (*configure.GetConfigurationsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Configure.GetConfigurations(ctx, &configure.GetConfigurationsParams{})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetConfigurationsByCategory returns configurations filtered by category
func GetConfigurationsByCategory(category string) (map[string]interface{}, error) {
	response, err := GetConfigurations()
	if err != nil {
		return nil, err
	}

	// Validate category
	validCategories := utils.GetValidCategories()
	isValid := false
	for _, validCat := range validCategories {
		if validCat == category {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, fmt.Errorf("invalid category '%s'. Valid categories: %v", category, validCategories)
	}

	// Convert and filter
	configs := utils.ConvertToConfigurations(response.Payload)
	filteredConfigs := utils.GetConfigurationsByCategory(configs, category)

	return filteredConfigs, nil
}

// GetAllCategorizedConfigurations returns all configurations grouped by category
func GetAllCategorizedConfigurations() (map[string]map[string]interface{}, error) {
	response, err := GetConfigurations()
	if err != nil {
		return nil, err
	}

	configs := utils.ConvertToConfigurations(response.Payload)
	categorizedConfigs := utils.GetAllCategorizedConfigurations(configs)

	return categorizedConfigs, nil
}

func UpdateConfigurations(config *utils.HarborConfig) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	params := &configure.UpdateConfigurationsParams{
		Configurations: &config.Configurations,
	}

	values := reflect.ValueOf(params.Configurations).Elem()

	fmt.Println(values.NumField())
	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		fmt.Printf("Field %d: %s, Type: %s, Value: %v\n", i, values.Type().Field(i).Name, field.Type(), field.Elem())
	}

	_, err = client.Configure.UpdateConfigurations(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
