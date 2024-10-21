package info

import (
	"fmt"
	"os"
	
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"github.com/go-openapi/strfmt"
)

var columns = []table.Column{
	{Title: "Field", Width: 30},
	{Title: "Value", Width: 30},
}

func PrintSystemInfo(info *models.GeneralInfo) {
	var rows []table.Row
	currentTime := getStringFromDateTime(info.CurrentTime)

	rows = append(rows, table.Row{"Harbor Version", getString(info.HarborVersion)})
	rows = append(rows, table.Row{"Auth Mode", getString(info.AuthMode)})
	rows = append(rows, table.Row{"Primary Auth Mode", formatBool(info.PrimaryAuthMode)})
	rows = append(rows, table.Row{"Self Registration", formatBool(info.SelfRegistration)})
	rows = append(rows, table.Row{"Has CA Root", formatBool(info.HasCaRoot)})
	rows = append(rows, table.Row{"Notification Enabled", formatBool(info.NotificationEnable)})
	rows = append(rows, table.Row{"Read Only", formatBool(info.ReadOnly)})
	rows = append(rows, table.Row{"External URL", getString(info.ExternalURL)})
	rows = append(rows, table.Row{"Project Creation Restriction", getString(info.ProjectCreationRestriction)})
	rows = append(rows, table.Row{"Current Time", currentTime})
	rows = append(rows, table.Row{"Registry URL", getString(info.RegistryURL)})
	rows = append(rows, table.Row{"Registry Storage Provider", getString(info.RegistryStorageProviderName)})
	rows = append(rows, table.Row{"OIDC Provider Name", getString(info.OIDCProviderName)})
	rows = append(rows, table.Row{"Banner Message", getString(info.BannerMessage)})
	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getString(s *string) string {
	if s == nil {
		return "N/A"
	}
	if *s == "" {
		return "(empty)"
	}
	return *s
}

func getStringFromDateTime(dt *strfmt.DateTime) string {
	if dt == nil {
		return "N/A"
	}
	return dt.String()
}

func formatBool(b *bool) string {
	if b == nil {
		return "N/A"
	}
	if *b {
		return "Yes"
	}
	return "No"
}
