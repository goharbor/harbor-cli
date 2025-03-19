package quota

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

// Function to get project ref details
func getRefProjectName(ref models.QuotaRefObject) (string, error) {
	if refMap, ok := ref.(map[string]interface{}); ok {
		projectName, _ := refMap["name"].(string)
		return projectName, nil
	}
	return "", fmt.Errorf("Error: Ref is not of expected type")
}

func QuotaList(quotas []*models.Quota, choice chan<- int64) {
	items := make([]list.Item, len(quotas))
	entityMap := make(map[string]int64, len(quotas))

	for i, p := range quotas {
		projectName, _ := getRefProjectName(p.Ref)
		items[i] = selection.Item(fmt.Sprintf("%v", projectName))
		entityMap[projectName] = p.ID
	}

	m := selection.NewModel(items, "Quota")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		chosen := entityMap[p.Choice]
		choice <- chosen
	}
}
