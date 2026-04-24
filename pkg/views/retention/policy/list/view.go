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
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Policy ID", Width: tablelist.WidthM},
	{Title: "Rule", Width: tablelist.WidthS},
	{Title: "Action", Width: tablelist.WidthM},
	{Title: "Template", Width: tablelist.WidthM},
	{Title: "Repository", Width: tablelist.Width3XL},
	{Title: "Tag", Width: tablelist.Width3XL},
	{Title: "Cron", Width: tablelist.WidthL},
}

func ListPolicy(policy *models.RetentionPolicy) {
	if policy == nil || len(policy.Rules) == 0 {
		fmt.Println("No retention policy found.")
		return
	}

	cron := ""
	if policy.Trigger != nil && policy.Trigger.Settings != nil {
		switch settings := policy.Trigger.Settings.(type) {
		case map[string]interface{}:
			if value, ok := settings["cron"]; ok && value != nil {
				cron = fmt.Sprintf("%v", value)
			}
		case map[string]string:
			cron = settings["cron"]
		}
	}

	rows := make([]table.Row, 0, len(policy.Rules))
	for i, rule := range policy.Rules {
		ruleLabel := fmt.Sprintf("%d", i+1)
		if rule.ID > 0 {
			ruleLabel = fmt.Sprintf("%d", rule.ID)
		}

		repoSelectors := flattenScopeSelectors(rule.ScopeSelectors)
		tagSelectors := flattenTagSelectors(rule.TagSelectors)

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", policy.ID),
			ruleLabel,
			rule.Action,
			rule.Template,
			strings.Join(repoSelectors, ", "),
			strings.Join(tagSelectors, ", "),
			cron,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func flattenScopeSelectors(selectors map[string][]models.RetentionSelector) []string {
	if len(selectors) == 0 {
		return []string{"-"}
	}

	items := make([]string, 0)
	for key, values := range selectors {
		for _, value := range values {
			items = append(items, fmt.Sprintf("%s %s:%s", value.Decoration, key, value.Pattern))
		}
	}

	if len(items) == 0 {
		return []string{"-"}
	}

	return items
}

func flattenTagSelectors(selectors []*models.RetentionSelector) []string {
	if len(selectors) == 0 {
		return []string{"-"}
	}

	items := make([]string, 0, len(selectors))
	for _, value := range selectors {
		if value == nil {
			continue
		}
		items = append(items, fmt.Sprintf("%s %s", value.Decoration, value.Pattern))
	}

	if len(items) == 0 {
		return []string{"-"}
	}

	return items
}
