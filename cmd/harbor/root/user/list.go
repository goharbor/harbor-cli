package user

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/user/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UserListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list users",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			response, err := api.ListUsers()
			if err != nil {
				log.Errorf("failed to list users: %v", err)
				return
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(response, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.ListUsers(response.Payload)
			}
		},
	}

	return cmd

}
