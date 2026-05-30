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
package task

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

var ErrUserAborted = errors.New("user aborted selection")

func ReplicationTasksList(tasks []*models.ReplicationTask, choice chan<- int64, errChan chan<- error) {
	itemsList := make([]list.Item, len(tasks))
	for i, p := range tasks {
		displayName := fmt.Sprintf("ID: %d, Status: %s, Operation: %s, Src: %s, Dst: %s, Start Time: %s",
			p.ID, p.Status, p.Operation, p.SrcResource, p.DstResource, p.StartTime.String())
		itemsList[i] = selection.Item(displayName)
	}

	m := selection.NewModel(itemsList, "Select a Replication Task")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		errChan <- fmt.Errorf("error running selection program: %w", err)
		return
	}

	if model, ok := p.(selection.Model); ok {
		if model.Choice == "" {
			errChan <- errors.New("user aborted selection")
			return
		}
		// Extract the ID from model.Choice
		var taskID int64
		_, err = fmt.Sscanf(model.Choice, "ID: %d", &taskID)
		if err != nil {
			errChan <- fmt.Errorf("error parsing task ID: %w", err)
			return
		}
		choice <- taskID
		return
	}

	errChan <- errors.New("unexpected program result")
}
