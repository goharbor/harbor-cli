package summary

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Summary", Width: 25},
	{Title: "Value", Width: 15},
}

func SecuritySummary(summary *models.SecuritySummary) {
	var rows []table.Row

	rows = append(rows, table.Row{"Critical Count", fmt.Sprintf("%d",summary.CriticalCnt)})
	rows = append(rows, table.Row{"Fixable Count", fmt.Sprintf("%d",summary.FixableCnt)})
	rows = append(rows, table.Row{"High Count", fmt.Sprintf("%d",summary.HighCnt)})
	rows = append(rows, table.Row{"Low Count", fmt.Sprintf("%d",summary.LowCnt)})
	rows = append(rows, table.Row{"Medium Count", fmt.Sprintf("%d",summary.MediumCnt)})
	rows = append(rows, table.Row{"Scanned Count", fmt.Sprintf("%d",summary.ScannedCnt)})
	rows = append(rows, table.Row{"None Count", fmt.Sprintf("%d",summary.NoneCnt)})
	rows = append(rows, table.Row{"Unknown Count", fmt.Sprintf("%d",summary.UnknownCnt)})
	rows = append(rows, table.Row{"Total Artifact", fmt.Sprintf("%d",summary.TotalArtifact)})
	rows = append(rows, table.Row{"Total Vulnerability", fmt.Sprintf("%d",summary.TotalVuls)})
	
	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
