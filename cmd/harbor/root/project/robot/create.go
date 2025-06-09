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
	"encoding/json"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/config"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateRobotCommand() *cobra.Command {
	var (
		opts         create.CreateView
		all          bool
		exportToFile bool
		configFile   string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create robot",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var permissions []models.Permission

			if configFile != "" {
				fmt.Println("Loading configuration from: ", configFile)
				loadedOpts, loadErr := config.LoadRobotConfigFromFile(configFile)
				if loadErr != nil {
					return fmt.Errorf("failed to load robot config from file: %v", loadErr)
				}
				opts = *loadedOpts
				permissions = make([]models.Permission, len(opts.Permissions[0].Access))
				for i, access := range opts.Permissions[0].Access {
					permissions[i] = models.Permission{
						Resource: access.Resource,
						Action:   access.Action,
					}
				}
			}

			if opts.ProjectName == "" {
				opts.ProjectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("%v", utils.ParseHarborErrorMsg(err))
				}
				if opts.ProjectName == "" {
					return fmt.Errorf("project name cannot be empty")
				}
			}

			if len(args) == 0 {
				if (opts.Name == "" || opts.Duration == 0) && configFile == "" {
					fmt.Println("Opening interactive form for robot creation")
					create.CreateRobotView(&opts)
				}

				if opts.Duration == 0 {
					msg := fmt.Errorf("duration cannot be 0")
					return fmt.Errorf("failed to create robot: %v", utils.ParseHarborErrorMsg(msg))
				}

				if len(permissions) == 0 {
					if all {
						perms, _ := api.GetPermissions()
						permission := perms.Payload.Project

						choices := []models.Permission{}
						for _, perm := range permission {
							choices = append(choices, *perm)
						}
						permissions = choices
					} else {
						permissions = prompt.GetRobotPermissionsFromUser()
						if len(permissions) == 0 {
							msg := fmt.Errorf("no permissions selected, robot account needs at least one permission")
							return fmt.Errorf("failed to create robot: %v", utils.ParseHarborErrorMsg(msg))
						}
					}
				}

				// []Permission to []*Access
				var accesses []*models.Access
				for _, perm := range permissions {
					access := &models.Access{
						Action:   perm.Action,
						Resource: perm.Resource,
					}
					accesses = append(accesses, access)
				}
				// convert []models.permission to []*model.Access
				perm := &create.RobotPermission{
					Namespace: opts.ProjectName,
					Access:    accesses,
				}
				opts.Permissions = []*create.RobotPermission{perm}
			}
			response, err := api.CreateRobot(opts, "project")
			if err != nil {
				return fmt.Errorf("failed to create robot: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				name := response.Payload.Name
				res, _ := api.GetRobot(response.Payload.ID)
				utils.SavePayloadJSON(name, res.Payload)
				return nil
			}
			name, secret := response.Payload.Name, response.Payload.Secret

			if exportToFile {
				exportSecretToFile(name, secret, response.Payload.CreationTime.String(), response.Payload.ExpiresAt)
				return nil
			} else {
				create.CreateRobotSecretView(name, secret)
				err = clipboard.WriteAll(response.Payload.Secret)
				if err != nil {
					logrus.Errorf("failed to write to clipboard")
					return nil
				}
				fmt.Println("secret copied to clipboard.")
				return nil
			}
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&all, "all-permission", "a", false, "Select all permissions for the robot account")
	flags.BoolVarP(&exportToFile, "export-to-file", "e", false, "Choose to export robot account to file")

	flags.StringVarP(&opts.ProjectName, "project", "", "", "set project name")
	flags.StringVarP(&opts.Name, "name", "", "", "name of the robot account")
	flags.StringVarP(&opts.Description, "description", "", "", "description of the robot account")
	flags.Int64VarP(&opts.Duration, "duration", "", 0, "set expiration of robot account in days")
	flags.StringVarP(&configFile, "robot-config-file", "r", "", "YAML file with robot configuration")
	return cmd
}

func exportSecretToFile(name, secret, creationTime string, expiresAt int64) {
	secretJson := config.RobotSecret{
		Name:         name,
		ExpiresAt:    expiresAt,
		CreationTime: creationTime,
		Secret:       secret,
	}
	filename := fmt.Sprintf("%s-secret.json", name)
	jsonData, err := json.MarshalIndent(secretJson, "", "  ")
	if err != nil {
		logrus.Errorf("Failed to marshal secret to JSON: %v", err)
	} else {
		if err := os.WriteFile(filename, jsonData, 0600); err != nil {
			logrus.Errorf("Failed to write secret to file: %v", err)
		} else {
			fmt.Printf("Secret saved to %s\n", filename)
		}
	}

}
