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
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Attribute", Width: tablelist.WidthL * 2},
	{Title: "Value", Width: tablelist.WidthL * 2},
}

type Volumes struct {
	Free  uint64 `json:"free,omitempty"`
	Total uint64 `json:"total,omitempty"`
}

type CLIInfoView struct {
	Username           string   `json:"username"`
	RegistryAddress    string   `json:"registry_address"`
	IsSysAdmin         bool     `json:"is_sys_admin"`
	PreviouslyLoggedIn []string `json:"previously_logged_in"`
	CLIVersion         string   `json:"cli_version"`
	OSInfo             string   `json:"os"`
}

type SystemInfo struct {
	Statistics *models.Statistic   `json:"statistics"`
	SystemInfo *models.GeneralInfo `json:"system_info"`
	VolumeInfo *Volumes            `json:"storage"`
	CLIInfo    *CLIInfoView        `json:"cli_info"`
}

func CreateSystemInfo(
	generalInfo *models.GeneralInfo,
	stats *models.Statistic,
	volumes *models.SystemInfo,
	cliinfo *api.CLIInfo,
	cliVersion string,
	osInfo string,
) SystemInfo {
	return SystemInfo{
		Statistics: stats,
		SystemInfo: generalInfo,
		VolumeInfo: &Volumes{
			Free: func() uint64 {
				if volumes.Storage != nil && len(volumes.Storage) > 0 {
					return volumes.Storage[0].Free
				}
				return 0
			}(),
			Total: func() uint64 {
				if volumes.Storage != nil && len(volumes.Storage) > 0 {
					return volumes.Storage[0].Total
				}
				return 0
			}(),
		},
		CLIInfo: &CLIInfoView{
			Username:           cliinfo.Username,
			RegistryAddress:    cliinfo.RegistryAddress,
			IsSysAdmin:         cliinfo.IsSysAdmin,
			PreviouslyLoggedIn: cliinfo.PreviouslyLoggedIn,
			CLIVersion:         cliVersion,
			OSInfo:             osInfo,
		},
	}
}

func ListInfo(info *SystemInfo) {
	renderSectionTable("System Info", info.SystemInfo)
	renderSectionTable("Statistics", info.Statistics)
	renderSectionTable("Storage", info.VolumeInfo)
	renderSectionTable("Harbor CLI Info", info.CLIInfo)
}

func renderSectionTable(title string, data interface{}) {
	var rows []table.Row
	createRows(data, &rows)
	fmt.Printf("\n  %s:\n", title)
	runTable(columns, rows)
}

func createRows(data interface{}, rows *[]table.Row) {
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		fmt.Println("Error: Expected a struct or pointer to a struct")
		return
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			field = field.Elem()
		}

		switch field.Kind() {
		case reflect.Struct:
			if field.Type() == reflect.TypeOf(strfmt.DateTime{}) {
				*rows = append(*rows, table.Row{fieldName, field.Interface().(strfmt.DateTime).String()})
			} else {
				createRows(field.Interface(), rows)
			}
		case reflect.Slice:
			sliceVal := reflect.ValueOf(field.Interface())
			for j := 0; j < sliceVal.Len(); j++ {
				item := fmt.Sprintf("%v", sliceVal.Index(j).Interface())
				if j == 0 {
					*rows = append(*rows, table.Row{fieldName, item})
				} else {
					*rows = append(*rows, table.Row{"", item})
				}
			}
		default:
			var value string
			if fieldName == "IsSysAdmin" {
				value = roleString(field.Bool())
			} else {
				value = fmt.Sprintf("%v", field.Interface())
			}
			*rows = append(*rows, table.Row{fieldName, value})
		}
	}
}

func runTable(columns []table.Column, rows []table.Row) {
	model := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running table:", err)
		os.Exit(1)
	}
}

func roleString(isSysAdmin bool) string {
	if isSysAdmin {
		return "Yes"
	}
	return "No"
}
