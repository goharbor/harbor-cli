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
	{Title: "ID", Width: 4},
	{Title: "Project", Width: 12},
	{Title: "Owner Name", Width: 12},
	{Title: "Storage", Width: 32},
	{Title: "Creation Time", Width: 20},
}

// Function to get project details
func GetProjectDetails(ref models.QuotaRefObject) (string, string, error) {
	if refMap, ok := ref.(map[string]interface{}); ok {
		projectName, _ := refMap["name"].(string)
		ownerName, _ := refMap["owner_name"].(string)
		return projectName, ownerName, nil
	}
	return "", "", fmt.Errorf("Error: Ref is not of expected type")
}

// Function to convert bytes to human-readable storage
func BytesToStorageString(bytes int64) string {
	const (
		mebibyte = 1024 * 1024
		gibibyte = 1024 * mebibyte
	)

	mib := float64(bytes) / float64(mebibyte)

	if mib >= 1024 {
		gib := mib / 1024
		return fmt.Sprintf("%.1f GiB", gib)
	}

	return fmt.Sprintf("%.2f MiB", mib)
}

// Function to calculate storage
func CalculateStorage(hard models.ResourceList, used models.ResourceList) (string, string) {
	var storageUsed, storageGiven string

	if hard["storage"] == -1 {
		storageGiven = "Unlimited"
	} else {
		storageGiven = BytesToStorageString(hard["storage"])
	}

	if used["storage"] == 0 {
		storageUsed = "0 MiB"
	} else {
		storageUsed = BytesToStorageString(used["storage"])
	}

	return storageUsed, storageGiven
}

// Function to format storage
func FormatStorage(hard models.ResourceList, used models.ResourceList) string {
	storageUsed, storageGiven := CalculateStorage(hard, used)
	return fmt.Sprintf("%v of %v", storageUsed, storageGiven)
}

// ListQuotas in table format
func ListQuotas(quotas []*models.Quota) {
	var rows []table.Row
	for _, quota := range quotas {
		projectName, ownerName, err := GetProjectDetails(quota.Ref)
		if err != nil {
			fmt.Println(err)
			continue
		}

		storage := FormatStorage(quota.Hard, quota.Used)

		createdTime, _ := utils.FormatCreatedTime(quota.CreationTime.String())
		rows = append(rows, table.Row{
			strconv.FormatInt(quota.ID, 10),
			projectName,
			ownerName,
			storage,
			createdTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
