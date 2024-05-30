package member

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func RoleList(roles []string, choice chan<- int64) {
	items := make([]list.Item, len(roles))
	entityMap := make(map[string]int64, len(roles))

	for i, p := range roles {
		items[i] = selection.Item(p)
		entityMap[p] = int64(i)
	}

	m := selection.NewModel(items, "Role")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		if id, exists := entityMap[p.Choice]; exists {
			id = id + 1
			choice <- id
		} else {
			os.Exit(1)
		}
	}
}

func MemberList(member []*models.ProjectMemberEntity, choice chan<- int64) {
	items := make([]list.Item, len(member))
	entityMap := make(map[string]int64, len(member))
	for i, p := range member {
		items[i] = selection.Item(p.EntityName)
		entityMap[p.EntityName] = p.ID
	}

	m := selection.NewModel(items, "Member")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		if id, exists := entityMap[p.Choice]; exists {
			choice <- id
		} else {
			os.Exit(1)
		}
	}
}
