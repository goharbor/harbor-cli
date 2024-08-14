package usergroup

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/usergroup/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UserGroupUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update user group",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			userGroupsResp, err := api.ListUserGroups()
			if err != nil {
				log.Errorf("failed to list user groups: %v", err)
				return
			}
			input, err := update.UpdateUserGroupView(userGroupsResp)
			if err != nil {
				log.Errorf("failed to get user input: %v", err)
				return
			}
			err = api.UpdateUserGroup(input.GroupID, input.GroupName, input.GroupType)
			if err != nil {
				log.Errorf("failed to update user group: %v", err)
			} else {
				log.Infof("User group `%s` updated successfully", input.GroupName)
			}
		},
	}

	return cmd
}