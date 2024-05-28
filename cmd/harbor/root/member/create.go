package member

import (
	// "context"

	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/member"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"

	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/goharbor/harbor-cli/pkg/views/member/create"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateMemberCommand() *cobra.Command {
	var opts create.CreateView
	opts.MemberUser = &models.UserEntity{} // Initialize MemberUser
	opts.MemberGroup = &models.UserGroup{} // Initialize MemberGroup

	cmd := &cobra.Command{
		Use:   "create [ProjectName Or ID]",
		Short: "create member",
		Long:  "create member for the project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			createView := &create.CreateView{
				ProjectNameOrID: args[0],
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
				err = runCreateMember(*createView)
			} else {
				err = createMemberView(createView)
			}

			if err != nil {
				log.Errorf("failed to create user: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.RoleID, "roleID", "", 0, "Role ID")
	flags.StringVarP(&opts.MemberUser.Username, "username", "", "", "Username")
	flags.StringVarP(&opts.MemberGroup.GroupName, "groupname", "", "", "Group Name")
	flags.Int64VarP(&opts.MemberGroup.ID, "groupid", "", 0, "Group ID")
	flags.Int64VarP(&opts.MemberUser.UserID, "userID", "", 0, "User ID")

	return cmd
}

func createMemberView(createView *create.CreateView) error {
	create.CreateMemberView(createView)
	return runCreateMember(*createView)
}

func runCreateMember(opts create.CreateView) error {
	credentialName := viper.GetString("current-credential-name")

	client := utils.GetClientByCredentialName(credentialName)

	ctx := context.Background()

	response, err := client.Member.CreateProjectMember(
		ctx, &member.CreateProjectMemberParams{
			ProjectMember: &models.ProjectMember{
				RoleID:      int64(opts.RoleID + 1),
				MemberUser:  opts.MemberUser,
				MemberGroup: opts.MemberGroup,
			},
			ProjectNameOrID: opts.ProjectNameOrID,
		},
	)
	if err != nil {
		return err
	}

	if response != nil {
		log.Info("Member created successfully")
	}

	return nil
}
