package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "CVE ID", Width: 12},
	{Title: "Repository Name", Width: 12},
	{Title: "Digest", Width: 12},
	{Title: "Tags", Width: 12},
	{Title: "CVSS3", Width: 5},
	{Title: "Severity", Width: 10},
	{Title: "Package", Width: 10},
	{Title: "Current version", Width: 15},
	{Title: "Fixed in version", Width: 15},
}

func ListVulnerability(vulnerability []*models.VulnerabilityItem) {
	var rows []table.Row
	for _, vul := range vulnerability {
		var tags string
		if len(tags) != 0 {
			tags = strings.Join(vul.Tags, ",")
		}
		rows = append(rows, table.Row{
			vul.CVEID,
			vul.RepositoryName,
			vul.Digest,
			tags,
			fmt.Sprintf("%.1f", vul.CvssV3Score),
			vul.Severity,
			vul.Package,
			vul.Version,
			vul.FixedVersion,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
