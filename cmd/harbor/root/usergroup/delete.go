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
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	delete "github.com/goharbor/harbor-cli/pkg/views/usergroup/delete"
	"github.com/spf13/cobra"
)

func UserGroupDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [groupID]",
		Short: "delete user group",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var groupID int64
			var err error

			if len(args) > 0 && args[0] != "" {
				groupID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert argument to ID: %v", err)
				}
			} else {
				response, err := api.ListUserGroups()
				if err != nil {
					return fmt.Errorf("failed to list user groups: %v", err)
				}

				opts, err := delete.DeleteUserGroupView(response)
				if err != nil {
					return fmt.Errorf("failed to obtain user selection: %v", err)
				}

				groupID = opts.ID
			}

			err = api.DeleteUserGroup(groupID)
			if err != nil {
				return fmt.Errorf("failed to delete user group: %v", err)
			}

			return nil
		},
	}

	return cmd
}
