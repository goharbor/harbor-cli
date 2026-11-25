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
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

// stubProgram replaces the interactive UI with one that finishes as result.
func stubProgram(t *testing.T, result func(m selection.Model) selection.Model) {
	t.Helper()
	orig := runProgram
	runProgram = func(m tea.Model) (tea.Model, error) {
		return result(m.(selection.Model)), nil
	}
	t.Cleanup(func() { runProgram = orig })
}

func TestRetentionList(t *testing.T) {
	rules := []*models.RetentionRule{
		{ID: 7, Action: "retain", Template: "always"},
		{ID: 9, Action: "retain", Template: "latestPushedK"},
	}

	t.Run("cancel returns ErrUserAborted", func(t *testing.T) {
		stubProgram(t, func(m selection.Model) selection.Model {
			m.Aborted = true
			return m
		})

		id, err := RetentionList(rules)
		if !errors.Is(err, ErrUserAborted) {
			t.Fatalf("RetentionList() error = %v, want %v", err, ErrUserAborted)
		}
		if id != 0 {
			t.Errorf("RetentionList() id = %d, want 0", id)
		}
	})

	t.Run("empty choice returns error", func(t *testing.T) {
		stubProgram(t, func(m selection.Model) selection.Model { return m })

		if _, err := RetentionList(rules); err == nil {
			t.Fatal("RetentionList() error = nil, want error")
		}
	})

	t.Run("choice returns selected rule ID", func(t *testing.T) {
		stubProgram(t, func(m selection.Model) selection.Model {
			m.Choice = string(m.List.Items()[1].(selection.Item))
			return m
		})

		id, err := RetentionList(rules)
		if err != nil {
			t.Fatalf("RetentionList() error = %v", err)
		}
		if id != 9 {
			t.Errorf("RetentionList() id = %d, want 9", id)
		}
	})
}
