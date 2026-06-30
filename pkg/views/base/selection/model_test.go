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
package selection

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelUpdateMarksSelectionAborted(t *testing.T) {
	for _, key := range []string{"q", "esc", "ctrl+c"} {
		t.Run(key, func(t *testing.T) {
			model := NewModel([]list.Item{Item("project")}, "Project")

			updated, cmd := model.Update(keyMsg(key))
			selectionModel, ok := updated.(Model)
			require.True(t, ok)

			assert.True(t, selectionModel.Aborted)
			assert.Empty(t, selectionModel.Choice)
			assert.NotNil(t, cmd)
		})
	}
}

func TestModelUpdateSelectsChoice(t *testing.T) {
	model := NewModel([]list.Item{Item("project")}, "Project")

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectionModel, ok := updated.(Model)
	require.True(t, ok)

	assert.False(t, selectionModel.Aborted)
	assert.Equal(t, "project", selectionModel.Choice)
	assert.NotNil(t, cmd)
}

func TestSelectedChoiceReturnsCancelError(t *testing.T) {
	model := Model{Aborted: true}

	choice, err := model.SelectedChoice()

	assert.Empty(t, choice)
	assert.ErrorIs(t, err, ErrUserAborted)
}

func keyMsg(key string) tea.KeyMsg {
	switch key {
	case "q":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
}
