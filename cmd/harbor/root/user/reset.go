package user

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/user/reset"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UserResetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "reset [username]",
		Short: "reset user's password",
		Long:  "Resets the password for a specific user by providing their username",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var UserID int64
			resetView := &reset.ResetView{}

			if len(args) > 0 {
				UserID, err = api.GetUsersIdByName(args[0])
				if err != nil {
					logrus.Error(err)
					return
				}
			} else {
				UserID = prompt.GetUserIdFromUser()
			}

			reset.ResetUserView(resetView)

			err = api.ResetPassword(UserID, *resetView)

			if err != nil {
				logrus.Errorf("failed to reset user's password: %v", err)
			}

		},
	}

	return cmd

}
