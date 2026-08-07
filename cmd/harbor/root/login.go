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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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
	authMode         string
)

const (
	cliAuthModeDB   = "db"
	cliAuthModeLDAP = "ldap"
	cliAuthModeOIDC = "oidc"

	harborAuthModeDB   = "db_auth"
	harborAuthModeLDAP = "ldap_auth"
	harborAuthModeOIDC = "oidc_auth"
)

func resetLoginOptions() {
	serverAddress = ""
	Username = ""
	Password = ""
	Name = ""
	passwordStdin = false
	skipVerifyClient = false
	authMode = ""
}

// LoginCommand creates a new `harbor login` command
func LoginCommand() *cobra.Command {
	resetLoginOptions()

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

			effectiveMode, err := resolveLoginAuthMode(loginView)
			if err != nil {
				return err
			}
			if effectiveMode == cliAuthModeOIDC {
				return RunOIDCLogin(serverAddress)
			}

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
	flags.StringVarP(&Password, "password", "p", "", "Password (not recommended, use --password-stdin for better security)")
	flags.BoolVar(&passwordStdin, "password-stdin", false, "Take the password from stdin")
	flags.StringVar(&authMode, "auth-mode", "", "Authentication mode (db, ldap, oidc)")
	flags.BoolVarP(&skipVerifyClient, "skip-verify-client", "", false, "Skip whether the clients basic auth credentials shall be validated against the Harbor server during login. This is not recommended as it may lead to storing invalid credentials. Use this flag if you want to skip validation of credentials during login, for example, when the Harbor server is not reachable at the moment of login but you still want to store the credentials for later use.")

	cmd.MarkFlagsMutuallyExclusive("password", "password-stdin")

	return cmd
}

func resolveLoginAuthMode(loginView login.LoginView) (string, error) {
	requestedMode, err := normalizeCLIAuthMode(authMode)
	if err != nil {
		return "", err
	}
	log.Debugf("resolving login auth mode for server=%q requested_mode=%q username_provided=%t password_provided=%t", loginView.Server, requestedMode, loginView.Username != "", loginView.Password != "")

	if requestedMode == "" {
		if loginView.Server == "" {
			return "", nil
		}
		if loginView.Username != "" || loginView.Password != "" {
			return "", nil
		}

		harborMode, err := getHarborAuthMode(loginView.Server)
		if err != nil {
			return "", fmt.Errorf("unable to determine Harbor auth_mode from /api/v2.0/systeminfo. Please retry with --auth-mode db, --auth-mode ldap, or --auth-mode oidc: %w", err)
		}
		log.Debugf("auto-detected Harbor auth_mode=%q for server=%q", harborMode, loginView.Server)
		switch harborMode {
		case harborAuthModeDB:
			log.Debug("selected CLI auth mode db from Harbor auth_mode db_auth")
			return cliAuthModeDB, nil
		case harborAuthModeLDAP:
			log.Debug("selected CLI auth mode ldap from Harbor auth_mode ldap_auth")
			return cliAuthModeLDAP, nil
		case harborAuthModeOIDC:
			log.Debug("selected CLI auth mode oidc from Harbor auth_mode oidc_auth")
			return cliAuthModeOIDC, nil
		default:
			return "", fmt.Errorf("unsupported Harbor auth_mode %q returned by /api/v2.0/systeminfo", harborMode)
		}
	}

	if loginView.Server == "" {
		return "", fmt.Errorf("server address is required when --auth-mode is set")
	}

	if requestedMode == cliAuthModeOIDC && (loginView.Username != "" || loginView.Password != "") {
		return "", fmt.Errorf("--auth-mode oidc cannot be used with --username, --password, or --password-stdin")
	}

	harborMode, err := getHarborAuthMode(loginView.Server)
	if err != nil {
		return "", fmt.Errorf("failed to determine Harbor auth_mode for --auth-mode %s: %w", requestedMode, err)
	}
	if err := validateAuthModeCombination(requestedMode, harborMode); err != nil {
		return "", err
	}
	log.Debugf("validated explicit CLI auth mode %q against Harbor auth_mode %q", requestedMode, harborMode)

	return requestedMode, nil
}

func normalizeCLIAuthMode(mode string) (string, error) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	switch mode {
	case "":
		return "", nil
	case cliAuthModeDB, cliAuthModeLDAP, cliAuthModeOIDC:
		return mode, nil
	default:
		return "", fmt.Errorf("invalid auth mode %q. Valid values are: db, ldap, oidc", mode)
	}
}

func validateAuthModeCombination(requestedMode, harborMode string) error {
	switch requestedMode {
	case cliAuthModeOIDC:
		if harborMode != harborAuthModeOIDC {
			return fmt.Errorf("OIDC login is not available because Harbor auth_mode is %s", harborMode)
		}
	case cliAuthModeLDAP:
		if harborMode != harborAuthModeLDAP {
			return fmt.Errorf("LDAP login is not available because Harbor auth_mode is %s", harborMode)
		}
	case cliAuthModeDB:
		if harborMode == harborAuthModeLDAP {
			return fmt.Errorf("DB login is not available because Harbor auth_mode is %s", harborMode)
		}
	default:
		return fmt.Errorf("unsupported auth mode %q", requestedMode)
	}
	return nil
}

func getHarborAuthMode(server string) (string, error) {
	server = utils.FormatUrl(server)
	if err := utils.ValidateURL(server); err != nil {
		return "", fmt.Errorf("invalid server URL: %w", err)
	}

	endpoint, err := joinServerPath(server, "/api/v2.0/systeminfo")
	if err != nil {
		return "", err
	}
	log.Debugf("querying Harbor system info at %s", endpoint)

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Get(endpoint) //nolint:gosec // endpoint is user-provided Harbor server URL.
	if err != nil {
		return "", fmt.Errorf("failed to query Harbor system info: %w", err)
	}
	defer resp.Body.Close()
	log.Debugf("received Harbor system info response status %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("unexpected status %d from /api/v2.0/systeminfo: %s", resp.StatusCode, string(body))
	}

	var payload struct {
		AuthMode string `json:"auth_mode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("failed to decode Harbor system info response: %w", err)
	}
	if payload.AuthMode == "" {
		return "", fmt.Errorf("Harbor system info response is missing auth_mode")
	}
	log.Debugf("Harbor system info reported auth_mode=%q", payload.AuthMode)

	return payload.AuthMode, nil
}

func joinServerPath(serverAddress, path string) (string, error) {
	u, err := url.Parse(serverAddress)
	if err != nil {
		return "", fmt.Errorf("failed to parse server URL: %w", err)
	}
	basePath := u.Path
	for len(basePath) > 0 && basePath[len(basePath)-1] == '/' {
		basePath = basePath[:len(basePath)-1]
	}
	u.Path = basePath + path
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
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

func RunOIDCLogin(serverAddress string) error {
	if serverAddress == "" {
		return fmt.Errorf("server address is required for OIDC login")
	}
	serverAddress = utils.FormatUrl(serverAddress)
	if err := utils.ValidateURL(serverAddress); err != nil {
		return fmt.Errorf("invalid server URL: %w", err)
	}
	log.Debugf("starting Harbor CLI OIDC login for server %s", serverAddress)

	loginResp, err := utils.InitiateOIDCLogin(serverAddress)
	if err != nil {
		return err
	}
	log.Debug("received Harbor CLI OIDC login redirect URL and poll token")

	fmt.Printf("Open this URL in your browser to authenticate:\n\n  %s\n\n", loginResp.RedirectURL)
	fmt.Print("Waiting for authentication...\n")

	tokenResp, err := utils.PollForOIDCToken(serverAddress, loginResp.PollToken, 10*time.Minute)
	if err != nil {
		return err
	}
	log.Debugf("Harbor CLI OIDC login completed for user %s", tokenResp.Username)

	harborData, err := utils.GetCurrentHarborData()
	if err != nil {
		return fmt.Errorf("failed to get current harbor data: %w", err)
	}

	if err := utils.AddOIDCCredentials(serverAddress, tokenResp.Username, tokenResp.IDToken, tokenResp.RefreshToken, tokenResp.ExpiresAt, harborData.ConfigPath); err != nil {
		return fmt.Errorf("failed to store OIDC credential: %w", err)
	}

	fmt.Printf("Login successful for %s at %s\n", tokenResp.Username, serverAddress)
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
