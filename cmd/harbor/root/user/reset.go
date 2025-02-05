package user

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/user/reset"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UserResetCmd() *cobra.Command {
	var UserID int64
	cmd := &cobra.Command{
		Use:   "reset [username]",
		Short: "reset user's password",
		Long:  "Resets the password for a specific user by providing their username",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			resetView := &reset.ResetView{}
			flags := cmd.Flags()

			if len(args) > 0 {
				UserID, _ = api.GetUsersIdByName(args[0])
			} else {
				if !flags.Changed("userID") {
					UserID = prompt.GetUserIdFromUser()
				}
			}

			existingUser := api.GetUserProfileById(UserID)
			if existingUser == nil {
				logrus.Errorf("user is not found")
				return
			}

			reset.ResetUserView(resetView)

			err = api.ResetPassword(UserID, *resetView)

			if err != nil {
				logrus.Errorf("failed to reset user's password: %v", err)
			}

		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&UserID, "userID", "", -1, "ID of the user")

	return cmd

}
