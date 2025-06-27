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
	"github.com/charmbracelet/huh"
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
			var err error
			var permissions []models.Permission
			var projectPermissionsMap map[string][]models.Permission = make(map[string][]models.Permission)

			if len(args) == 0 {
				if configFile != "" {
					fmt.Println("Loading configuration from: ", configFile)
					loadedOpts, loadErr := config.LoadRobotConfigFromFile(configFile)
					if loadErr != nil {
						return fmt.Errorf("failed to load robot config from file: %v", loadErr)
					}
					logrus.Info("Successfully loaded robot configuration")

					opts = *loadedOpts

					// Extract system-level permissions
					var systemPermFound bool
					for _, perm := range opts.Permissions {
						if perm.Kind == "system" && perm.Namespace == "/" {
							systemPermFound = true
							permissions = make([]models.Permission, len(perm.Access))
							for i, access := range perm.Access {
								permissions[i] = models.Permission{
									Resource: access.Resource,
									Action:   access.Action,
								}
							}
						} else if perm.Kind == "project" {
							// Handle project-specific permissions
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

					if !systemPermFound {
						return fmt.Errorf("system robot configuration must include system-level permissions")
					}

					// Skip interactive permission collection
					logrus.Infof("Loaded system robot with %d system permissions and %d project-specific permissions",
						len(permissions), len(projectPermissionsMap))
				}

				if (opts.Name == "" || opts.Duration == 0) && configFile == "" {
					create.CreateRobotView(&opts)
				}

				if opts.Duration == 0 {
					msg := fmt.Errorf("duration cannot be 0")
					return fmt.Errorf("failed to create robot: %v", utils.ParseHarborErrorMsg(msg))
				}

				if len(permissions) == 0 {
					if all {
						perms, _ := api.GetPermissions()
						permission := perms.Payload.System

						choices := []models.Permission{}
						for _, perm := range permission {
							choices = append(choices, *perm)
						}
						permissions = choices
					} else {
						permissions = prompt.GetRobotPermissionsFromUser("system")
						if len(permissions) == 0 {
							msg := fmt.Errorf("no permissions selected, robot account needs at least one permission")
							return fmt.Errorf("failed to create robot: %v", utils.ParseHarborErrorMsg(msg))
						}
					}
				}

				var accessesSystem []*models.Access
				for _, perm := range permissions {
					access := &models.Access{
						Resource: perm.Resource,
						Action:   perm.Action,
					}
					accessesSystem = append(accessesSystem, access)
				}

				permissionMode, err := promptPermissionMode()
				if err != nil {
					return fmt.Errorf("error selecting permission mode: %v", err)
				}

				if permissionMode == "list" {
					selectedProjects, err := getMultipleProjectsFromUser()
					if err != nil {
						return fmt.Errorf("error selecting projects: %v", err)
					}

					if len(selectedProjects) > 0 {
						fmt.Println("Select permissions to apply to all selected projects:")
						projectPermissions := prompt.GetRobotPermissionsFromUser("project")

						for _, projectName := range selectedProjects {
							projectPermissionsMap[projectName] = projectPermissions
						}
					}
				} else if permissionMode == "per_project" {
					if opts.ProjectName == "" {
						for {
							projectName, err := prompt.GetProjectNameFromUser()
							if err != nil {
								return fmt.Errorf("%v", utils.ParseHarborErrorMsg(err))
							}
							if projectName == "" {
								return fmt.Errorf("project name cannot be empty")
							}
							projectPermissionsMap[projectName] = prompt.GetRobotPermissionsFromUser("project")
							moreProjects, err := promptMoreProjects()
							if err != nil {
								return fmt.Errorf("error asking for more projects: %v", err)
							}
							if !moreProjects {
								break
							}
						}
					} else {
						projectPermissions := prompt.GetRobotPermissionsFromUser("project")
						projectPermissionsMap[opts.ProjectName] = projectPermissions
					}
				} else if permissionMode == "none" {
					fmt.Println("Creating robot with system-level permissions only (no project-specific permissions)")
				}

				var mergedPermissions []*create.RobotPermission
				for projectName, projectPermissions := range projectPermissionsMap {
					var accessesProject []*models.Access
					for _, perm := range projectPermissions {
						access := &models.Access{
							Resource: perm.Resource,
							Action:   perm.Action,
						}
						accessesProject = append(accessesProject, access)
					}
					mergedPermissions = append(mergedPermissions, &create.RobotPermission{
						Namespace: projectName,
						Access:    accessesProject,
						Kind:      "project",
					})
				}

				mergedPermissions = append(mergedPermissions, &create.RobotPermission{
					Namespace: "/",
					Access:    accessesSystem,
					Kind:      "system",
				})
				opts.Permissions = mergedPermissions
			}

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

func getMultipleProjectsFromUser() ([]string, error) {
	allProjects, err := api.ListAllProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %v", err)
	}

	var selectedProjects []string
	projectOptions := []huh.Option[string]{}

	for _, p := range allProjects.Payload {
		projectOptions = append(projectOptions,
			huh.NewOption(p.Name, p.Name))
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Multiple Project Selection").
				Description("Select the projects to assign the same permissions to this robot account."),
			huh.NewMultiSelect[string]().
				Title("Select projects").
				Options(projectOptions...).
				Value(&selectedProjects),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(80).Run()

	return selectedProjects, err
}

func promptMoreProjects() (bool, error) {
	var addMore bool = false
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Project Selection").
				Description("You can add permissions for multiple projects to this robot account."),
			huh.NewSelect[bool]().
				Title("Do you want to select (more) projects?").
				Description("Select 'Yes' to add (another) project, 'No' to continue with current selection.").
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(&addMore),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).WithHeight(10).Run()

	return addMore, err
}

func promptPermissionMode() (string, error) {
	var permissionMode string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Permission Mode").
				Description("Select how you want to assign permissions to projects:"),
			huh.NewSelect[string]().
				Title("Permission Mode").
				Description("Choose 'List' to select multiple projects with common permissions, or 'Per Project' for individual project permissions.").
				Options(
					huh.NewOption("No project permissions (system-level only)", "none"),
					huh.NewOption("Per Project", "per_project"),
					huh.NewOption("List", "list"),
				).
				Value(&permissionMode),
		),
	).WithTheme(huh.ThemeCharm()).WithWidth(60).WithHeight(10).Run()

	return permissionMode, err
}
