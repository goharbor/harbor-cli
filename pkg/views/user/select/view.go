package user

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views"
)

func UserList(users []*models.UserResp, choice chan<- int64) {
	itemsList := make([]list.Item, len(users))

	items := map[string]int64{}

	for i, r := range users {
		items[r.Username] = r.UserID
		itemsList[i] = views.Item(r.Username)
	}

	m := views.NewModel(itemsList, "User")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(views.Model); ok {
		choice <- items[p.Choice]
	}

}
