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
package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	GreenStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // ANSI 32
	RedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // ANSI 31
	BoldStyle         = lipgloss.NewStyle().Bold(true)                      // ANSI 1
	YellowStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // ANSI 33
	BlueStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("4")) // ANSI 34
	GrayStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // ANSI 37
)

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).Padding(0, 1)

func RedText(strs ...string) string {
	var msg strings.Builder
	for _, str := range strs {
		msg.WriteString(str)
	}
	return RedStyle.Render(msg.String())
}
