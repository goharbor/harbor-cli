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
package robot

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func ListPermissions(perms *models.Permissions, ch chan<- []models.Permission) {
	permissions := perms.Project
	choices := []models.Permission{}

	// collect all possible permissions
	for _, perm := range permissions {
		choices = append(choices, *perm)
	}
	selects := &[]models.Permission{}

	km := huh.NewDefaultKeyMap()
	km.MultiSelect.SelectAll = key.NewBinding(
		key.WithKeys("a"), key.WithHelp("a", "all"),
	)
	km.MultiSelect.SelectNone = key.NewBinding(
		key.WithKeys("A"), key.WithHelp("A", "none"),
	)

	// groups will hold the multi-select groups for each resource
	// selections will hold the selected actions for each resource
	var (
		groups     []*huh.Group
		selections = make(map[string]*[]string)
	)

	// resActs groups the actions w.r.t. their resources
	// e.g. {"repository": ["read", "write"], "project": ["read"]}
	// this is needed for easier options creation in the multi-select
	resActs := map[string][]string{}
	for _, p := range permissions {
		resActs[p.Resource] = append(resActs[p.Resource], p.Action)
	}

	resources := make([]string, 0, len(resActs))
	for r := range resActs {
		resources = append(resources, r)
	}
	sort.Strings(resources)

	// loop over all resources and create a multi select for each
	// within each resource a multi-select for the actions is created
	// e.g. for "repository" resource, a multi-select with options "read", "write" is created
	// the selected actions are stored in selections map
	// e.g. selections["repository"] = &[]string{"read", "write"}
	for _, res := range resources {
		acts := resActs[res]
		sort.Strings(acts)
		opts := make([]huh.Option[string], len(acts))
		for i, a := range acts {
			opts[i] = huh.NewOption(a, a)
		}
		pick := []string{}
		selections[res] = &pick
		ms := huh.NewMultiSelect[string]().
			Options(opts...).
			Title(strings.ToUpper(res)).
			Value(&pick).
			WithKeyMap(km)
		groups = append(groups, huh.NewGroup(ms))
	}

	// all huh groups are displaye in the form grid
	form := huh.NewForm(groups...).
		WithLayout(huh.LayoutGrid(4, 6)).
		WithTheme(huh.ThemeDracula())
	if err := form.Run(); err != nil {
		fmt.Println("error:", err)
		return
	}

	// we have created a resource: [actions] like mapping earlier
	// now we have to convert it back to []models.Permission
	for res, acts := range selections {
		for _, act := range *acts {
			*selects = append(*selects, models.Permission{
				Resource: res,
				Action:   act,
			})
		}
	}

	//ToDo: This has to generalized for a multi-select view and moved out of this file

	// m := multiselect.NewModel(choices, selects)

	// _, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	// if err != nil {
	// 	fmt.Println("Error running program:", err)
	// }
	// Get selected permissions
	ch <- *selects
}

func ListRobot(robots []*models.Robot, choice chan<- int64) {
	itemsList := make([]list.Item, len(robots))

	items := map[string]int64{}

	for i, r := range robots {
		items[r.Name] = r.ID
		itemsList[i] = selection.Item(r.Name)
	}

	m := selection.NewModel(itemsList, "Robot")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- items[p.Choice]
	}
}
