package user

import (
	// "context"

	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"

	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/goharbor/harbor-cli/pkg/views/user/create"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UserCreateCmd() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create user",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			createView := &create.CreateView{
				Email:    opts.Email,
				Realname: opts.Realname,
				Comment:  opts.Comment,
				Password: opts.Password,
				Username: opts.Username,
			}

			if opts.Email != "" && opts.Realname != "" && opts.Comment != "" && opts.Password != "" && opts.Username != "" {
				err = runCreateUser(opts)
			} else {
				err = createUserView(createView)
			}

			if err != nil {
				log.Errorf("failed to create user: %v", err)
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Email, "email", "", "", "Email")
	flags.StringVarP(&opts.Realname, "realname", "", "", "Realname")
	flags.StringVarP(&opts.Comment, "comment", "", "", "Comment")
	flags.StringVarP(&opts.Password, "password", "", "", "Password")
	flags.StringVarP(&opts.Username, "username", "", "", "Username")

	return cmd
}

func createUserView(createView *create.CreateView) error {
	create.CreateUserView(createView)
	return runCreateUser(*createView)

}

func runCreateUser(opts create.CreateView) error {
	credentialName := viper.GetString("current-credential-name")

	client := utils.GetClientByCredentialName(credentialName)

	ctx := context.Background()

	response, err := client.User.CreateUser(ctx, &user.CreateUserParams{
		UserReq: &models.UserCreationReq{
			Email:    opts.Email,
			Realname: opts.Realname,
			Comment:  opts.Comment,
			Password: opts.Password,
			Username: opts.Username,
		},
	})

	if err != nil {
		return err
	}

	if response != nil {
		log.Info("User created successfully")
	}

	return nil
}
