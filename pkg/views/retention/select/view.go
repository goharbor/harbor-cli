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
package retention

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

var ErrUserAborted = errors.New("user aborted selection")

func RetentionList(rules []*models.RetentionRule) (int, error) {
	itemsList := make([]list.Item, len(rules))
	items := map[string]int{}

	for i, rule := range rules {
		scopeStrs := []string{}
		tagStrs := []string{}

		for k, v := range rule.ScopeSelectors {
			for _, scope := range v {
				scopeStrs = append(scopeStrs, fmt.Sprintf("%s: [%s %s]", k, scope.Decoration, scope.Pattern))
			}
		}

		for _, tag := range rule.TagSelectors {
			tagStrs = append(tagStrs, fmt.Sprintf("%s {%v}: %s", tag.Kind, tag.Extras, tag.Pattern))
		}

		// Compose detailed display string
		display := fmt.Sprintf(
			"ID: %d | Action: %s | Disabled: %v | Params: %s | Priority: %d | Scope: %s | Tags: %s | Template: %s",
			rule.ID,
			rule.Action,
			rule.Disabled,
			formatParams(rule.Params),
			rule.Priority,
			strings.Join(scopeStrs, ", "),
			strings.Join(tagStrs, ", "),
			rule.Template,
		)

		items[display] = i
		itemsList[i] = selection.Item(display)
	}

	m := selection.NewModel(itemsList, "Select a Retention Rule")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return 0, fmt.Errorf("error running selection UI: %w", err)
	}

	if model, ok := p.(selection.Model); ok {
		if model.Aborted {
			return 0, ErrUserAborted
		}
		if model.Choice == "" {
			return 0, errors.New("no retention rule selected")
		}
		return items[model.Choice], nil
	}

	return 0, errors.New("unexpected program result")
}

func formatParams(params map[string]interface{}) string {
	parts := []string{}
	for k, v := range params {
		parts = append(parts, fmt.Sprintf("%s: %v", k, v))
	}
	return strings.Join(parts, ", ")
}
