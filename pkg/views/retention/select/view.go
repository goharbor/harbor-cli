package retention

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func RetentionList(retentionRules []*models.RetentionRule, choice chan<- int64) {
	itemsList := make([]list.Item, 0)
	items := map[string]int64{}

	for i, rule := range retentionRules {
		tagSelectors := make([]string, len(rule.TagSelectors))
		scopeSelectors := make([]string, 0)

		for i, tag := range rule.TagSelectors {
			tagSelectors[i] = fmt.Sprintf("%s %s", tag.Decoration, tag.Pattern)
		}

		for scopeKey, selectors := range rule.ScopeSelectors {
			for _, selector := range selectors {
				scopeSelectors = append(scopeSelectors, fmt.Sprintf("%s %s (%s)", selector.Decoration, selector.Pattern, scopeKey))
			}
		}

		for _, scope := range scopeSelectors {
			for _, tag := range tagSelectors {
				display := fmt.Sprintf("Action: %s | Repo: %s | Tag: %s", rule.Action, scope, tag)
				itemsList = append(itemsList, selection.Item(display))
				items[display] = int64(i)
			}
		}
	}

	if len(itemsList) == 0 {
		fmt.Println("No retention rules found.")
		os.Exit(0)
	}

	m := selection.NewModel(itemsList, "Retention Rule")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- items[p.Choice]
	}
}
