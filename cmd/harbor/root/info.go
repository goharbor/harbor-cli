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
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/systeminfo"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/cmd/harbor/internal/version"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show the current credential information",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentCredential := viper.GetString("current-credential-name")
			if currentCredential == "" {
				return fmt.Errorf("no active credentials found")
			}

			var registryAddress string
			creds := viper.Get("credentials").([]interface{})
			for _, cred := range creds {
				c := cred.(map[string]interface{})
				if c["name"] == currentCredential {
					registryAddress = c["serveraddress"].(string)
					break
				}
			}

			if registryAddress == "" {
				return fmt.Errorf("registry address not found for current credential: %s", currentCredential)
			}

			ctx, client, err := utils.ContextWithClient()
			if err != nil {
				return fmt.Errorf("failed to create Harbor client: %v", err)
			}

			userInfo, err := client.User.GetCurrentUserInfo(ctx, &user.GetCurrentUserInfoParams{})
			if err != nil {
				return fmt.Errorf("failed to get current user info: %v", err)
			}

			isSysAdmin := userInfo.Payload.SysadminFlag

			sysInfo, err := client.Systeminfo.GetSystemInfo(ctx, &systeminfo.GetSystemInfoParams{})
			if err != nil {
				return fmt.Errorf("failed to get system info: %v", err)
			}
			harborVersion := sysInfo.Payload.HarborVersion

			fmt.Println("\nHarbor CLI Info:")
			fmt.Println("==================")
			fmt.Printf("Logged in as: %s\n", userInfo.Payload.Username)
			fmt.Printf("Registry: %s\n", registryAddress)
			fmt.Printf("Harbor Version: %s\n", *harborVersion)
			fmt.Printf("Connected as Admin: %s\n", roleString(isSysAdmin))

			// Previously logged-in registries
			fmt.Println("\nPreviously Logged in to the following registries:")
			previousRegistriesMap := make(map[string]struct{})
			for _, cred := range creds {
				c := cred.(map[string]interface{})
				if registry, ok := c["serveraddress"].(string); ok {
					previousRegistriesMap[registry] = struct{}{}
				}
			}
			for registry := range previousRegistriesMap {
				fmt.Printf("- %s\n", registry)
			}

			fmt.Printf("\nCLI Version: %s\n", version.Version)
			fmt.Printf("OS: %s\n", version.System)

			return nil
		},
	}
	return cmd
}

func roleString(isSysAdmin bool) string {
	if isSysAdmin {
		return "Yes"
	}
	return "No"
}
