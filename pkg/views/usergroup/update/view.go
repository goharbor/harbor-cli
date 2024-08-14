package update

import (
	"fmt"
	"strconv"

	sdkUserGroup "github.com/goharbor/go-client/pkg/sdk/v2.0/client/usergroup"
 
	"github.com/charmbracelet/huh"
 
)

type UpdateUserGroupInput struct {
	GroupID   int64
	GroupName string
	GroupType int64
}

func UpdateUserGroupView(userGroups *sdkUserGroup.ListUserGroupsOK) (*UpdateUserGroupInput, error) {
	var options []huh.Option[string]
	for _, ug := range userGroups.Payload {
		option := huh.NewOption(fmt.Sprintf("%d - %s", ug.ID, ug.GroupName), strconv.FormatInt(ug.ID, 10))
		options = append(options, option)
	}

	var selectedGroupID string
	var groupName string
	var groupType int64

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a user group to update").
				Options(options...).
				Value(&selectedGroupID),
			huh.NewInput().
				Title("Enter new group name").
				Value(&groupName),
			huh.NewSelect[int64]().
				Title("Select group type").
				Options(
					huh.NewOption("LDAP", int64(1)),
					huh.NewOption("HTTP", int64(2)),
					huh.NewOption("OIDC", int64(3)),
				).
				Value(&groupType),
		),
	)

	err := form.Run()
	if err != nil {
		return nil, fmt.Errorf("form input error: %v", err)
	}
	groupID, err := strconv.ParseInt(selectedGroupID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %v", err)
	}

	return &UpdateUserGroupInput{
		GroupID:   groupID,
		GroupName: groupName,
		GroupType: groupType,
	}, nil
}