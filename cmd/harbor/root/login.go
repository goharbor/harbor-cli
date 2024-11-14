package root

import (
	"context"
	"fmt"
	"os"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/login"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	serverAddress string
	Username      string
	Password      string
	Name          string
	passwordStdin bool
)

// LoginCommand creates a new `harbor login` command
func LoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [server]",
		Short: "Log in to Harbor registry",
		Long:  "Authenticate with Harbor Registry.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				serverAddress = args[0]
			}

			if passwordStdin {
				fmt.Print("Password: ")
				passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					return fmt.Errorf("failed to read password from stdin: %v", err)
				}
				fmt.Println()
				Password = string(passwordBytes)
			}

			loginView := login.LoginView{
				Server:   serverAddress,
				Username: Username,
				Password: Password,
				Name:     Name,
			}

			// autogenerate name
			if loginView.Name == "" && loginView.Server != "" && loginView.Username != "" {
				loginView.Name = fmt.Sprintf("%s@%s", loginView.Username, utils.SanitizeServerAddress(loginView.Server))
			}

			var err error

			if loginView.Server != "" && loginView.Username != "" && loginView.Password != "" {
				err = runLogin(loginView)
			} else {
				err = createLoginView(&loginView)
			}

			if err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&Name, "name", "", "", "name for the set of credentials")
	flags.StringVarP(&Username, "username", "u", "", "Username")
	flags.StringVarP(&Password, "password", "p", "", "Password")
	flags.BoolVar(&passwordStdin, "password-stdin", false, "Take the password from stdin")

	return cmd
}

func createLoginView(loginView *login.LoginView) error {
	if loginView == nil {
		loginView = &login.LoginView{
			Server:   "",
			Username: "",
			Password: "",
			Name:     "",
		}
	}
	login.CreateView(loginView)

	return runLogin(*loginView)
}

func runLogin(opts login.LoginView) error {
	opts.Server = utils.FormatUrl(opts.Server)

	clientConfig := &harbor.ClientSetConfig{
		URL:      opts.Server,
		Username: opts.Username,
		Password: opts.Password,
	}
	client := utils.GetClientByConfig(clientConfig)

	ctx := context.Background()
	_, err := client.User.GetCurrentUserInfo(ctx, &user.GetCurrentUserInfoParams{})
	if err != nil {
		return fmt.Errorf("login failed, please check your credentials: %s", err)
	}

	cred := utils.Credential{
		Name:          opts.Name,
		Username:      opts.Username,
		Password:      opts.Password,
		ServerAddress: opts.Server,
	}

	if err = utils.AddCredentialsToConfigFile(cred, utils.DefaultConfigPath); err != nil {
		return fmt.Errorf("failed to store the credential: %s", err)
	}
	return nil
}
