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

// Define the SystemInfo struct that includes both Statistic and GeneralInfo
type SystemInfo struct {
	Statistics *models.Statistic   `json:"statistics"`
	SystemInfo *models.GeneralInfo `json:"system_info"`
}

func CreateSystemInfo(generalInfo *models.GeneralInfo, stats *models.Statistic) SystemInfo {
	return SystemInfo{
		Statistics: stats,
		SystemInfo: generalInfo,
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
	var columns []table.Column

	columns = column

	createRows(info.SystemInfo, &rows)

	fmt.Println("\nSystem Info:")
	mSystem := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(mSystem).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	var rows2 []table.Row
	createRows(info.Statistics, &rows2)
	fmt.Println("\nStatistics:")

	mStats := tablelist.NewModel(columns, rows2, len(rows2))
	if _, err := tea.NewProgram(mStats).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
