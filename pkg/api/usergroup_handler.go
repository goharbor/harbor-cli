package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/usergroup"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
"fmt"
	log "github.com/sirupsen/logrus"
	"encoding/json"
)


func CreateUserGroup(groupName string, groupType int64, ldapGroupDn string) error {
    ctx, client, err := utils.ContextWithClient()
    if err != nil {
        return fmt.Errorf("failed to create client: %v", err)
    }

    userGroup := &models.UserGroup{
        GroupName: groupName,
        GroupType: groupType,
    }

    if groupType == 1 {
        userGroup.LdapGroupDn = ldapGroupDn
    }

    _, err = client.Usergroup.CreateUserGroup(ctx, &usergroup.CreateUserGroupParams{
        Usergroup: userGroup,
    })
	
    if err != nil {
		switch e := err.(type) {
		case *usergroup.CreateUserGroupBadRequest:
			payload, _ := json.MarshalIndent(e.Payload, "", "  ")
			return fmt.Errorf("bad request: %s", string(payload))
		case *usergroup.CreateUserGroupConflict:
			payload, _ := json.MarshalIndent(e.Payload, "", "  ")
			return fmt.Errorf("conflict: %s", string(payload))
		default:
			return fmt.Errorf("failed to create user group: %v", err)
		}
		}

    return nil

}

func DeleteUserGroup(groupId int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Usergroup.DeleteUserGroup(ctx, &usergroup.DeleteUserGroupParams{GroupID: groupId})
	if err != nil {
		return err
	}
	log.Infof("User group deleted successfully with id %d", groupId)
	return nil
}

func GetUserGroup(groupId int64) (*usergroup.GetUserGroupOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Usergroup.GetUserGroup(ctx, &usergroup.GetUserGroupParams{GroupID: groupId})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ListUserGroups() (*usergroup.ListUserGroupsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Usergroup.ListUserGroups(ctx, &usergroup.ListUserGroupsParams{})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func SearchUserGroups(groupName string) (*usergroup.SearchUserGroupsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Usergroup.SearchUserGroups(ctx, &usergroup.SearchUserGroupsParams{Groupname: groupName})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func UpdateUserGroup(groupId int64, groupName string, groupType int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Usergroup.UpdateUserGroup(ctx, &usergroup.UpdateUserGroupParams{
		GroupID: groupId,
		Usergroup: &models.UserGroup{
			GroupName: groupName,
			GroupType: groupType,
		},
	})

	if err != nil {
		return err
	}
	return nil
}