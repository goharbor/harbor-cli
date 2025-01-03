package user

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/user/create"
)

func UserCreateCmd() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create user",
		Args:  cobra.ExactArgs(0),
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
				err = api.CreateUser(opts)
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
	return api.CreateUser(*createView)

}
