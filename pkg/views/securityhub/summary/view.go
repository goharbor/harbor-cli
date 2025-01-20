package summary

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns_total = []table.Column{
	{Title: "Critical", Width: 10},
	{Title: "High", Width: 10},
	{Title: "Medium", Width: 10},
	{Title: "Low", Width: 10},
	{Title: "N/A", Width: 10},
	{Title: "None", Width: 10},
}

func GetTotalVulnerability(item *models.SecuritySummary) {
	var rows []table.Row
	fmt.Printf("%d total with %d fixable\n", item.TotalVuls, item.FixableCnt)
	rows = append(rows, table.Row{
		fmt.Sprintf("%d", item.CriticalCnt),
		fmt.Sprintf("%d", item.HighCnt),
		fmt.Sprintf("%d", item.MediumCnt),
		fmt.Sprintf("%d", item.LowCnt),
		fmt.Sprintf("%d", item.UnknownCnt),
		fmt.Sprintf("%d", item.NoneCnt),
	})

	m := tablelist.NewModel(columns_total, rows, len(rows)+1)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

var columns_artifact = []table.Column{
	{Title: "Repository Name", Width: 15},
	{Title: "Digest", Width: 15},
	{Title: "Critical", Width: 10},
	{Title: "High", Width: 10},
	{Title: "Medium", Width: 10},
}

func ShowMostDangerousArtifacts(item *models.SecuritySummary) {
	var rows []table.Row
	for _, artifact := range item.DangerousArtifacts {
		rows = append(rows, table.Row{
			artifact.RepositoryName,
			artifact.Digest[:16],
			fmt.Sprintf("%d", item.CriticalCnt),
			fmt.Sprintf("%d", item.HighCnt),
			fmt.Sprintf("%d", item.MediumCnt),
		})
	}

	m := tablelist.NewModel(columns_artifact, rows, len(rows)+1)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

var columns_cve = []table.Column{
	{Title: "CVE ID", Width: 15},
	{Title: "Severity", Width: 15},
	{Title: "CVSS3", Width: 5},
	{Title: "Package", Width: 15},
}

func ShowMostDangerousCVE(item *models.SecuritySummary) {
	var rows []table.Row
	for _, cve := range item.DangerousCves {
		rows = append(rows, table.Row{
			cve.CVEID,
			cve.Severity,
			fmt.Sprintf("%.1f", cve.CvssScoreV3),
			cve.Package,
		})
	}

	m := tablelist.NewModel(columns_cve, rows, len(rows)+1)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
