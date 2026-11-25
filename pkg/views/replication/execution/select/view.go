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
package execution

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

var ErrUserAborted = errors.New("user aborted selection")

func ReplicationExecutionList(executions []*models.ReplicationExecution) (int64, error) {
	labelIDs := make(map[string]int64, len(executions))
	itemsList := make([]list.Item, len(executions))
	for i, p := range executions {
		label := fmt.Sprintf("ID: %d, Status: %s, Trigger: %s, Start Time: %s, Succeed: %d, Total: %d",
			p.ID, p.Status, p.Trigger, p.StartTime.String(), p.Succeed, p.Total)
		labelIDs[label] = p.ID
		itemsList[i] = selection.Item(label)
	}

	m := selection.NewModel(itemsList, "Select a Replication Execution")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return 0, fmt.Errorf("error running selection program: %w", err)
	}

	model, ok := p.(selection.Model)
	if !ok {
		return 0, errors.New("unexpected program result")
	}
	if model.Aborted {
		return 0, ErrUserAborted
	}
	id, ok := labelIDs[model.Choice]
	if !ok {
		return 0, errors.New("no replication execution selected")
	}
	return id, nil
}
