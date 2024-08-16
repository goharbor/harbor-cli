package statisticview

import (
    "fmt"
    "github.com/charmbracelet/bubbles/table"
    "github.com/goharbor/go-client/pkg/sdk/v2.0/models"
    "github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
    tea "github.com/charmbracelet/bubbletea"
    "os"
)

var columns = []table.Column{
    {Title: "Metric", Width: 30},
    {Title: "Value", Width: 20},
}

func PrintStatistics(stats *models.Statistic) {
    rows := []table.Row{
        {"Total Projects", fmt.Sprintf("%d", stats.TotalProjectCount)},
        {"Public Projects", fmt.Sprintf("%d", stats.PublicProjectCount)},
        {"Private Projects", fmt.Sprintf("%d", stats.PrivateProjectCount)},
        {"Total Repositories", fmt.Sprintf("%d", stats.TotalRepoCount)},
        {"Total Storage Consumption", fmt.Sprintf("%d", stats.TotalStorageConsumption)},
    }

    m := tablelist.NewModel(columns, rows, len(rows))

    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}
