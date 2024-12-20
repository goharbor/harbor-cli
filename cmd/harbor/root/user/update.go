package user

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/user/update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateUserCmd() *cobra.Command {

	opts := &models.UserProfile{}

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "update user's profile",
		Long:    "update user's profile by username",
		Example: "harbor user update [username]",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var UserId int64

			if len(args) > 0 {
				UserId, err = api.GetUsersIdByName(args[0])
				if err != nil {
					logrus.Errorf("fail to get user id by username: %s", args[0])
					return
				}
			} else {
				UserId = prompt.GetUserIdFromUser()
			}

			existingUser := api.GetUserProfileById(UserId)
			if existingUser == nil {
				logrus.Errorf("user is not found")
				return
			}

			updateView := &models.UserProfile{
				Comment:  existingUser.Comment,
				Email:    existingUser.Email,
				Realname: existingUser.Realname,
			}

			flags := cmd.Flags()
			if flags.Changed("comment") {
				updateView.Comment = opts.Comment
			}
			if flags.Changed("email") {
				updateView.Email = opts.Email
			}
			if flags.Changed("realname") {
				updateView.Realname = opts.Realname
			}

			update.UpdateUserProfileView(updateView)
			err = api.UpdateUserProfile(updateView, UserId)
			if err != nil {
				logrus.Errorf("failed to update user's profile: %v", err)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Comment, "comment", "m", "", "Comment of the user")
	flags.StringVarP(&opts.Email, "email", "e", "", "Email of the user")
	flags.StringVarP(&opts.Realname, "realname", "r", "", "Realname of the user")

	return cmd
}
