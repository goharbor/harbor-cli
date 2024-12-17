package user

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/user/reset"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UserResetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "reset [username]",
		Short: "reset user's password",
		Long:  "reset user's password by username",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var userId int64
			resetView := &reset.ResetView{}

			if len(args) > 0 {
				userId, _ = api.GetUsersIdByName(args[0])
			} else {
				userId = prompt.GetUserIdFromUser()
			}

			reset.ResetUserView(resetView)

			err = api.ResetPassword(userId, *resetView)

			if err != nil {
				log.Errorf("failed to reset user's password: %v", err)
			}

		},
	}

	return cmd

}
