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

package create

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

type MemberUser struct {
	UserID   int
	Username string
}

type MemberGroup struct {
	ID          int
	GroupName   string
	GroupType   int
	LdapGroupDN string
}

type CreateView struct {
	AuthMode      string
	XIsResourceID bool
	ProjectName   string
	RoleID        int
	RoleName      string
	MemberUser    *models.UserEntity
	MemberGroup   *models.UserGroup
}

type Role struct {
	ID      int
	Name    string
	Aliases []string
}

// Roles is the single source of truth for project member roles in Harbor.
var Roles = []Role{
	{ID: 1, Name: "Project Admin", Aliases: []string{"admin", "projectadmin", "project_admin"}},
	{ID: 2, Name: "Developer", Aliases: []string{"developer"}},
	{ID: 3, Name: "Guest", Aliases: []string{"guest"}},
	{ID: 4, Name: "Maintainer", Aliases: []string{"maintainer"}},
	{ID: 5, Name: "Limited Guest", Aliases: []string{"limitedguest", "limited_guest"}},
}

// RoleOptions maps role names to role ids for backward compatibility.
var RoleOptions = map[string]int{
	"Admin":         1,
	"Project Admin": 1,
	"Developer":     2,
	"Guest":         3,
	"Maintainer":    4,
	"LimitedGuest":  5,
	"Limited Guest": 5,
}

// RoleSelectOptions builds interactive options for huh select components.
func RoleSelectOptions() []huh.Option[int] {
	roleSelectOptions := make([]huh.Option[int], 0, len(Roles))
	for _, role := range Roles {
		roleSelectOptions = append(roleSelectOptions, huh.NewOption(role.Name, role.ID))
	}
	return roleSelectOptions
}

// ResolveRoleName maps a role name or alias to its 1-based Role ID.
func ResolveRoleName(roleName string) (int, error) {
	normalized := strings.ToLower(strings.TrimSpace(roleName))
	normalized = strings.ReplaceAll(normalized, "_", "")
	normalized = strings.ReplaceAll(normalized, " ", "")

	for _, role := range Roles {
		for _, alias := range role.Aliases {
			aliasNorm := strings.ReplaceAll(alias, "_", "")
			if normalized == aliasNorm {
				return role.ID, nil
			}
		}
	}

	return 0, fmt.Errorf("invalid role name %q: allowed values are Project_Admin, Developer, Guest, Maintainer, Limited_Guest", roleName)
}

// IsValidRoleID reports whether the given role ID exists in Roles.
func IsValidRoleID(roleID int) bool {
	for _, role := range Roles {
		if role.ID == roleID {
			return true
		}
	}
	return false
}

// AllowedRolesString returns a formatted string of all allowed role IDs and names.
func AllowedRolesString() string {
	parts := make([]string, 0, len(Roles))
	for _, role := range Roles {
		parts = append(parts, fmt.Sprintf("%d (%s)", role.ID, role.Name))
	}
	return strings.Join(parts, ", ")
}

func CreateMemberView(createView *CreateView) {
	roleSelectOptions := RoleSelectOptions()

	groupOptions := []string{"None", "LDAP group", "HTTP group", "OIDC group"}
	var groupSelectOptions []huh.Option[int]
	for id, name := range groupOptions {
		groupSelectOptions = append(groupSelectOptions, huh.NewOption(name, id))
	}

	groups := []*huh.Group{}

	// Show role select only if not already given
	if createView.RoleID == 0 && createView.RoleName == "" {
		groups = append(groups,
			huh.NewGroup(
				huh.NewSelect[int]().
					Description("Select a Role").
					Title("Role").
					Value(&createView.RoleID).
					Options(roleSelectOptions...),
			))
	}

	if createView.MemberUser.UserID == 0 && createView.MemberUser.Username == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Value(&createView.MemberUser.Username).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("username cannot be empty.")
					}

					if !utils.ValidateUserName(str) {
						return errors.New("Invalid username. Must be 1-255 characters long and cannot contain special characters: , \" ~ # % $")
					}

					return nil
				}),
		))
	}

	// show group info only if LDAP
	if createView.AuthMode == "ldap" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Group Name (optional)").
				Value(&createView.MemberGroup.GroupName).
				Validate(func(str string) error {
					return nil
				}),
		))
	}

	// if groupname is populated, show groupType
	if createView.MemberGroup.GroupName != "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Group Name (optional)").
				Value(&createView.MemberGroup.GroupName).
				Validate(func(str string) error {
					return nil
				}),
		), huh.NewGroup(
			huh.NewInput().
				Title("DN of LDAP group (optional)").
				Value(&createView.MemberGroup.GroupName).
				Validate(func(str string) error {
					return nil
				}),
		))
	}

	theme := huh.ThemeCharm()
	err := huh.NewForm(groups...).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}
