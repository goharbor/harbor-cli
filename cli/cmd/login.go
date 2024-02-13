package cmd

import (
	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/spf13/cobra"
)

// LoginCommand creates a new `harbor login` command
func LoginCommand() *cobra.Command {
	var opts config.LoginOptions

	cmd := &cobra.Command{
		Use:   "login [SERVER]",
		Short: "Log in to Harbor registry",
		Long:  "Authenticate with Harbor Registry.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ServerAddress = args[0]
			return api.RunLogin(&opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "name for the set of credentials")
	flags.StringVarP(&opts.ServerAddress, "url", "", "", "server address")

	flags.StringVarP(&opts.Username, "username", "u", "", "Username")
	if err := cmd.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
	flags.StringVarP(&opts.Password, "password", "p", "", "Password")
	if err := cmd.MarkFlagRequired("password"); err != nil {
		panic(err)
	}

	return cmd
}