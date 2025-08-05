package root

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/password/change"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func PasswordCommand() *cobra.Command {
	var opts change.PasswordChangeView

	cmd := &cobra.Command{
		Use:   "password",
		Short: "Change your password",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			change.ChangePasswordView(&opts)

			err := UpdatePassword(&opts)
			if err != nil {
				log.Errorf("Error changing password: %v", err)
				fmt.Printf("Error changing password: %v\n", err)
				return
			}
			log.Info("Password updated successfully.")
			fmt.Println("Password updated successfully.")
		},
	}

	return cmd
}

func UpdatePassword(opts *change.PasswordChangeView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	userResp, err := client.User.GetCurrentUserInfo(ctx, &user.GetCurrentUserInfoParams{})
	if err != nil {
		return err
	}
	userId := userResp.Payload.UserID
	response, err := client.User.UpdateUserPassword(ctx, &user.UpdateUserPasswordParams{
		Password: &models.PasswordReq{
			OldPassword: opts.OldPassword,
			NewPassword: opts.NewPassword,
		},
		UserID: userId,
	})
	if err != nil {
		return err
	}

	if response != nil {
		log.Infof("Password change successful")
	}

	return nil
}
