package user

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/password/reset"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UserPasswordChangeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "password",
		Short: "Reset user password by name or id",
		Long:  "Allows admin to reset the password for a specified user or select interactively if no username is provided.",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var userID int64
			var err error
			var opts reset.PasswordChangeView

			if len(args) > 0 {
				userID, err = api.GetUsersIdByName(args[0])
				if err != nil {
					log.Errorf("Failed to get user ID for '%s': %v", args[0], err)
					return
				}
			} else {
				userID = prompt.GetUserIdFromUser()
			}

			reset.ChangePasswordView(&opts)

			if err := api.ResetPassword(userID, &opts); err != nil {
				if isUnauthorizedError(err) {
					fmt.Println("Permission denied: Admin privileges are required to execute this command.")
				} else {
					fmt.Printf("Failed to change password for user ID %d: %v", userID, err)
				}
			} else {
				fmt.Printf("Password successfully changed for user ID %d", userID)
			}

		},
	}
}
