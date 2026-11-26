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

// Permission mode values returned by the mode selects below.
const (
	ModeNone       = "none"
	ModeKeep       = "keep"
	ModeClear      = "clear"
	ModeList       = "list"
	ModePerProject = "per_project"

	ModifyAdd     = "add"
	ModifyModify  = "modify"
	ModifyReplace = "replace"
)

// SelectProjectsForPermissions asks which projects should share one set of
// permissions.
func SelectProjectsForPermissions(projectNames []string) ([]string, error) {
	var selected []string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Multiple Project Selection").
				Description("Select the projects to assign the same permissions to this robot account."),
			huh.NewMultiSelect[string]().
				Title("Select projects").
				Options(stringOptions(projectNames)...).
				Value(&selected),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(80).Run()
	return selected, err
}

// SelectProjectsToModify asks which of the robot's projects to re-permission.
func SelectProjectsToModify(projectNames []string) ([]string, error) {
	var selected []string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select projects to modify").
				Options(stringOptions(projectNames)...).
				Value(&selected),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(80).Run()
	return selected, err
}

// ConfirmMoreProjects asks whether to add permissions for another project.
func ConfirmMoreProjects() (bool, error) {
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

// SelectCreatePermissionMode asks how project permissions should be assigned
// while creating a robot.
func SelectCreatePermissionMode() (string, error) {
	var mode string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Permission Mode").
				Description("Select how you want to assign permissions to projects:"),
			huh.NewSelect[string]().
				Title("Permission Mode").
				Description("Choose 'List' to select multiple projects with common permissions, or 'Per Project' for individual project permissions.").
				Options(
					huh.NewOption("No project permissions (system-level only)", ModeNone),
					huh.NewOption("Per Project", ModePerProject),
					huh.NewOption("List", ModeList),
				).
				Value(&mode),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).WithHeight(10).Run()
	return mode, err
}

// SelectUpdatePermissionMode asks how project permissions should be handled
// while updating a robot. The offered choices depend on whether the robot
// already has project permissions.
func SelectUpdatePermissionMode(hasExisting bool) (string, error) {
	options := []huh.Option[string]{
		huh.NewOption("No project permissions (system-level only)", ModeClear),
		huh.NewOption("Per Project (individual permissions)", ModePerProject),
		huh.NewOption("List (same permissions for multiple projects)", ModeList),
	}
	if hasExisting {
		options = []huh.Option[string]{
			huh.NewOption("Keep existing project permissions", ModeKeep),
			huh.NewOption("Clear all project permissions", ModeClear),
			huh.NewOption("Per Project (individual permissions)", ModePerProject),
			huh.NewOption("List (same permissions for multiple projects)", ModeList),
		}
	}

	var mode string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Project Permission Mode").
				Description("Select how you want to handle project permissions:"),
			huh.NewSelect[string]().
				Title("Permission Mode").
				Options(options...).
				Value(&mode),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).WithHeight(10).Run()
	return mode, err
}

// ConfirmUpdatePermissions asks whether to touch the robot's permissions at all.
func ConfirmUpdatePermissions() (bool, error) {
	return confirmNo("Do you want to update permissions?")
}

// ConfirmUpdateSystemPermissions asks whether to change system-level permissions.
func ConfirmUpdateSystemPermissions() (bool, error) {
	return confirmNo("Do you want to update system permissions?")
}

// ConfirmReplaceExistingProjects asks whether new project permissions replace
// the existing ones or are added next to them.
func ConfirmReplaceExistingProjects() (bool, error) {
	var replace bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("What do you want to do with existing project permissions?").
				Options(
					huh.NewOption("Keep existing and add new", false),
					huh.NewOption("Replace all existing with new selection", true),
				).
				Value(&replace),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).Run()
	return replace, err
}

// SelectModifyMode asks how existing project permissions should be modified.
func SelectModifyMode() (string, error) {
	var mode string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How do you want to modify project permissions?").
				Options(
					huh.NewOption("Add new projects only", ModifyAdd),
					huh.NewOption("Modify existing projects", ModifyModify),
					huh.NewOption("Replace all existing with new projects", ModifyReplace),
				).
				Value(&mode),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).Run()
	return mode, err
}

// confirmNo shows a yes/no select defaulting to "No".
func confirmNo(title string) (bool, error) {
	var answer bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title(title).
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(&answer),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).Run()
	return answer, err
}

func stringOptions(values []string) []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(values))
	for _, v := range values {
		options = append(options, huh.NewOption(v, v))
	}
	return options
}
