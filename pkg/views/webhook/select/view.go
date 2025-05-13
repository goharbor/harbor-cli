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

func WebhookList(webhooks []*models.WebhookPolicy) (models.WebhookPolicy, error) {
	items := make([]list.Item, len(webhooks))
	for i, item := range webhooks {
		items[i] = selection.Item(item.Name)
	}

	m := selection.NewModel(items, "Webhook")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return models.WebhookPolicy{}, fmt.Errorf("error running selection program: %w", err)
	}

	if model, ok := p.(selection.Model); ok {
		if model.Aborted {
			return models.WebhookPolicy{}, errors.New("user aborted selection")
		}
		if model.Choice == "" {
			return models.WebhookPolicy{}, errors.New("no webhook selected")
		}
		for _, webhook := range webhooks {
			if webhook.Name == model.Choice {
				return *webhook, nil
			}
		}
	}

	return models.WebhookPolicy{}, errors.New("unexpected program result")
}
