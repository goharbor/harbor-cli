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
	{Title: "Artifact Digest", Width: 20},
	{Title: "Type", Width: 12},
	{Title: "Size", Width: 12},
	{Title: "Vulnerabilities", Width: 15},
	{Title: "Push Time", Width: 12},
}

func ListArtifacts(artifacts []*models.Artifact) {
	var rows []table.Row
	for _, artifact := range artifacts {
		pushTime, _ := utils.FormatCreatedTime(artifact.PushTime.String())
		artifactSize := utils.FormatSize(artifact.Size)
		var totalVulnerabilities int64
		for _, scan := range artifact.ScanOverview {
			totalVulnerabilities += scan.Summary.Total
		}
		rows = append(rows, table.Row{
			strconv.FormatInt(int64(artifact.ID), 10),
			artifact.Digest[:16],
			artifact.Type,
			artifactSize,
			strconv.FormatInt(totalVulnerabilities, 10),
			pushTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
