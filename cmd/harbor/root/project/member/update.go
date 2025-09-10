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

package member

import (
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

// NewGetRegistryCommand creates a new `harbor get registry` command
func UpdateMemberCommand() *cobra.Command {
	var opts api.UpdateMemberOptions
	var roleID int64
	var isID bool

	cmd := &cobra.Command{
		Use:     "update [ProjectName] [member ID]",
		Short:   "update member by ID or name",
		Long:    "update member in a project by MemberID",
		Example: "  harbor project member update my-project [memberID] --roleid 2",
		Args:    cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) == 1 {
				opts.ProjectNameOrID = args[0]
			} else if len(args) == 2 {
				opts.ProjectNameOrID = args[0]
				opts.ID, _ = strconv.ParseInt(args[1], 0, 64)
			}

			if opts.ProjectNameOrID == "" {
				opts.ProjectNameOrID, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
			}

			if opts.ID == 0 {
				opts.ID = prompt.GetMemberIDFromUser(opts.ProjectNameOrID)
			}

			if roleID == 0 {
				roleID = prompt.GetRoleIDFromUser()
			}
			opts.RoleID = &models.RoleRequest{
				RoleID: roleID,
			}

			// when set true parses projectNameOrID as projectName
			// else it parses as an integer ID
			opts.XIsResourceName = !isID

			err = api.UpdateMember(opts)
			if err != nil {
				return fmt.Errorf("failed to get members list: %v", err)
			}

			fmt.Printf("successfully updated user with ID %d with role ID %d for project %s\n", opts.ID, opts.RoleID, opts.ProjectNameOrID)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "parses projectName as an ID")
	flags.Int64VarP(&opts.ID, "member-id", "", 0, "Member ID")
	flags.Int64VarP(&roleID, "roleid", "", 0, "Role to be updated")
	return cmd
}
