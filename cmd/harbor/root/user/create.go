package user

import (
	// "context"

	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"

	// "github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// {
// 	"email": "string",
// 	"realname": "string",
// 	"comment": "string",
// 	"password": "string",
// 	"username": "string"
//   }

type createUserOptions struct {
	email    string
	realname string
	comment  string
	password string
	username string
}

func UserCreateCmd() *cobra.Command {
	var opts createUserOptions

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create user",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)

			if err != nil {
				return err
			}

			return runCreateUser(credentialName, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.email, "email", "", "", "Email")
	flags.StringVarP(&opts.realname, "realname", "", "", "Realname")
	flags.StringVarP(&opts.comment, "comment", "", "", "Comment")
	flags.StringVarP(&opts.password, "password", "", "", "Password")
	flags.StringVarP(&opts.username, "username", "", "", "Username")

	return cmd
}

func runCreateUser(credentialName string, opts createUserOptions) error {

	client := utils.GetClientByCredentialName(credentialName)

	ctx := context.Background()

	response, err := client.User.CreateUser(ctx, &user.CreateUserParams{
		UserReq: &models.UserCreationReq{
			Email:    opts.email,
			Realname: opts.realname,
			Comment:  opts.comment,
			Password: opts.password,
			Username: opts.username,
		},
	})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
