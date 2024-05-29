package member

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"

	"github.com/goharbor/harbor-cli/pkg/views/member/create"
	"github.com/spf13/cobra"
)

func CreateMemberCommand() *cobra.Command {
	var opts create.CreateView
	opts.MemberUser = &models.UserEntity{} // Initialize MemberUser
	opts.MemberGroup = &models.UserGroup{} // Initialize MemberGroup

	cmd := &cobra.Command{
		Use:     "create [ProjectName Or ID]",
		Short:   "create project member",
		Long:    "create project member by Name",
		Example: "  harbor member create my-project --username user --roleid 1",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				opts.ProjectNameOrID = args[0]
			} else {
				opts.ProjectNameOrID = prompt.GetProjectNameFromUser()
			}

			createView := &create.CreateView{
				ProjectNameOrID: opts.ProjectNameOrID,
				RoleID:          opts.RoleID,
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

			if opts.RoleID != 0 && opts.MemberUser.Username != "" {
				err = api.CreateMember(*createView)
			} else {
				err = api.CreateMemberView(createView)
			}

			if err != nil {
				log.Errorf("failed to create user: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.RoleID, "roleid", "", 0, "Role ID")
	flags.StringVarP(&opts.MemberUser.Username, "username", "", "", "Username")
	flags.StringVarP(&opts.MemberGroup.GroupName, "groupname", "", "", "Group Name")
	flags.StringVarP(&opts.MemberGroup.LdapGroupDn, "ldapdn", "", "", "DN of LDAP Group")
	flags.Int64VarP(&opts.MemberGroup.ID, "groupid", "", 0, "Group ID")
	flags.Int64VarP(&opts.MemberUser.UserID, "userid", "", 0, "User ID")
	flags.Int64VarP(&opts.MemberGroup.GroupType, "grouptype", "", 0, "Group Type")

	return cmd
}
