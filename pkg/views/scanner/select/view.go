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
package scanner

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func ScannerList(scanners []*models.ScannerRegistration) (string, error) {
	itemsList := make([]list.Item, len(scanners))

	items := map[string]string{}

	for i, s := range scanners {
		items[s.Name] = s.UUID
		itemsList[i] = selection.Item(s.Name)
	}

	m := selection.NewModel(itemsList, "Scanner")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		return "", fmt.Errorf("error running selection program: %w", err)
	}

	choice, err := selection.Selected(p)
	if err != nil {
		return "", err
	}
	return items[choice], nil
}
