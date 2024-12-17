package list

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: 6},
	{Title: "Tag Name", Width: 20},
	{Title: "Pull Time", Width: 18},
	{Title: "Push Time", Width: 18},
}

func ListTagArtifact(artifacts []*models.Tag) {
	var rows []table.Row
	for _, tag := range artifacts {
		pushTime, _ := utils.FormatCreatedTime(tag.PushTime.String())
		pullTime, _ := utils.FormatCreatedTime(tag.PullTime.String())
		rows = append(rows, table.Row{
			strconv.FormatInt(int64(tag.ID), 10),
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
