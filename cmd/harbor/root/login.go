// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package root

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/ping"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/login"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	serverAddress    string
	Username         string
	Password         string
	Name             string
	passwordStdin    bool
	skipVerifyClient bool
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
				passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd())) // #nosec G115 - fd fits in int on all supported platforms
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

			var err error
			var config *utils.HarborConfig
			config, err = utils.GetCurrentHarborConfig()
			if err != nil {
				return fmt.Errorf("failed to get current harbor config: %s", err)
			}
			if err := ProcessLogin(loginView, config); err != nil {
				return fmt.Errorf("login failed: %w", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&Username, "username", "u", "", "Username")
	flags.StringVarP(&Name, "context-name", "n", "", "Login context name (optional)")
	flags.StringVarP(&Password, "password", "p", "", "Password")
	flags.BoolVar(&passwordStdin, "password-stdin", false, "Take the password from stdin")
	flags.BoolVarP(&skipVerifyClient, "skip-verify-client", "", false, "Skip whether the clients basic auth credentials shall be validated against the Harbor server during login. This is not recommended as it may lead to storing invalid credentials. Use this flag if you want to skip validation of credentials during login, for example, when the Harbor server is not reachable at the moment of login but you still want to store the credentials for later use.")

	cmd.MarkFlagsMutuallyExclusive("password", "password-stdin")

	return cmd
}

// ProcessLogin applies a simplified decision logic to run login or launch an interactive view.
func ProcessLogin(loginView login.LoginView, config *utils.HarborConfig) error {
	// Auto-generate the name if not provided.
	if loginView.Name == "" && loginView.Server != "" && loginView.Username != "" {
		loginView.Name = fmt.Sprintf("%s@%s", loginView.Username, utils.SanitizeServerAddress(loginView.Server))
	}

	// If complete credentials are provided (overrides), run login using them directly.
	if loginView.Server != "" && loginView.Username != "" && loginView.Password != "" {
		return RunLogin(loginView)
	}

	// If nothing matches, launch the interactive view.
	return CreateLoginView(&loginView)
}

// CreateLoginView launches the interactive login view.
// In this implementation, it calls login.CreateView and then tries to run login.
func CreateLoginView(loginView *login.LoginView) error {
	if loginView == nil {
		loginView = &login.LoginView{
			Server:   "",
			Username: "",
			Password: "",
			Name:     "",
		}
	}
	login.CreateView(loginView)

	return RunLogin(*loginView)
}

// RunLogin attempts to log in using the provided LoginView credentials.
func RunLogin(opts login.LoginView) error {
	opts.Server = utils.FormatUrl(opts.Server)

	clientConfig := &harbor.ClientSetConfig{
		URL:      opts.Server,
		Username: opts.Username,
		Password: opts.Password,
	}
	err := utils.ValidateURL(opts.Server)
	if err != nil {
		return fmt.Errorf("invalid server URL: %w", err)
	}
	client := utils.GetClientByConfig(clientConfig)

	if !skipVerifyClient {
		if err := validateClientConnection(client); err != nil {
			return err
		}
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
				fmt.Printf("Login successful for %s at %s\n", opts.Username, opts.Server)
				return nil
			} else {
				log.Warn("Credentials already exist in the config file but the password is different. Updating the password.")
				if err = utils.UpdateCredentialsInConfigFile(cred, configPath); err != nil {
					log.Fatalf("failed to update the credential: %s", err)
				}
				fmt.Printf("Login successful for %s at %s\n", opts.Username, opts.Server)
				return nil
			}
		} else {
			log.Warn("Credentials already exist in the config file but more than one field was different. Updating the credentials.")
			if err = utils.UpdateCredentialsInConfigFile(cred, configPath); err != nil {
				log.Fatalf("failed to update the credential: %s", err)
			}
			fmt.Printf("Login successful for %s at %s\n", opts.Username, opts.Server)
			return nil
		}
	}

	if err = utils.AddCredentialsToConfigFile(cred, configPath); err != nil {
		return fmt.Errorf("failed to store the credential: %s", err)
	}
	log.Debugf("Credentials successfully added to the config file.")
	fmt.Printf("Login successful for %s at %s\n", opts.Username, opts.Server)
	return nil
}

func validateClientConnection(client *client.HarborAPI) error {
	ctx := context.Background()

	// Primary check: GetCurrentUserInfo requires auth → 401 for bad creds.
	_, err := client.User.GetCurrentUserInfo(ctx, &user.GetCurrentUserInfoParams{})
	if err == nil {
		return nil
	}

	errorCode := utils.ParseHarborErrorCode(err)
	// 401/403 = definite auth failure
	if errorCode == "401" || errorCode == "403" {
		return fmt.Errorf("authentication failed, check your credentials: %v", utils.ParseHarborErrorMsg(err))
	}

	// For other errors (e.g. 412 for robot/OIDC accounts, 5xx),
	// fall back to secondary endpoints to verify creds and reachability.
	_, projectErr := client.Project.ListProjects(ctx, &project.ListProjectsParams{
		Page:     new(int64(1)),
		PageSize: new(int64(1)),
	})
	_, pingErr := client.Ping.GetPing(ctx, &ping.GetPingParams{})

	// If either secondary check returns 401/403, creds are bad.
	if projectErr != nil {
		projCode := utils.ParseHarborErrorCode(projectErr)
		if projCode == "401" || projCode == "403" {
			return fmt.Errorf("authentication failed, check your credentials: %v", utils.ParseHarborErrorMsg(projectErr))
		}
	}

	// Both passed → creds valid, server reachable
	if projectErr == nil && pingErr == nil {
		return nil
	}

	// Build diagnostic message
	var results []string
	if projectErr != nil {
		results = append(results, fmt.Sprintf("ListProjects failed: %v", projectErr))
	} else {
		results = append(results, "ListProjects succeeded")
	}
	if pingErr != nil {
		results = append(results, fmt.Sprintf("Ping failed: %v", pingErr))
	} else {
		results = append(results, "Ping succeeded")
	}
	return fmt.Errorf("server error (status %s): %v (%s)", errorCode, utils.ParseHarborErrorMsg(err), strings.Join(results, "; "))
}
