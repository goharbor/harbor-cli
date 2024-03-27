package user

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	// "github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func UserListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list users",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			credentialName, _ := cmd.Flags().GetString(constants.CredentialNameOption)
			runListUsers(credentialName)

		},
	}

	return cmd

}

func runListUsers(credentialName string) {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, _ := client.User.ListUsers(ctx, &user.ListUsersParams{})

	utils.PrintPayloadInJSONFormat(response)
}
