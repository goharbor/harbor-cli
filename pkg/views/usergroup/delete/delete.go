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
