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
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

var ErrUserAborted = errors.New("user aborted selection")

func ProjectList(projects []*models.Project) (string, error) {
	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = selection.Item(p.Name)
	}

	m := selection.NewModel(items, "Project")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return "", fmt.Errorf("error running selection program: %w", err)
	}

	if model, ok := p.(selection.Model); ok {
		if model.Aborted {
			return "", ErrUserAborted
		}
		if model.Choice == "" {
			return "", errors.New("no project selected")
		}
		return model.Choice, nil
	}

	return "", errors.New("unexpected program result")
}

func ProjectListWithId(projects []*models.Project) (int64, error) {
	items := make([]list.Item, len(projects))
	itemsMap := make(map[string]int64)

	for i, p := range projects {
		items[i] = selection.Item(p.Name)
		itemsMap[p.Name] = int64(p.ProjectID)
	}

	m := selection.NewModel(items, "Project")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return 0, fmt.Errorf("error running selection program: %w", err)
	}

	if model, ok := p.(selection.Model); ok {
		if model.Aborted {
			return 0, ErrUserAborted
		}
		if model.Choice == "" {
			return 0, errors.New("no project selected")
		}
		return itemsMap[model.Choice], nil
	}

	return 0, errors.New("unexpected program result")
}
