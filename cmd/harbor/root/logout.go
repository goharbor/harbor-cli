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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// LogoutCommand creates a new `harbor logout` command
func LogoutCommand() *cobra.Command {
	var skipConfirm bool

	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "Log out from Harbor registry",
		Long:    "Remove the current credential from the local CLI config.",
		Args:    cobra.NoArgs,
		Example: `  harbor logout`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				return fmt.Errorf("failed to get current harbor config: %s", err)
			}

			if config.CurrentCredentialName == "" {
				fmt.Println("Not logged in.")
				return nil
			}

			currentName := config.CurrentCredentialName

			if !skipConfirm {
				fmt.Printf("Log out from '%s'? [y/N]: ", currentName)

				reader := bufio.NewReader(os.Stdin)
				answer, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read user input: %v", err)
				}

				answer = strings.TrimSpace(strings.ToLower(answer))
				if answer != "y" && answer != "yes" {
					fmt.Println("Logout cancelled.")
					return nil
				}
			}

			for i, cred := range config.Credentials {
				if cred.Name == currentName {
					config.Credentials = append(config.Credentials[:i], config.Credentials[i+1:]...)
					break
				}
			}

			config.CurrentCredentialName = ""

			if err := utils.UpdateConfigFile(config); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Logged out from '%s'.\n", currentName)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
