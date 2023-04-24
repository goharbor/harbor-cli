package login

import (
	"context"
	"fmt"

	"github.com/akshatdalton/harbor-cli/cmd/utils"
	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/spf13/cobra"
)

type loginOptions struct {
	name          string
	serverAddress string
	username      string
	password      string
}

// NewLoginCommand creates a new `harbor login` command
func NewLoginCommand() *cobra.Command {
	var opts loginOptions

	cmd := &cobra.Command{
		Use:   "login [SERVER]",
		Short: "Log in to Harbor registry",
		Long:  "Authenticate with Harbor Registry.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.serverAddress = args[0]
			return runLogin(opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.name, "name", "", "", "name for the set of credentials")

	flags.StringVarP(&opts.username, "username", "u", "", "Username")
	if err := cmd.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
	flags.StringVarP(&opts.password, "password", "p", "", "Password")
	if err := cmd.MarkFlagRequired("password"); err != nil {
		panic(err)
	}

	return cmd
}

func runLogin(opts loginOptions) error {
	clientConfig := &harbor.ClientSetConfig{
		URL:      opts.serverAddress,
		Username: opts.username,
		Password: opts.password,
	}
	client := utils.GetClientByConfig(clientConfig)

	ctx := context.Background()
	_, err := client.User.GetCurrentUserInfo(ctx, &user.GetCurrentUserInfoParams{})
	if err != nil {
		return fmt.Errorf("login failed, please check your credentials: %s", err)
	}

	cred := utils.Credential{
		Name:          opts.name,
		Username:      opts.username,
		Password:      opts.password,
		ServerAddress: opts.serverAddress,
	}

	if err = utils.StoreCredential(cred, true); err != nil {
		return fmt.Errorf("Failed to store the credential: %s", err)
	}
	return nil
}
