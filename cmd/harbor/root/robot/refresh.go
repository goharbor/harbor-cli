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
package robot

import (
	"fmt"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func RefreshSecretCommand() *cobra.Command {
	var (
		robotID     int64
		secret      string
		secretStdin bool
	)
	cmd := &cobra.Command{
		Use:   "refresh [robotID]",
		Short: "refresh robot secret by id",
		Long: `Refresh the secret for an existing robot account in Harbor.

This command generates a new secret for a robot account, effectively revoking 
the old secret and requiring updates to any systems using the robot's credentials.

The command supports multiple ways to identify the robot account:
- By providing the robot ID directly as an argument
- Without any arguments, which will prompt for both project and robot selection

You can specify the new secret in several ways:
- Let Harbor generate a random secret (default)
- Provide a custom secret with the --secret flag
- Pipe a secret via stdin using the --secret-stdin flag

After refreshing, the new secret will be:
- Displayed on screen
- Copied to clipboard for immediate use
- Usable immediately for authentication

Important considerations:
- The old secret will be invalidated immediately
- Any systems using the old credentials will need to be updated
- There is no way to recover the old secret after refreshing

Examples:
  # Refresh robot secret by ID (generates a random secret)
  harbor-cli project robot refresh 123

  # Refresh with a custom secret
  harbor-cli project robot refresh 123 --secret "MyCustomSecret123"

  # Provide secret via stdin (useful for scripting)
  echo "MySecretFromScript123" | harbor-cli project robot refresh 123 --secret-stdin

  # Interactive refresh (will prompt for project and robot selection)
  harbor-cli project robot refresh`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) == 1 {
				robotID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Fatalf("failed to parse robot ID: %v", err)
				}
			} else {
				robotID, err = prompt.GetRobotIDFromUser(-1)
				if err != nil {
					log.Fatalf("failed to get robot ID from user: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			if secret != "" {
				err = utils.ValidatePassword(secret)
				if err != nil {
					log.Fatalf("Invalid secret: %v\n", err)
				}
			}
			if secretStdin {
				secret = getSecret()
			}

			response, err := api.RefreshSecret(secret, robotID)
			if err != nil {
				errorCode := utils.ParseHarborErrorCode(err)
				if errorCode == "403" {
					log.Fatalf("Permission denied: (Project) Admin privileges are required to execute this command.\n")
				} else {
					log.Fatalf("failed to refresh robot secret: %v\n", utils.ParseHarborErrorMsg(err))
				}
			}

			log.Info("Secret updated successfully.")

			if response.Payload.Secret != "" {
				secret = response.Payload.Secret
				create.CreateRobotSecretView("", secret)

				err = clipboard.WriteAll(response.Payload.Secret)
				if err != nil {
					log.Fatalf("failed to write the secret to the clipboard: %v", err)
				}
				fmt.Println("secret copied to clipboard.")
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&secret, "secret", "", "", "secret")
	flags.BoolVarP(&secretStdin, "secret-stdin", "", false, "Take the robot secret from stdin")

	return cmd
}

// getSecret from commandline
func getSecret() string {
	secret, err := utils.GetSecretStdin("Enter your secret: ")
	if err != nil {
		log.Fatalf("Error reading secret: %v\n", err)
	}

	if err := utils.ValidatePassword(secret); err != nil {
		log.Fatalf("Invalid secret: %v\n", err)
	}
	return secret
}
