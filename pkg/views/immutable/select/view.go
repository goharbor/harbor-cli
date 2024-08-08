package immutable

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func ImmutableList(immutablerule []*models.ImmutableRule, choice chan<- int64) {
	itemsList := make([]list.Item, len(immutablerule))

	items := map[string]int64{}
	for i, r := range immutablerule {
		scopeSelectors := make([]string, len(r.ScopeSelectors))
		tagSelectors := make([]string, len(r.TagSelectors))
		for _, scope := range r.ScopeSelectors {
			for i, repo := range scope {
				scopeSelectors[i] = fmt.Sprintf("%s %s",repo.Decoration, repo.Pattern)
			}
		}
		for i, tag := range r.TagSelectors {
			tagSelectors[i] = fmt.Sprintf("%s %s",tag.Decoration, tag.Pattern)
		}
		for _, scope := range scopeSelectors {
			for _, tag := range tagSelectors {
				immutablename := fmt.Sprintf("for the %s, tags %s", scope, tag)
				items[immutablename] = r.ID
				itemsList[i] = selection.Item(immutablename)
			}
		}
	}

	m := selection.NewModel(itemsList, "Immutable Rule")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- items[p.Choice]
	}

}