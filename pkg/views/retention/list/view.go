package list

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: 6},
	{Title: "Action", Width: 12},
	{Title: "Disabled", Width: 10},
	{Title: "Params", Width: 20},
	{Title: "Priority", Width: 10},
	{Title: "Scope Selectors", Width: 20},
	{Title: "Tag Selectors", Width: 20},
	{Title: "Template", Width: 14},
}

func ListRetentionRules(rules []*models.RetentionRule) {
	var rows []table.Row
	for _, rule := range rules {
		params := ""
		for k, v := range rule.Params {
			params += fmt.Sprintf("%s: %v, ", k, v)
		}

		scopeSelectors := ""
		for k, v := range rule.ScopeSelectors {
			scopeSelectors += fmt.Sprintf("%s: %v, ", k, v)
		}

		tagSelectors := ""
		for _, ts := range rule.TagSelectors {
			tagSelectors += fmt.Sprintf("%v, ", ts)
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", rule.ID),
			rule.Action,
			fmt.Sprintf("%v", rule.Disabled),
			params,
			fmt.Sprintf("%d", rule.Priority),
			scopeSelectors,
			tagSelectors,
			rule.Template,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
