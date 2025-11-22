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
package delete

import (
	"fmt"

	sdkUserGroup "github.com/goharbor/go-client/pkg/sdk/v2.0/client/usergroup"

	"github.com/charmbracelet/huh"
)

type DeleteUserGroupInput struct {
	ID      int64
	Name    string
	Confirm bool
}

func DeleteUserGroupView(userGroups *sdkUserGroup.ListUserGroupsOK) (*DeleteUserGroupInput, error) {
	var options []huh.Option[DeleteUserGroupInput]
	for _, ug := range userGroups.Payload {
		option := huh.NewOption(fmt.Sprintf("%d - %s", ug.ID, ug.GroupName), DeleteUserGroupInput{
			ID:      ug.ID,
			Name:    ug.GroupName,
			Confirm: false,
		})
		options = append(options, option)
	}

	var selectedGroup DeleteUserGroupInput

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[DeleteUserGroupInput]().
				Title("Select a user group to update").
				Options(options...).
				Value(&selectedGroup),
		),
	)

	err := form.Run()
	if err != nil {
		return nil, fmt.Errorf("form input error: %v", err)
	}

	err = huh.NewConfirm().
		Title(fmt.Sprintf("Delete Usergroup '%s'", selectedGroup.Name)).
		Value(&selectedGroup.Confirm).
		Run()
	if err != nil {
		return nil, fmt.Errorf("form input error: %v", err)
	}

	return &selectedGroup, nil
}
