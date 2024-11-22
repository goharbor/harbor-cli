package user

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ElevateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elevate",
		Short: "elevate user",
		Long:  "elevate user to admin role",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var userId int64
			if len(args) > 0 {
				userId, _ = api.GetUsersIdByName(args[0])

			} else {
				userId = prompt.GetUserIdFromUser()
			}

			confirm, err := views.ConfirmElevation()
			if confirm {
				err = api.ElevateUser(userId)
			} else {
				log.Error("Permission denied for elevate user to admin.")
			}
			if err != nil {
				log.Errorf("failed to elevate user: %v", err)
			}

		},
	}

	return cmd
}
