package user

import (
	"context"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UserDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete user",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				userId, _ := strconv.ParseInt(args[0], 10, 64)
				err = runDeleteUser(userId)

			} else {
				userId := utils.GetUserIdFromUser()
				err = runDeleteUser(userId)
			}

			if err != nil {
				log.Errorf("failed to delete user: %v", err)
			}

		},
	}

	return cmd

}

func runDeleteUser(userId int64) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	_, err := client.User.DeleteUser(ctx, &user.DeleteUserParams{UserID: userId})
	if err != nil {
		return err
	}
	log.Info("user deleted successfully")
	return nil
}
