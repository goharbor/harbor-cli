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

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	robotpkg "github.com/goharbor/harbor-cli/pkg/robot"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UpdateRobotCommand() *cobra.Command {
	var (
		robotID    int64
		opts       update.UpdateView
		all        bool
		configFile string
	)

	cmd := &cobra.Command{
		Use:   "update [robotID]",
		Short: "update robot by id",
		Long: `Update an existing robot account within Harbor.

Robot accounts are non-human users that can be used for automation purposes
such as CI/CD pipelines, scripts, or other automated processes that need
to interact with Harbor. This command allows you to modify an existing robot's
properties including its name, description, duration, and permissions.

This command supports both interactive and non-interactive modes:
- With robot ID: directly updates the specified robot
- Without ID: walks through robot selection interactively

The update process will:
1. Identify the robot account to be updated
2. Load its current configuration
3. Apply the requested changes
4. Save the updated configuration

This command can update both system and project-specific permissions:
- System permissions apply across the entire Harbor instance
- Project permissions apply to specific projects

Configuration can be loaded from:
- Interactive prompts (default)
- Command line flags
- YAML/JSON configuration file

Note: Updating a robot does not regenerate its secret. If you need a new
secret, consider deleting the robot and creating a new one instead.

Examples:
  # Update robot by ID with a new description
  harbor-cli robot update 123 --description "Updated CI/CD pipeline robot"

  # Update robot's duration (extend lifetime)
  harbor-cli robot update 123 --duration 180

  # Update with all permissions
  harbor-cli robot update 123 --all-permission

  # Update from configuration file
  harbor-cli robot update 123 --robot-config-file ./robot-config.yaml

  # Interactive update (will prompt for robot selection and changes)
  harbor-cli robot update`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			// Get robot ID from args or interactive prompt
			if len(args) == 1 {
				robotID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse robot ID: %v", err)
				}
			} else {
				robotID = prompt.GetRobotIDFromUser(-1)
			}

			// Get current robot configuration
			robot, err := api.GetRobot(robotID)
			if err != nil {
				return fmt.Errorf("failed to get robot: %v", utils.ParseHarborErrorMsg(err))
			}

			// Initialize update view with current values
			bot := robot.Payload
			opts.ID = bot.ID
			opts.Level = bot.Level
			opts.Name = bot.Name
			opts.Secret = bot.Secret
			opts.Description = bot.Description
			opts.Duration = *bot.Duration
			opts.Disable = bot.Disable
			opts.Editable = bot.Editable
			opts.CreationTime = bot.CreationTime

			// Extract current permissions (both system and project)
			var permissions []models.Permission
			var projectPermissionsMap = make(map[string][]models.Permission)

			// Separate system and project permissions
			for _, perm := range bot.Permissions {
				if perm.Kind == "system" && perm.Namespace == "/" {
					for _, access := range perm.Access {
						permissions = append(permissions, models.Permission{
							Resource: access.Resource,
							Action:   access.Action,
						})
					}
				} else if perm.Kind == "project" {
					var projectPerms []models.Permission
					for _, access := range perm.Access {
						projectPerms = append(projectPerms, models.Permission{
							Resource: access.Resource,
							Action:   access.Action,
						})
					}
					projectPermissionsMap[perm.Namespace] = projectPerms
				}
			}

			logrus.Infof("Loaded robot with %d system permissions and %d project-specific permissions",
				len(permissions), len(projectPermissionsMap))

			// Handle configuration from file or interactive input
			if configFile != "" {
				if err := robotpkg.LoadFromConfigFileForUpdate(&opts, configFile, &permissions, projectPermissionsMap); err != nil {
					return err
				}
			} else {
				if err := handleInteractiveInputForUpdate(&opts, all, &permissions, projectPermissionsMap); err != nil {
					return err
				}
			}

			// Build system access permissions
			var accessesSystem []*models.Access
			for _, perm := range permissions {
				accessesSystem = append(accessesSystem, &models.Access{
					Resource: perm.Resource,
					Action:   perm.Action,
				})
			}

			// Build merged permissions structure
			opts.Permissions = robotpkg.BuildMergedPermissions(projectPermissionsMap, accessesSystem)

			// Update robot and handle response
			return updateRobotAndHandleResponse(&opts)
		},
	}

	addUpdateFlags(cmd, &opts, &all, &configFile)
	return cmd
}

func handleInteractiveInputForUpdate(opts *update.UpdateView, all bool, permissions *[]models.Permission, projectPermissionsMap map[string][]models.Permission) error {
	// Show interactive form for updating basic details
	update.UpdateRobotView(opts)

	// Validate duration
	if opts.Duration == 0 {
		return fmt.Errorf("failed to update robot: %v", utils.ParseHarborErrorMsg(fmt.Errorf("duration cannot be 0")))
	}

	// Ask if user wants to update permissions
	var updatePerms bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Do you want to update permissions?").
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(&updatePerms),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).Run()

	if err != nil {
		return fmt.Errorf("error asking about permission updates: %v", err)
	}

	if !updatePerms {
		logrus.Info("Keeping existing permissions")
		return nil
	}

	// Get system permissions (update flow: ask for confirmation to change)
	if err := robotpkg.GetSystemPermissions(true, false, all, permissions); err != nil {
		return err
	}

	return robotpkg.GetProjectPermissions(true, nil, projectPermissionsMap)
}

func updateRobotAndHandleResponse(opts *update.UpdateView) error {
	err := api.UpdateRobot(opts)
	if err != nil {
		return fmt.Errorf("failed to update robot: %v", utils.ParseHarborErrorMsg(err))
	}

	logrus.Infof("Successfully updated robot account '%s' (ID: %d)", opts.Name, opts.ID)

	// Handle output format
	if formatFlag := viper.GetString("output-format"); formatFlag != "" {
		res, _ := api.GetRobot(opts.ID)
		utils.SavePayloadJSON(opts.Name, res.Payload)
	}

	return nil
}

func addUpdateFlags(cmd *cobra.Command, opts *update.UpdateView, all *bool, configFile *string) {
	flags := cmd.Flags()
	flags.BoolVarP(all, "all-permission", "a", false, "Select all permissions for the robot account")
	flags.StringVarP(&opts.Name, "name", "", "", "name of the robot account")
	flags.StringVarP(&opts.Description, "description", "", "", "description of the robot account")
	flags.Int64VarP(&opts.Duration, "duration", "", 0, "set expiration of robot account in days")
	flags.StringVarP(configFile, "robot-config-file", "r", "", "YAML/JSON file with robot configuration")
}
