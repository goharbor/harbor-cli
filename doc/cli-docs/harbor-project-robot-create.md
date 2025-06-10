---
title: harbor project robot create
weight: 15
---
## harbor project robot create

### Description

##### create robot

### Synopsis

Create a new robot account within a Harbor project.

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
  harbor-cli project robot create --project myproject --name ci-robot --export-to-file

```sh
harbor project robot create [flags]
```

### Options

```sh
  -a, --all-permission             Select all permissions for the robot account
      --description string         description of the robot account
      --duration int               set expiration of robot account in days
  -e, --export-to-file             Choose to export robot account to file
  -h, --help                       help for create
      --name string                name of the robot account
      --project string             set project name
  -r, --robot-config-file string   YAML/JSON file with robot configuration
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project robot](harbor-project-robot.md)	 - Manage robot accounts

