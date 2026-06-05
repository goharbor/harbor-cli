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
package instance

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

var ErrUserAborted = errors.New("user aborted selection")

func InstanceList(instances []*models.Instance) (string, error) {
	items := make([]list.Item, len(instances))
	for i, instance := range instances {
		items[i] = selection.Item(instance.Name)
	}

	m := selection.NewModel(items, "Instance")

	p, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if model, ok := p.(selection.Model); ok {
		if model.Aborted {
			return "", ErrUserAborted
		}
		if model.Choice == "" {
			return "", errors.New("no instance selected")
		}
		return model.Choice, nil
	}

	return "", errors.New("unexpected program result")
}
