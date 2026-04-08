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
package tablelist

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
	"github.com/goharbor/harbor-cli/pkg/views"
)

const (
	WidthXS  = 4
	WidthS   = 8
	WidthM   = 12
	WidthL   = 16
	WidthXL  = 20
	WidthXXL = 24
	Width3XL = 30
)

func NewModel(columns []table.Column, rows []table.Row, height int) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(s.Cell.GetForeground()).
		Background(s.Cell.GetBackground()).
		Bold(false)
	t.SetStyles(s)

	return t
}

// Render returns the table rendered as a string with the base style applied.
func Render(t table.Model) string {
	return views.BaseStyle.Render(t.View()) + "\n"
}
