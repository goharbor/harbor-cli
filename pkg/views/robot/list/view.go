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
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Name", Width: tablelist.WidthL},
	{Title: "Status", Width: tablelist.WidthL},
	{Title: "Permissions", Width: tablelist.WidthM},
	{Title: "Creation Time", Width: tablelist.WidthXL},
	{Title: "Expires in", Width: tablelist.WidthM},
	{Title: "Description", Width: tablelist.WidthXL},
}

func ListRobots(robots []*models.Robot) {
	var rows []table.Row

	for _, robot := range robots {
		var enabledStatus string
		var expires string

		if robot.Disable {
			enabledStatus = views.GreenANSI + "Disabled" + views.ResetANSI
		} else {
			enabledStatus = views.GreenANSI + "Enabled" + views.ResetANSI
		}

		TotalPermissions := strconv.FormatInt(int64(len(robot.Permissions[0].Access)), 10)

		if robot.ExpiresAt == -1 {
			expires = "Never"
		} else {
			expires = remainingTime(robot.ExpiresAt)
		}

		createdTime, _ := utils.FormatCreatedTime(robot.CreationTime.String())
		rows = append(rows, table.Row{
			robot.Name,
			enabledStatus,
			TotalPermissions,
			createdTime,
			expires,
			string(robot.Description),
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
