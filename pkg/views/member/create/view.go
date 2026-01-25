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
	"log"

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

// map role names to role ids
var RoleOptions = map[string]int{
	"Admin":        1,
	"Developer":    2,
	"Guest":        3,
	"Maintainer":   4,
	"LimitedGuest": 5,
}

func CreateMemberView(createView *CreateView) {
	roleOptions := []string{"Project Admin", "Developer", "Guest", "Maintainer", "Limited Guest"}
	var roleSelectOptions []huh.Option[int]
	for id, name := range roleOptions {
		roleSelectOptions = append(roleSelectOptions, huh.NewOption(name, id))
	}

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
						return errors.New("Username cannot be empty.")
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
