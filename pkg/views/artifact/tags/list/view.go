package list

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Name", Width: 12},
	{Title: "Pull Time", Width: 30},
	{Title: "Push Time", Width: 30},
}

func ListTags(tags []*models.Tag) {
	var rows []table.Row
	for _, tag := range tags {

		pullTime, _ := utils.FormatCreatedTime(tag.PullTime.String())
		pushTime, _ := utils.FormatCreatedTime(tag.PushTime.String())
		rows = append(rows, table.Row{
			tag.Name,
			pullTime,
			pushTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
