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
