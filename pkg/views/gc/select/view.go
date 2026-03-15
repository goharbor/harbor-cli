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
package gcselect

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

var ErrUserAborted = errors.New("user aborted selection")

func GCJobList(history []*models.GCHistory) (int64, error) {
	items := make([]list.Item, len(history))
	jobsMap := make(map[string]int64)

	for i, job := range history {
		creationTime, _ := utils.FormatCreatedTime(job.CreationTime.String())
		displayName := fmt.Sprintf("ID: %d | Status: %s | Created: %s",
			job.ID, job.JobStatus, creationTime)
		items[i] = selection.Item(displayName)
		jobsMap[displayName] = job.ID
	}

	m := selection.NewModel(items, "GC Job")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return 0, fmt.Errorf("error running selection program: %w", err)
	}

	if model, ok := p.(selection.Model); ok {
		if model.Aborted {
			return 0, ErrUserAborted
		}
		if model.Choice == "" {
			return 0, errors.New("no GC job selected")
		}
		return jobsMap[model.Choice], nil
	}

	return 0, errors.New("unexpected program result")
}
