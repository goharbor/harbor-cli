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
package user

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func UserList(users []*models.UserResp) (int64, error) {
	var itemList []list.Item
	items := map[string]int64{}

	if len(users) == 0 {
		msg := views.RedText("Operation failed:")
		out := views.
			BaseStyle.
			BorderForeground(lipgloss.Color("1")).
			Render(msg, "No users found in the registry.")
		fmt.Println(out)
		return 0, fmt.Errorf("No users in the registry")
	}

	for _, r := range users {
		items[r.Username] = r.UserID
		itemList = append(itemList, selection.Item(r.Username))
	}

	m := selection.NewModel(itemList, "User")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		return 0, err
	}
	if p, ok := p.(selection.Model); ok {
		return items[p.Choice], nil
	}
	return 0, fmt.Errorf("failed to get user selection")
}
