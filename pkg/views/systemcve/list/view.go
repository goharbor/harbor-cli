package list

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: 6},
	{Title: "CVE Name", Width: 18},
	{Title: "Expires At", Width: 18},
	{Title: "Creation Time", Width: 24},
}

func ListSystemCve(systemcve *models.CVEAllowlist) {
	var rows []table.Row
	var expiresAtStr string
	for _, cve := range systemcve.Items {
		CveName := cve.CVEID

		if systemcve.ExpiresAt != nil && *systemcve.ExpiresAt != 0 {
			expiresAt := time.Unix(int64(*systemcve.ExpiresAt), 0)
			expiresAtStr = expiresAt.Format("01/02/2006")
		} else {
			expiresAtStr = "Never expires"
		}

		createdTime, _ := utils.FormatCreatedTime(systemcve.CreationTime.String())
		rows = append(rows, table.Row{
			strconv.FormatInt(systemcve.ID, 10),
			CveName,
			expiresAtStr,
			createdTime,
		})
	}
	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
