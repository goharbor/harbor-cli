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
	"github.com/charmbracelet/huh"
)

type PermissionSelection struct {
	Resource string
	Action   string
}

// ChooseProjectPermissionMode returns the selected mode for project perms.
// If hasExisting is true, it offers keep/clear/list/per_project; otherwise clear/list/per_project.
func ChooseProjectPermissionMode(hasExisting bool) (string, error) {
	var permissionMode string
	var options []huh.Option[string]

	if hasExisting {
		options = []huh.Option[string]{
			huh.NewOption("Keep existing project permissions", "keep"),
			huh.NewOption("Clear all project permissions", "clear"),
			huh.NewOption("Per Project (individual permissions)", "per_project"),
			huh.NewOption("List (same permissions for multiple projects)", "list"),
		}
	} else {
		options = []huh.Option[string]{
			huh.NewOption("No project permissions (system-level only)", "clear"),
			huh.NewOption("Per Project (individual permissions)", "per_project"),
			huh.NewOption("List (same permissions for multiple projects)", "list"),
		}
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Project Permission Mode").
				Description("Select how you want to handle project permissions:"),
			huh.NewSelect[string]().
				Title("Permission Mode").
				Options(options...).
				Value(&permissionMode),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).WithHeight(10).Run()

	return permissionMode, err
}

// AskMoreProjects asks if the user wants to keep adding projects.
func AskMoreProjects() (bool, error) {
	var addMore bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Project Selection").
				Description("You can add permissions for multiple projects to this robot account."),
			huh.NewSelect[bool]().
				Title("Do you want to select (more) projects?").
				Description("Select 'Yes' to add (another) project, 'No' to continue with current selection.").
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(&addMore),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).WithHeight(10).Run()

	return addMore, err
}

// ConfirmReplaceExisting asks whether to replace existing project permissions.
func ConfirmReplaceExisting() (bool, error) {
	var replaceExisting bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("What do you want to do with existing project permissions?").
				Options(
					huh.NewOption("Keep existing and add new", false),
					huh.NewOption("Replace all existing with new selection", true),
				).
				Value(&replaceExisting),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).Run()
	return replaceExisting, err
}

// ChooseModifyMode asks how to modify project permissions when some exist.
// Returns one of: add | modify | replace
func ChooseModifyMode() (string, error) {
	var modifyMode string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How do you want to modify project permissions?").
				Options(
					huh.NewOption("Add new projects only", "add"),
					huh.NewOption("Modify existing projects", "modify"),
					huh.NewOption("Replace all existing with new projects", "replace"),
				).
				Value(&modifyMode),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).Run()
	return modifyMode, err
}

// SelectProjects shows a multi-select of project names.
// The actual project list is provided externally by the existing prompt package.
// Here we invoke the existing provider to avoid duplicating project retrieval.
func SelectProjects(getProjects func() ([]string, error)) ([]string, error) {
	projects, err := getProjects()
	if err != nil {
		return nil, err
	}
	var selected []string
	var opts []huh.Option[string]
	for _, p := range projects {
		opts = append(opts, huh.NewOption(p, p))
	}
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select projects").
				Options(opts...).
				Value(&selected),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(80).Run()
	return selected, err
}

// AskUpdateSystemPerms asks in update flow if system permissions should be updated.
func AskUpdateSystemPerms() (bool, error) {
	var updateSystem bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Do you want to update system permissions?").
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(&updateSystem),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).Run()
	return updateSystem, err
}
