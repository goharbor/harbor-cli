package root

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/login"
	log "github.com/sirupsen/logrus"
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

			// Check whether there is already a config with credentials
			var err error
			var config *utils.HarborConfig
			var creds []utils.Credential
			config, err = utils.GetCurrentHarborConfig()
			if err != nil {
				return fmt.Errorf("failed to get current harbor config: %s", err)
			}
			currentCredentialName := config.CurrentCredentialName

			if loginView.Server != "" && loginView.Username != "" && loginView.Password != "" {
				err = runLogin(loginView)
				// check whether the user entered a new credential than already in config
			} else if currentCredentialName != "" && loginView.Name != "" && loginView.Name != currentCredentialName {
				var resolvedLoginView login.LoginView
				creds = config.Credentials
				for _, cred := range creds {
					if cred.Name == loginView.Name {
						resolvedLoginView = login.LoginView{
							Server:   cred.ServerAddress,
							Username: cred.Username,
							Password: cred.Password,
							Name:     cred.Name,
						}
					}
				}
				if resolvedLoginView.Server != "" && resolvedLoginView.Username != "" && resolvedLoginView.Password != "" {
					err = runLogin(resolvedLoginView)
				} else {
					err = createLoginView(&loginView)
				}
			} else if currentCredentialName != "" && loginView.Name == "" && loginView.Server == "" && loginView.Username == "" && loginView.Password == "" {
				var resolvedLoginView login.LoginView
				creds := config.Credentials
				key, err := utils.GetEncryptionKey()
				if err != nil {
					return fmt.Errorf("failed to get encryption key: %w", err)
				}
				var targetCred *utils.Credential
				for _, cred := range creds {
					if strings.EqualFold(cred.Name, currentCredentialName) {
						targetCred = &cred
						break
					}
				}
				if targetCred == nil {
					return fmt.Errorf("no matching credential found for '%s'", currentCredentialName)
				}
				decryptedPassword, err := utils.Decrypt(key, string(targetCred.Password))
				if err != nil {
					return fmt.Errorf("failed to decrypt password: %w", err)
				}
				resolvedLoginView = login.LoginView{
					Server:   targetCred.ServerAddress,
					Username: targetCred.Username,
					Password: decryptedPassword,
					Name:     targetCred.Name,
				}
				err = runLogin(resolvedLoginView)
				if err != nil {
					return fmt.Errorf("failed to run login: %w", err)
				}
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
	if err := utils.GenerateEncryptionKey(); err != nil {
		fmt.Println("Encryption key already exists or could not be created:", err)
	}

	key, err := utils.GetEncryptionKey()
	if err != nil {
		fmt.Println("Error getting encryption key:", err)
		return fmt.Errorf("failed to get encryption key: %s", err)
	}

	encryptedPassword, err := utils.Encrypt(key, []byte(opts.Password))
	if err != nil {
		fmt.Println("Error encrypting password:", err)
		return fmt.Errorf("failed to encrypt password: %s", err)
	}

	cred := utils.Credential{
		Name:          opts.Name,
		Username:      opts.Username,
		Password:      encryptedPassword,
		ServerAddress: opts.Server,
	}
	harborData, err := utils.GetCurrentHarborData()
	if err != nil {
		return fmt.Errorf("failed to get current harbor data: %s", err)
	}
	configPath := harborData.ConfigPath
	log.Debugf("Checking if credentials already exist in the config file...")
	existingCred, err := utils.GetCredentials(opts.Name)
	if err == nil {
		if existingCred.Username == opts.Username && existingCred.ServerAddress == opts.Server {
			if existingCred.Password == encryptedPassword {
				log.Warn("Credentials already exist in the config file. They were not added again.")
				return nil
			} else {
				log.Warn("Credentials already exist in the config file but the password is different. Updating the password.")
				if err = utils.UpdateCredentialsInConfigFile(cred, configPath); err != nil {
					log.Fatalf("failed to update the credential: %s", err)
				}
				return nil
			}
		} else {
			log.Warn("Credentials already exist in the config file but more than one field was different. Updating the credentials.")
			if err = utils.UpdateCredentialsInConfigFile(cred, configPath); err != nil {
				log.Fatalf("failed to update the credential: %s", err)
			}
			return nil
		}
	}

	if err = utils.AddCredentialsToConfigFile(cred, configPath); err != nil {
		return fmt.Errorf("failed to store the credential: %s", err)
	}
	log.Debugf("Credentials successfully added to the config file.")
	return nil
}
