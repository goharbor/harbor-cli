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
	"io"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: tablelist.WidthXS},
	{Title: "Name", Width: tablelist.WidthL},
	{Title: "Administrator", Width: tablelist.WidthL},
	{Title: "Email", Width: tablelist.WidthXXL},
	{Title: "Registration Time", Width: tablelist.WidthL},
}

func MakeUserRows(users []*models.UserResp) []table.Row {
	var rows []table.Row
	for _, user := range users {
		isAdmin := "No"
		if user.SysadminFlag {
			isAdmin = "Yes"
		}
		createdTime, _ := utils.FormatCreatedTime(user.CreationTime.String())
		rows = append(rows, table.Row{
			strconv.FormatInt(int64(user.UserID), 10), // UserID
			user.Username,
			isAdmin,
			user.Email,
			createdTime,
		})
	}
	return rows
}
func ListUsers(w io.Writer, users []*models.UserResp) error {
	rows := MakeUserRows(users)
	m := tablelist.NewModel(columns, rows, len(rows))
	opts := []tea.ProgramOption{tea.WithOutput(w)}
	if w != os.Stdout {
		opts = append(opts, tea.WithInput(nil))
	}
	p := tea.NewProgram(m, opts...)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to render user list: %w", err)
	}
	return nil
}
