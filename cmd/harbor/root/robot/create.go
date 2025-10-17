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
	robotpkg "github.com/goharbor/harbor-cli/pkg/robot"
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
		Long: `Create a new robot account within Harbor.

Robot accounts are non-human users that can be used for automation purposes
such as CI/CD pipelines, scripts, or other automated processes that need
to interact with Harbor. They have specific permissions and a defined lifetime.

This command creates system-level robots that can have permissions spanning 
multiple projects, making them suitable for automation tasks that need access 
across your Harbor instance.

This command supports both interactive and non-interactive modes:
- Without flags: opens an interactive form for configuring the robot
- With flags: creates a robot with the specified parameters
- With config file: loads robot configuration from YAML or JSON

A robot account requires:
- A unique name
- A set of system permissions
- Optional project-specific permissions
- A duration (lifetime in days)

The generated robot credentials can be:
- Displayed on screen
- Copied to clipboard (default)
- Exported to a JSON file with the -e flag

Examples:
  # Interactive mode
  harbor-cli robot create

  # Non-interactive mode with all flags
  harbor-cli robot create --name ci-robot --description "CI pipeline" --duration 90

  # Create with all permissions
  harbor-cli robot create --name ci-robot --all-permission

  # Load from configuration file
  harbor-cli robot create --robot-config-file ./robot-config.yaml

  # Export secret to file
  harbor-cli robot create --name ci-robot --export-to-file`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var permissions []models.Permission
			var projectPermissionsMap = make(map[string][]models.Permission)
			var accessesSystem []*models.Access

			// Handle config file or interactive input
			if configFile != "" {
				if err := robotpkg.LoadFromConfigFileForCreate(&opts, configFile, &permissions, projectPermissionsMap); err != nil {
					return fmt.Errorf("failed to load robot config from file: %v", err)
				}
				logrus.Info("Successfully loaded robot configuration")
			} else {
				if err := handleInteractiveInput(&opts, all, &permissions, projectPermissionsMap); err != nil {
					return err
				}
			}

			accessesSystem = robotpkg.PermissionsToAccess(permissions)

			// Build merged permissions structure
			opts.Permissions = robotpkg.BuildMergedPermissions(projectPermissionsMap, accessesSystem)
			opts.Level = "system"

			// Create robot and handle response
			return createRobotAndHandleResponse(&opts, exportToFile)
		},
	}

	addFlags(cmd, &opts, &all, &exportToFile, &configFile)
	return cmd
}

func handleInteractiveInput(opts *create.CreateView, all bool, permissions *[]models.Permission, projectPermissionsMap map[string][]models.Permission) error {
	// Show interactive form if needed
	if opts.Name == "" || opts.Duration == 0 {
		create.CreateRobotView(opts)
	}

	// Validate duration
	if opts.Duration == 0 {
		return fmt.Errorf("failed to create robot: %v", utils.ParseHarborErrorMsg(fmt.Errorf("duration cannot be 0")))
	}

	// Get system permissions (create flow: no update confirmation)
	if err := robotpkg.GetSystemPermissions(false, true, all, permissions); err != nil {
		return err
	}

	return robotpkg.GetProjectPermissions(false, opts, projectPermissionsMap)
}

func createRobotAndHandleResponse(opts *create.CreateView, exportToFile bool) error {
	response, err := api.CreateRobot(*opts)
	if err != nil {
		return fmt.Errorf("failed to create robot: %v", utils.ParseHarborErrorMsg(err))
	}

	logrus.Infof("Successfully created robot account '%s' (ID: %d)",
		response.Payload.Name, response.Payload.ID)

	// Handle output format
	if formatFlag := viper.GetString("output-format"); formatFlag != "" {
		res, _ := api.GetRobot(response.Payload.ID)
		utils.SavePayloadJSON(response.Payload.Name, res.Payload)
		return nil
	}

	// Handle secret output
	name, secret := response.Payload.Name, response.Payload.Secret

	if exportToFile {
		logrus.Info("Exporting robot credentials to file")
		exportSecretToFile(name, secret, response.Payload.CreationTime.String(), response.Payload.ExpiresAt)
		return nil
	}

	create.CreateRobotSecretView(name, secret)
	if err := clipboard.WriteAll(secret); err != nil {
		logrus.Errorf("failed to write to clipboard")
	} else {
		fmt.Println("secret copied to clipboard.")
	}

	return nil
}

func addFlags(cmd *cobra.Command, opts *create.CreateView, all *bool, exportToFile *bool, configFile *string) {
	flags := cmd.Flags()
	flags.BoolVarP(all, "all-permission", "a", false, "Select all permissions for the robot account")
	flags.BoolVarP(exportToFile, "export-to-file", "e", false, "Choose to export robot account to file")
	flags.StringVarP(&opts.ProjectName, "project", "", "", "set project name")
	flags.StringVarP(&opts.Name, "name", "", "", "name of the robot account")
	flags.StringVarP(&opts.Description, "description", "", "", "description of the robot account")
	flags.Int64VarP(&opts.Duration, "duration", "", 0, "set expiration of robot account in days")
	flags.StringVarP(configFile, "robot-config-file", "r", "", "YAML/JSON file with robot configuration")
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
		return
	}

	if err := os.WriteFile(filename, jsonData, 0600); err != nil {
		logrus.Errorf("Failed to write secret to file: %v", err)
		return
	}

	fmt.Printf("Secret saved to %s\n", filename)
}
