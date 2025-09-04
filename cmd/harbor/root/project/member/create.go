package member

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/sirupsen/logrus"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"

	"github.com/goharbor/harbor-cli/pkg/views/member/create"
	"github.com/spf13/cobra"
)

func CreateMemberCommand() *cobra.Command {
	var opts create.CreateView
	opts.MemberUser = &models.UserEntity{} // Initialize MemberUser
	opts.MemberGroup = &models.UserGroup{} // Initialize MemberGroup

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create project member",
		Long:    "create project member by Name",
		Example: "  harbor project member create my-project --username user --role Developer",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) > 0 {
				ok, checkErr := api.CheckProject(args[0]) // verifying project name
				if checkErr != nil {
					return fmt.Errorf("failed to verify project name: %v", checkErr)
				}

				if ok {
					opts.ProjectName = args[0]
				} else {
					return fmt.Errorf("invalid project name: %s", args[0])
				}
			} else {
				opts.ProjectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
			}

			createView := &create.CreateView{
				ProjectName: opts.ProjectName,
				RoleID:      opts.RoleID,
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

			return nil
		},
	}

	flags := cmd.Flags()
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

func setRoleIDFromRoleName(opts *create.CreateView) {
	if opts.RoleName != "" {
		if id, ok := create.RoleOptions[opts.RoleName]; ok {
			opts.RoleID = id
		}
	}
}
