package volume

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Volume", Width: 8},
	{Title: "Total", Width: 20},
	{Title: "Free", Width: 15},
}

func PrintVolumeInfo(storage []*models.Storage) {
	var rows []table.Row

	if len(storage) == 0 {
		fmt.Println("No volume information available.")
		return
	}

	for i, vol := range storage {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			formatSize(vol.Total),
			formatSize(vol.Free),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

func formatSize(size uint64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := uint64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
