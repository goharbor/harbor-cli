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
package view

import (
	"fmt"
	"os"
	"reflect"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Property", Width: tablelist.Width3XL},
	{Title: "Editable", Width: tablelist.WidthM},
	{Title: "Value", Width: tablelist.Width3XL},
}

func ViewConfigurations(configs *models.ConfigurationsResponse, category string) {
	var rows []table.Row
	apiConfigurationsResponseObject := reflect.ValueOf(configs).Elem()
	apiConfigurationsResponseType := apiConfigurationsResponseObject.Type()
	for i := 0; i < apiConfigurationsResponseObject.NumField(); i++ {
		fieldItem := apiConfigurationsResponseObject.Field(i).Elem()
		valueField := fieldItem.FieldByName("Value")
		var displayValue string
		if valueField.IsValid() {
			actualValue := valueField.Interface()
			displayValue = fmt.Sprintf("%v", actualValue)
		} else {
			displayValue = "<no value>"
		}
		editableField := fieldItem.FieldByName("Editable")
		var displayEditable bool
		if editableField.IsValid() {
			displayEditable = editableField.Interface().(bool)
		} else {
			displayEditable = false
		}
		if utils.IsCategory(apiConfigurationsResponseType.Field(i).Name, category) {
			rows = append(rows, table.Row{
				apiConfigurationsResponseType.Field(i).Name,
				fmt.Sprintf("%t", displayEditable),
				displayValue,
			})
		}
	}
	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
