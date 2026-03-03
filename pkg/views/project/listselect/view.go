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
package project

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/listselect"
)

var ErrUserAborted = errors.New("user aborted selection")

func ProjectsList(projects []*models.Project) ([]string, error) {
	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = listselect.Item(p.Name)
	}

	m := listselect.NewModel(items, "Project")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		return nil, fmt.Errorf("error running selection program: %w", err)
	}

	if model, ok := p.(listselect.Model); ok {
		if model.Aborted {
			return nil, ErrUserAborted
		}
		if len(model.Choices) == 0 {
			return nil, errors.New("no project selected")
		}
		return model.Choices, nil
	}

	return nil, errors.New("unexpected program result")
}

func ProjectsListWithId(projects []*models.Project) ([]int64, error) {
	items := make([]list.Item, len(projects))
	itemsMap := make(map[string]int64)

	for i, p := range projects {
		items[i] = listselect.Item(p.Name)
		itemsMap[p.Name] = int64(p.ProjectID)
	}

	m := listselect.NewModel(items, "Project")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return nil, fmt.Errorf("error running selection program: %w", err)
	}

	if model, ok := p.(listselect.Model); ok {
		if model.Aborted {
			return nil, ErrUserAborted
		}
		if len(model.Choices) == 0 {
			return nil, errors.New("no project selected")
		}
		var selectedIDs []int64
		for _, choice := range model.Choices {
			selectedIDs = append(selectedIDs, itemsMap[choice])
		}
		return selectedIDs, nil
	}

	return nil, errors.New("unexpected program result")
}
