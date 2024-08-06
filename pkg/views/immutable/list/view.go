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
	{Title: "ID", Width: 6},
	{Title: "Repository", Width: 26},
	{Title: "Tag", Width: 24},
}

func ListImmuRules(immutable []*models.ImmutableRule) {
	var rows []table.Row
	for _, regis := range immutable {
		scopeSelectors := make([]string, len(regis.ScopeSelectors))
		for _, scope := range regis.ScopeSelectors {
			for i, repo := range scope {
				scopeSelectors[i] = fmt.Sprintf("%s %s", repo.Decoration, repo.Pattern)
			}
		}

		tagSelectors := make([]string, len(regis.TagSelectors))
		for i, tag := range regis.TagSelectors {
			tagSelectors[i] = fmt.Sprintf("%s %s",tag.Decoration, tag.Pattern)
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", regis.ID),
			strings.Join(scopeSelectors, " "),
			strings.Join(tagSelectors, " "),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}