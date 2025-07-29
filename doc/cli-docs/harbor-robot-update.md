---
title: harbor robot update
weight: 40
---
## harbor robot update

### Description

##### update robot by id

### Synopsis

Update an existing robot account within Harbor.

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
  harbor-cli robot update

```sh
harbor robot update [robotID] [flags]
```

### Options

```sh
  -a, --all-permission             Select all permissions for the robot account
      --description string         description of the robot account
      --duration int               set expiration of robot account in days
  -h, --help                       help for update
      --name string                name of the robot account
  -r, --robot-config-file string   YAML/JSON file with robot configuration
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor robot](harbor-robot.md)	 - Manage robot accounts

