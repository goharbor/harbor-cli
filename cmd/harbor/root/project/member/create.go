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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/sirupsen/logrus"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"

	"github.com/goharbor/harbor-cli/pkg/views/member/create"
	"github.com/spf13/cobra"
)

func CreateMemberCommand() *cobra.Command {
	var opts create.CreateView
	opts.MemberUser = &models.UserEntity{} // Initialize MemberUser
	opts.MemberGroup = &models.UserGroup{} // Initialize MemberGroup
	var isID bool

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create project member",
		Long:    "create project member by Name",
		Example: "  harbor project member create my-project --username user --role Developer",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) > 0 {
				_, checkErr := api.GetProject(args[0], isID)
				if checkErr != nil {
					if utils.ParseHarborErrorCode(checkErr) == "404" {
						return fmt.Errorf("project %s does not exist", args[0])
					}
					return fmt.Errorf("failed to verify project: %v", utils.ParseHarborErrorMsg(checkErr))
				}
				opts.ProjectName = args[0]
			} else {
				opts.ProjectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
			}

			sysInfo, err := api.GetSystemInfo()
			if err != nil {
				fmt.Println("could not access server info")
			}

			createView := &create.CreateView{
				AuthMode:      *sysInfo.Payload.AuthMode,
				XIsResourceID: !isID,
				ProjectName:   opts.ProjectName,
				RoleID:        opts.RoleID,
				MemberUser: &models.UserEntity{
					UserID:   opts.MemberUser.UserID,
					Username: opts.MemberUser.Username,
				},
				MemberGroup: &models.UserGroup{
					ID:          opts.MemberGroup.ID,
					GroupName:   opts.MemberGroup.GroupName,
					GroupType:   opts.MemberGroup.GroupType,
					LdapGroupDn: opts.MemberGroup.LdapGroupDn,
				},
			}

			// check if role and member is valid
			if opts.RoleID != 0 && opts.MemberUser.Username != "" {
				err = api.CreateMember(*createView)
			} else {
				err = createMemberView(createView)
			}

			if err != nil {
				logrus.Errorf("failed to create user: %v", err)
			}

			fmt.Printf("successfully added user %s to project %s\n", createView.MemberUser.Username, opts.ProjectName)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "parses projectName as an ID")
	flags.StringVarP(&opts.RoleName, "role", "", "", "Role Name [one of Project_Admin, Developer, Guest, Maintainer, Limited_Guest]")
	flags.IntVarP(&opts.RoleID, "roleid", "", 0, "Role ID")
	flags.StringVarP(&opts.MemberUser.Username, "username", "", "", "Username")
	flags.StringVarP(&opts.MemberGroup.GroupName, "groupname", "", "", "Group Name")
	flags.StringVarP(&opts.MemberGroup.LdapGroupDn, "ldapdn", "", "", "DN of LDAP Group")
	flags.Int64VarP(&opts.MemberGroup.ID, "groupid", "", 0, "Group ID")
	flags.Int64VarP(&opts.MemberUser.UserID, "userid", "", 0, "User ID")
	flags.Int64VarP(&opts.MemberGroup.GroupType, "grouptype", "", 0, "Group Type")

	return cmd
}

func createMemberView(createView *create.CreateView) error {
	create.CreateMemberView(createView)
	return api.CreateMember(*createView)
}
