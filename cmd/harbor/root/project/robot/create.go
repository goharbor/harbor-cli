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

	"github.com/atotto/clipboard"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateRobotCommand() *cobra.Command {
	var (
		opts create.CreateView
		all  bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create robot",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if opts.ProjectName == "" {
				opts.ProjectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					logrus.Fatalf("%v", utils.ParseHarborErrorMsg(err))
				}
				if opts.ProjectName == "" {
					log.Fatalf("Project Name Cannot be empty")
				}
			}

			if len(args) == 0 {
				if opts.Name == "" || opts.Duration == 0 {
					create.CreateRobotView(&opts)
				}

				var permissions []models.Permission

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
				log.Fatalf("failed to create robot: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				name := response.Payload.Name
				res, _ := api.GetRobot(response.Payload.ID)
				utils.SavePayloadJSON(name, res.Payload)
				return
			}

			name, secret := response.Payload.Name, response.Payload.Secret
			create.CreateRobotSecretView(name, secret)
			err = clipboard.WriteAll(response.Payload.Secret)
			if err != nil {
				log.Errorf("failed to write to clipboard")
				return
			}
			fmt.Println("secret copied to clipboard.")
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(
		&all,
		"all-permission",
		"a",
		false,
		"Select all permissions for the robot account",
	)
	flags.StringVarP(&opts.ProjectName, "project", "", "", "set project name")
	flags.StringVarP(&opts.Name, "name", "", "", "name of the robot account")
	flags.StringVarP(&opts.Description, "description", "", "", "description of the robot account")
	flags.Int64VarP(&opts.Duration, "duration", "", 0, "set expiration of robot account in days")

	return cmd
}
