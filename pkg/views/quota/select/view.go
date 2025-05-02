// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package quota

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

// Function to get project ref details
func getRefProjectName(ref models.QuotaRefObject) (string, error) {
	if refMap, ok := ref.(map[string]interface{}); ok {
		projectName, _ := refMap["name"].(string)
		return projectName, nil
	}
	return "", fmt.Errorf("Error: Ref is not of expected type")
}

func QuotaList(quotas []*models.Quota, choice chan<- int64) {
	items := make([]list.Item, len(quotas))
	entityMap := make(map[string]int64, len(quotas))

	for i, p := range quotas {
		projectName, _ := getRefProjectName(p.Ref)
		items[i] = selection.Item(fmt.Sprintf("%v", projectName))
		entityMap[projectName] = p.ID
	}

	m := selection.NewModel(items, "Quota")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		chosen := entityMap[p.Choice]
		choice <- chosen
	}
}
