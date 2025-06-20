---
title: harbor project robot update
weight: 55
---
## harbor project robot update

### Description

##### update robot by id

### Synopsis

Update an existing robot account within a Harbor project.

Robot accounts are non-human users that can be used for automation purposes
such as CI/CD pipelines, scripts, or other automated processes that need
to interact with Harbor. This command allows you to modify an existing robot's
properties including its name, description, duration, and permissions.

This command supports both interactive and non-interactive modes:
- With robot ID: directly updates the specified robot
- With --project flag: helps select a robot from the specified project
- Without either: walks through project and robot selection interactively

The update process will:
1. Identify the robot account to be updated
2. Load its current configuration
3. Apply the requested changes
4. Save the updated configuration

Fields that can be updated:
- Name: The robot account's identifier
- Description: A human-readable description of the robot's purpose
- Duration: The lifetime of the robot account in days
- Permissions: The actions the robot is allowed to perform

Note: Updating a robot does not regenerate its secret. If you need a new
secret, consider deleting the robot and creating a new one instead.

Examples:
  # Update robot by ID with a new description
  harbor-cli project robot update 123 --description "Updated CI/CD pipeline robot"

  # Update robot's duration (extend lifetime)
  harbor-cli project robot update 123 --duration 180

  # Update by selecting from a specific project
  harbor-cli project robot update --project myproject

  # Update with all permissions
  harbor-cli project robot update 123 --all-permission

  # Interactive update (will prompt for robot selection and changes)
  harbor-cli project robot update

```sh
harbor project robot update [robotID] [flags]
```

### Options

```sh
  -a, --all-permission       Select all permissions for the robot account
      --description string   description of the robot account
      --duration int         set expiration of robot account in days
  -h, --help                 help for update
      --name string          name of the robot account
      --project string       set project name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project robot](harbor-project-robot.md)	 - Manage robot accounts

