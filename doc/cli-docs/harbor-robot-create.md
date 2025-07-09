---
title: harbor robot create
weight: 15
---
## harbor robot create

### Description

##### create robot

### Synopsis

Create a new robot account within Harbor.

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
  harbor-cli robot create --name ci-robot --export-to-file

```sh
harbor robot create [flags]
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

* [harbor robot](harbor-robot.md)	 - Manage robot accounts

