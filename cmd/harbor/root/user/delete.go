package user

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UserDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete user",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				userName, _ := api.GetUsersIdByName(args[0])
				err = api.DeleteUser(userName)

			} else {
				userId := prompt.GetUserIdFromUser()
				err = api.DeleteUser(userId)
			}

			if err != nil {
				log.Errorf("failed to delete user: %v", err)
			}

		},
	}

	return cmd

}
