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
	robot "github.com/goharbor/harbor-cli/pkg/robot"
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/spf13/cobra"
)

func CreateRobotCommand() *cobra.Command {
	var (
		opts         robotview.CreateView
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
  level: "project"          # Optional: defaults to "project"
  permissions:              # Required: One permission scope with the project as namespace
    - namespace: "project-name"
      kind: "project"
      access:
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
			return robot.RunCreate(robot.CreateOptions{
				Level:        robot.LevelProject,
				View:         &opts,
				All:          all,
				ConfigFile:   configFile,
				ExportToFile: exportToFile,
			})
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
