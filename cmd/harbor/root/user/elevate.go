package user

import (
	"context"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
				userId, _ = strconv.ParseInt(args[0], 10, 64)

			} else {
				userId = utils.GetUserIdFromUser()
			}

			// Todo : Ask for the confirmation before elevating the user to admin role

			err = runElevateUser(userId)

			if err != nil {
				log.Errorf("failed to elevate user: %v", err)
			}

		},
	}

	return cmd
}

func runElevateUser(userId int64) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	UserSysAdminFlag := &models.UserSysAdminFlag{
		SysadminFlag: true,
	}
	_, err := client.User.SetUserSysAdmin(ctx, &user.SetUserSysAdminParams{UserID: userId, SysadminFlag: UserSysAdminFlag})
	if err != nil {
		return err
	}
	log.Info("user elevated role to admin successfully")
	return nil
}
