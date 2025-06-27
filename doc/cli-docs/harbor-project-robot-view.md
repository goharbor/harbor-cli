---
title: harbor project robot view
weight: 50
---
## harbor project robot view

### Description

##### get robot by id

### Synopsis

View detailed information about a robot account in Harbor.

This command displays comprehensive information about a robot account including
its ID, name, description, creation time, expiration, and the permissions
it has been granted within its project.

The command supports multiple ways to identify the robot account:
- By providing the robot ID directly as an argument
- By specifying a project with the --project flag and selecting the robot interactively
- Without any arguments, which will prompt for both project and robot selection

The displayed information includes:
- Basic details (ID, name, description)
- Temporal information (creation date, expiration date, remaining time)
- Security details (disabled status)
- Detailed permissions breakdown by resource and action

Examples:
  # View robot by ID
  harbor-cli project robot view 123

  # View robot by selecting from a specific project
  harbor-cli project robot view --project myproject

  # Interactive selection (will prompt for project and robot)
  harbor-cli project robot view

```sh
harbor project robot view [robotID] [flags]
```

### Options

```sh
  -h, --help             help for view
      --project string   set project name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project robot](harbor-project-robot.md)	 - Manage robot accounts

