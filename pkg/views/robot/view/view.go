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
	"slices"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	config "github.com/goharbor/harbor-cli/pkg/config/robot"
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

var projectPermissionsColumns = []table.Column{
	{Title: "Resource", Width: tablelist.WidthXL},
	{Title: "Create", Width: tablelist.WidthS},
	{Title: "Delete", Width: tablelist.WidthS},
	{Title: "List", Width: tablelist.WidthS},
	{Title: "Pull", Width: tablelist.WidthS},
	{Title: "Push", Width: tablelist.WidthS},
	{Title: "Read", Width: tablelist.WidthS},
	{Title: "Stop", Width: tablelist.WidthS},
	{Title: "Update", Width: tablelist.WidthS},
}

var systemPermissionsColumns = []table.Column{
	{Title: "Resource", Width: tablelist.WidthXL},
	{Title: "Create", Width: tablelist.WidthS},
	{Title: "Delete", Width: tablelist.WidthS},
	{Title: "List", Width: tablelist.WidthS},
	{Title: "Read", Width: tablelist.WidthS},
	{Title: "Stop", Width: tablelist.WidthS},
	{Title: "Update", Width: tablelist.WidthS},
}

var projectResourceStrings = []string{
	"Accessory", "Artifact", "Artifact Addition", "Artifact Label",
	"Export CVE", "Immutable Tag", "Label", "Log", "Member",
	"Metadata", "Notification Policy", "Preheat Policy",
	"Project", "Quota", "Repository", "Robot Account", "SBOM",
	"Scan", "Scanner", "Tag", "Tag Retention",
}

var systemResourceStrings = []string{
	"Audit Log", "Catalog", "Garbage Collection", "JobService Monitor",
	"Label", "LDAP User", "Preheat Instance", "Project", "Purge Audit",
	"Quota", "Registry", "Replication", "Replication Adapter", "Replication Policy",
	"Robot", "Scan All", "Scanner", "Security Hub", "System Volumes",
	"User", "User Group",
}

func ViewRobot(robot *models.Robot) {
	var rows []table.Row
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

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	fmt.Printf("\n%sRobot Permissions:%s\n\n", views.BoldANSI, views.ResetANSI)

	var permissionsColumns []table.Column
	var resourceStrings []string
	var systemLevel bool

	if robot.Level == "system" {
		permissionsColumns = systemPermissionsColumns
		resourceStrings = systemResourceStrings
		systemLevel = true
		fmt.Printf("%sSystem-level robot with access across projects%s\n\n", views.BoldANSI, views.ResetANSI)
	} else {
		permissionsColumns = projectPermissionsColumns
		resourceStrings = projectResourceStrings
		systemLevel = false
		fmt.Printf("%sProject-level robot for project: %s%s\n\n", views.BoldANSI, robot.Permissions[0].Namespace, views.ResetANSI)
	}

	var permissionRows []table.Row
	resActs := map[string][]string{}

	for _, perm := range robot.Permissions {
		for _, access := range perm.Access {
			resActs[access.Resource] = append(resActs[access.Resource], access.Action)
		}
	}

	perms, err := config.GetAllAvailablePermissions()
	if err != nil {
		fmt.Printf("Error fetching available permissions: %v\n", err)
		os.Exit(1)
	}
	var availablePerms map[string][]string
	if systemLevel {
		availablePerms = perms.System
	} else {
		availablePerms = perms.Project
	}

	resourceMap := make(map[string]string)
	for _, displayName := range resourceStrings {
		kebabName := utils.ToKebabCase(displayName)
		resourceMap[kebabName] = displayName
	}

	for _, displayName := range resourceStrings {
		kebabName := utils.ToKebabCase(displayName)
		if _, exists := availablePerms[kebabName]; !exists {
			continue
		}

		row := table.Row{displayName}

		var actionsToCheck []string
		if systemLevel {
			actionsToCheck = []string{"create", "delete", "list", "read", "stop", "update"}
		} else {
			actionsToCheck = []string{"create", "delete", "list", "pull", "push", "read", "stop", "update"}
		}

		for _, action := range actionsToCheck {
			if slices.Contains(availablePerms[kebabName], action) {
				actions := resActs[kebabName]
				if slices.Contains(actions, action) {
					row = append(row, "✅")
				} else {
					row = append(row, "❌")
				}
			} else {
				row = append(row, " ")
			}
		}
		permissionRows = append(permissionRows, row)
	}

	if systemLevel && len(robot.Permissions) > 1 {
		t := tablelist.NewModel(permissionsColumns, permissionRows, len(permissionRows))
		if _, err := tea.NewProgram(t).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		fmt.Printf("\n%sProject-specific Permissions:%s\n", views.BoldANSI, views.ResetANSI)

		for _, perm := range robot.Permissions {
			if perm.Kind == "project" && perm.Namespace != "/" {
				fmt.Printf("\n%sProject: %s%s\n\n", views.BoldANSI, perm.Namespace, views.ResetANSI)
				projectRows := createProjectPermissionRows(perm, perms.Project)
				pt := tablelist.NewModel(projectPermissionsColumns, projectRows, len(projectRows))
				if _, err := tea.NewProgram(pt).Run(); err != nil {
					fmt.Println("Error running program:", err)
					os.Exit(1)
				}
			}
		}
	} else {
		t := tablelist.NewModel(permissionsColumns, permissionRows, len(permissionRows))
		if _, err := tea.NewProgram(t).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	}
}

func createProjectPermissionRows(perm *models.RobotPermission, availablePerms map[string][]string) []table.Row {
	var rows []table.Row
	resActs := map[string][]string{}

	for _, access := range perm.Access {
		resActs[access.Resource] = append(resActs[access.Resource], access.Action)
	}

	for _, displayName := range projectResourceStrings {
		kebabName := utils.ToKebabCase(displayName)
		if _, exists := availablePerms[kebabName]; !exists {
			continue
		}

		row := table.Row{displayName}
		for _, action := range []string{"create", "delete", "list", "pull", "push", "read", "stop", "update"} {
			if slices.Contains(availablePerms[kebabName], action) {
				actions := resActs[kebabName]
				if slices.Contains(actions, action) {
					row = append(row, "✅")
				} else {
					row = append(row, "❌")
				}
			} else {
				row = append(row, " ")
			}
		}
		rows = append(rows, row)
	}

	return rows
}

func remainingTime(unixTimestamp int64) string {
	now := time.Now()
	expirationTime := time.Unix(unixTimestamp, 0)
	duration := expirationTime.Sub(now)
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}
