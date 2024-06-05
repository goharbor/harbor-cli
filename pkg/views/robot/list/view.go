package list

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: 4},
	{Title: "Name", Width: 30},
	{Title: "Status", Width: 18},
	{Title: "Permissions", Width: 12},
	{Title: "Creation Time", Width: 16},
	{Title: "Expires in", Width: 14},
	{Title: "Description", Width: 12},
}

func ListRobots(robots []*models.Robot) {
	var rows []table.Row

	for _, robot := range robots {
		var enabledStatus string
		var expires string

		if robot.Disable {
			enabledStatus = views.RedStyle.Render("Disabled")
		} else {
			enabledStatus = views.GreenStyle.Render("Enabled")
		}

		TotalPermissions := strconv.FormatInt(int64(len(robot.Permissions[0].Access)), 10)

		if robot.ExpiresAt == -1 {
			expires = "Never"
		} else {
			expires = remainingTime(robot.ExpiresAt)
		}

		createdTime, _ := utils.FormatCreatedTime(robot.CreationTime.String())
		rows = append(rows, table.Row{
			strconv.FormatInt(robot.ID, 10),
			robot.Name,
			enabledStatus,
			TotalPermissions,
			createdTime,
			expires,
			robot.Description,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func remainingTime(unixTimestamp int64) string {
	// Get the current time
	now := time.Now()
	// Convert the Unix timestamp to time.Time
	expirationTime := time.Unix(unixTimestamp, 0)
	// Calculate the duration between now and the expiration time
	duration := expirationTime.Sub(now)

	// Calculate days, hours, minutes, and seconds
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	// Format the output string
	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}

func getStatusStyle(status string) lipgloss.Style {
	statusStyle := views.RedStyle
	if status == "healthy" {
		statusStyle = views.GreenStyle
	}
	return statusStyle
}
