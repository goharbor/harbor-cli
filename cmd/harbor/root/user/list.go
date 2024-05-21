package user

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/user/list"
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
			response := runListUsers()
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(response.Payload)
			} else {
				list.ListUsers(response.Payload)
			}
		},
	}

	return cmd

}

func runListUsers() *user.ListUsersOK {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, _ := client.User.ListUsers(ctx, &user.ListUsersParams{})

	return response
}
