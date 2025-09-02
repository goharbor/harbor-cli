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
	config "github.com/goharbor/harbor-cli/pkg/config/robot"
	rmodel "github.com/goharbor/harbor-cli/pkg/models/robot"
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
		Long: `Create a new robot account within a Harbor project.

Robot accounts are non-human users that can be used for automation purposes
such as CI/CD pipelines, scripts, or other automated processes that need
to interact with Harbor. They have specific permissions and a defined lifetime.

This command supports both interactive and non-interactive modes:
- Without flags: opens an interactive form for configuring the robot
- With flags: creates a robot with the specified parameters
- With config file: loads robot configuration from YAML or JSON

A robot account requires:
- A unique name
- A project where it will be created
- A set of permissions
- A duration (lifetime in days)

The generated robot credentials can be:
- Displayed on screen
- Copied to clipboard (default)
- Exported to a JSON file with the -e flag

Configuration File Format (YAML or JSON):
  name: "robot-name"        # Required: Name of the robot account
  description: "..."        # Optional: Description of the robot account
  duration: 90              # Required: Lifetime in days
  project: "project-name"   # Required: Project where the robot will be created
  permissions:              # Required: At least one permission must be specified
    - resource: "repository"  # Either specify a single resource
      actions: ["pull", "push"]
    - resources: ["artifact", "scan"]  # Or specify multiple resources
      actions: ["read"]
    - resource: "project"    # Use "*" as an action to grant all available actions
      actions: ["*"]

Examples:
  # Interactive mode
  harbor-cli project robot create

  # Non-interactive mode with all flags
  harbor-cli project robot create --project myproject --name ci-robot --description "CI pipeline" --duration 90

  # Create with all permissions
  harbor-cli project robot create --project myproject --name ci-robot --all-permission

  # Load from configuration file
  harbor-cli project robot create --robot-config-file ./robot-config.yaml

  # Export secret to file
  harbor-cli project robot create --project myproject --name ci-robot --export-to-file`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var permissions []models.Permission

			if configFile != "" {
				fmt.Println("Loading configuration from: ", configFile)
				loadedOpts, loadErr := config.LoadRobotConfigFromFile(configFile)
				if loadErr != nil {
					return fmt.Errorf("failed to load robot config from file: %v", loadErr)
				}
				logrus.Info("Successfully loaded robot configuration")
				opts = *loadedOpts
				if opts.ProjectName == "" {
					opts.ProjectName = opts.Permissions[0].Namespace
				}
				permissions = make([]models.Permission, len(opts.Permissions[0].Access))
				for i, access := range opts.Permissions[0].Access {
					permissions[i] = models.Permission{
						Resource: access.Resource,
						Action:   access.Action,
					}
				}
			}

			if opts.ProjectName == "" && configFile == "" {
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
						permissions = prompt.GetRobotPermissionsFromUser("project")
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
				perm := &rmodel.RobotPermission{
					Namespace: opts.ProjectName,
					Access:    accesses,
					Kind:      "project", // Default to project level
				}
				opts.Permissions = []*rmodel.RobotPermission{perm}
			}
			getProjectID, err := api.GetProject(opts.ProjectName, false)
			if err != nil {
				return fmt.Errorf("failed to get project: %v", utils.ParseHarborErrorMsg(err))
			}
			exists, err := api.CheckRoboWithNameExists(getProjectID.Payload.ProjectID, opts.Name)
			if err != nil {
				return fmt.Errorf("failed to get robot by name: %v", utils.ParseHarborErrorMsg(err))
			}
			if exists {
				return fmt.Errorf("robot account with name '%s' already exists in project '%s'", opts.Name, opts.ProjectName)
			}
			opts.Level = "project" // Default to project level
			response, err := api.CreateRobot(opts)
			if err != nil {
				return fmt.Errorf("failed to create robot: %v", utils.ParseHarborErrorMsg(err))
			}

			logrus.Infof("Successfully created robot account '%s' (ID: %d)",
				response.Payload.Name, response.Payload.ID)

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				name := response.Payload.Name
				res, _ := api.GetRobot(response.Payload.ID)
				utils.SavePayloadJSON(name, res.Payload)
				return nil
			}
			name, secret := response.Payload.Name, response.Payload.Secret

			if exportToFile {
				logrus.Info("Exporting robot credentials to file")
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
	flags.StringVarP(&configFile, "robot-config-file", "r", "", "YAML/JSON file with robot configuration")
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
