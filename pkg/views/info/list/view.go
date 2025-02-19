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
package list

import (
	"fmt"
	"os"
	"reflect"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

type Volumes struct {
	Free  uint64 `json:"free,omitempty"`
	Total uint64 `json:"total,omitempty"`
}

// Define the SystemInfo struct that includes both Statistic and GeneralInfo
type SystemInfo struct {
	Statistics *models.Statistic   `json:"statistics"`
	SystemInfo *models.GeneralInfo `json:"system_info"`
	VolumeInfo *Volumes            `json:"storage"`
}

func CreateSystemInfo(
	generalInfo *models.GeneralInfo,
	stats *models.Statistic,
	volumes *models.SystemInfo,
) SystemInfo {
	return SystemInfo{
		Statistics: stats,
		SystemInfo: generalInfo,
		VolumeInfo: &Volumes{
			Free:  volumes.Storage[0].Free,
			Total: volumes.Storage[0].Total,
		},
	}
}

func createRows(data interface{}, rows *[]table.Row) {
	val := reflect.ValueOf(data)

	// Dereference pointer if necessary
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		fmt.Println("Error: Expected a struct or a pointer to a struct")
		return
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// Skip if the field type is a struct
		if field.Kind() == reflect.Struct {
			createRows(field.Interface(), rows)
			continue
		}

		fieldName := typ.Field(i).Name
		// Initialize a string variable to store the field value
		var fieldValue string

		// Dereference pointer to access underlying value
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			field = field.Elem()
		}

		// Handle slices of pointers to structs
		if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Ptr &&
			field.Type().Elem().Elem().Kind() == reflect.Struct {
			for j := 0; j < field.Len(); j++ {
				createRows(field.Index(j).Interface(), rows)
			}
			continue
		}
		// Convert field value to string
		switch field.Kind() {
		case reflect.Struct:
			// Check if the field is of type strfmt.DateTime
			if field.Type() == reflect.TypeOf(strfmt.DateTime{}) {
				// Convert strfmt.DateTime to string
				timeStr := field.Interface().(strfmt.DateTime).String()
				fieldValue = timeStr
			} else {
				// Recursively print the struct fields
				createRows(field.Interface(), rows)
			}
		default:
			fieldValue = fmt.Sprintf("%v", field.Interface())
		}
		// Append field name and value to the rows slice
		*rows = append(*rows, table.Row{fieldName, fieldValue})
	}
}

var column = []table.Column{
	{Title: "Attribute", Width: 24},
	{Title: "Value", Width: 22},
}

func ListInfo(info *SystemInfo) {
	var rows []table.Row
	columns := column

	// Create SystemInfo Table
	createRows(info.SystemInfo, &rows)
	fmt.Println("\n  System Info:")
	mSystem := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(mSystem).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// Create Statistics Table
	var rows2 []table.Row
	createRows(info.Statistics, &rows2)
	fmt.Println("\n  Statistics:")
	mStats := tablelist.NewModel(columns, rows2, len(rows2))
	if _, err := tea.NewProgram(mStats).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// Create Storage Table
	var rows3 []table.Row
	createRows(info.VolumeInfo, &rows3)
	fmt.Println("\n  Storage:")
	mStorage := tablelist.NewModel(columns, rows3, len(rows3))
	if _, err := tea.NewProgram(mStorage).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
