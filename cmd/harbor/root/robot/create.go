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
	config "github.com/goharbor/harbor-cli/pkg/config/robot"
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
			var permissions []models.Permission
			var projectPermissionsMap = make(map[string][]models.Permission)
			var accessesSystem []*models.Access

			// Handle config file or interactive input
			if configFile != "" {
				if err := loadFromConfigFile(&opts, configFile, &permissions, projectPermissionsMap); err != nil {
					return err
				}
			} else {
				if err := handleInteractiveInput(&opts, all, &permissions, projectPermissionsMap); err != nil {
					return err
				}
			}

			// Build system access permissions
			for _, perm := range permissions {
				accessesSystem = append(accessesSystem, &models.Access{
					Resource: perm.Resource,
					Action:   perm.Action,
				})
			}

			// Build merged permissions structure
			opts.Permissions = buildMergedPermissions(projectPermissionsMap, accessesSystem)
			opts.Level = "system"

			// Create robot and handle response
			return createRobotAndHandleResponse(&opts, exportToFile)
		},
	}

	addFlags(cmd, &opts, &all, &exportToFile, &configFile)
	return cmd
}

func loadFromConfigFile(opts *create.CreateView, configFile string, permissions *[]models.Permission, projectPermissionsMap map[string][]models.Permission) error {
	fmt.Println("Loading configuration from: ", configFile)

	loadedOpts, err := config.LoadRobotConfigFromFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to load robot config from file: %v", err)
	}

	logrus.Info("Successfully loaded robot configuration")
	*opts = *loadedOpts

	// Extract system-level and project permissions
	var systemPermFound bool
	for _, perm := range opts.Permissions {
		if perm.Kind == "system" && perm.Namespace == "/" {
			systemPermFound = true
			*permissions = make([]models.Permission, len(perm.Access))
			for i, access := range perm.Access {
				(*permissions)[i] = models.Permission{
					Resource: access.Resource,
					Action:   access.Action,
				}
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

	if !systemPermFound {
		return fmt.Errorf("system robot configuration must include system-level permissions")
	}

	logrus.Infof("Loaded system robot with %d system permissions and %d project-specific permissions",
		len(*permissions), len(projectPermissionsMap))

	return nil
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

	// Get system permissions
	if err := getSystemPermissions(all, permissions); err != nil {
		return err
	}

	// Get project permissions
	return getProjectPermissions(opts, projectPermissionsMap)
}

func getSystemPermissions(all bool, permissions *[]models.Permission) error {
	if len(*permissions) == 0 {
		if all {
			perms, _ := api.GetPermissions()
			for _, perm := range perms.Payload.System {
				*permissions = append(*permissions, *perm)
			}
		} else {
			*permissions = prompt.GetRobotPermissionsFromUser("system")
			if len(*permissions) == 0 {
				return fmt.Errorf("failed to create robot: %v",
					utils.ParseHarborErrorMsg(fmt.Errorf("no permissions selected, robot account needs at least one permission")))
			}
		}
	}
	return nil
}

func getProjectPermissions(opts *create.CreateView, projectPermissionsMap map[string][]models.Permission) error {
	permissionMode, err := promptPermissionMode()
	if err != nil {
		return fmt.Errorf("error selecting permission mode: %v", err)
	}

	switch permissionMode {
	case "list":
		return handleMultipleProjectsPermissions(projectPermissionsMap)
	case "per_project":
		return handlePerProjectPermissions(opts, projectPermissionsMap)
	case "none":
		fmt.Println("Creating robot with system-level permissions only (no project-specific permissions)")
		return nil
	default:
		return fmt.Errorf("unknown permission mode: %s", permissionMode)
	}
}

func handleMultipleProjectsPermissions(projectPermissionsMap map[string][]models.Permission) error {
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

	return nil
}

func handlePerProjectPermissions(opts *create.CreateView, projectPermissionsMap map[string][]models.Permission) error {
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

	return nil
}

func buildMergedPermissions(projectPermissionsMap map[string][]models.Permission, accessesSystem []*models.Access) []*create.RobotPermission {
	var mergedPermissions []*create.RobotPermission

	// Add project permissions
	for projectName, projectPermissions := range projectPermissionsMap {
		var accessesProject []*models.Access
		for _, perm := range projectPermissions {
			accessesProject = append(accessesProject, &models.Access{
				Resource: perm.Resource,
				Action:   perm.Action,
			})
		}
		mergedPermissions = append(mergedPermissions, &create.RobotPermission{
			Namespace: projectName,
			Access:    accessesProject,
			Kind:      "project",
		})
	}

	// Add system permissions
	mergedPermissions = append(mergedPermissions, &create.RobotPermission{
		Namespace: "/",
		Access:    accessesSystem,
		Kind:      "system",
	})

	return mergedPermissions
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

func getMultipleProjectsFromUser() ([]string, error) {
	allProjects, err := api.ListAllProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %v", err)
	}

	var selectedProjects []string
	var projectOptions []huh.Option[string]

	for _, p := range allProjects.Payload {
		projectOptions = append(projectOptions, huh.NewOption(p.Name, p.Name))
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
	var addMore bool
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
