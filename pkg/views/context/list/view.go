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

package list

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Name", Width: tablelist.Width3XL},
	{Title: "Username", Width: tablelist.WidthL},
	{Title: "Server Address", Width: tablelist.WidthXXL},
}

func ListContexts(contexts []api.ContextListView, currentCredential string) {
	rows := selectActiveContext(contexts, currentCredential)

	var opts []tea.ProgramOption
	if !term.IsTerminal(int(os.Stdout.Fd())) { // #nosec G115 - fd fits in int on all supported platforms
		opts = append(opts, tea.WithoutRenderer(), tea.WithInput(nil))
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m, opts...).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func selectActiveContext(contexts []api.ContextListView, currentCredential string) []table.Row {
	var rows []table.Row

	for _, ctx := range contexts {
		if ctx.Name == currentCredential {
			rows = append(rows, table.Row{
				"* " + ctx.Name,
				ctx.Username,
				ctx.Server,
			})
		} else {
			rows = append(rows, table.Row{
				"  " + ctx.Name,
				ctx.Username,
				ctx.Server,
			})
		}
	}
	return rows
}
