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
package robot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

// SummaryColumns is the table layout shared by the robot list and detail views.
var SummaryColumns = []table.Column{
	{Title: "Name", Width: tablelist.WidthL},
	{Title: "Status", Width: tablelist.WidthL},
	{Title: "Permissions", Width: tablelist.WidthM},
	{Title: "Creation Time", Width: tablelist.WidthXL},
	{Title: "Expires in", Width: tablelist.WidthM},
	{Title: "Description", Width: tablelist.WidthXL},
}

// SummaryRow builds the table row shown for a robot in the list and detail
// views.
func SummaryRow(robot *models.Robot) table.Row {
	permissionCount := 0
	if len(robot.Permissions) > 0 {
		permissionCount = len(robot.Permissions[0].Access)
	}

	expires := "Never"
	if robot.ExpiresAt != -1 {
		expires = RemainingTime(robot.ExpiresAt)
	}

	createdTime, _ := utils.FormatCreatedTime(robot.CreationTime.String())
	return table.Row{
		robot.Name,
		RenderStatus(robot.Disable),
		strconv.Itoa(permissionCount),
		createdTime,
		expires,
		robot.Description,
	}
}

// RenderStatus renders a robot's enabled state for table output.
func RenderStatus(disabled bool) string {
	if disabled {
		return views.GreenStyle.Render("Disabled")
	}
	return views.GreenStyle.Render("Enabled")
}

// RemainingTime formats how much time is left until a Unix timestamp.
func RemainingTime(unixTimestamp int64) string {
	duration := time.Until(time.Unix(unixTimestamp, 0))
	if duration <= 0 {
		return "Expired"
	}

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}
