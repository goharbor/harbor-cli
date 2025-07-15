---
title: harbor robot view
weight: 25
---
## harbor robot view

### Description

##### get robot by id

### Synopsis

View detailed information about a robot account in Harbor.

This command displays comprehensive information about a robot account including
its ID, name, description, creation time, expiration, and the permissions
it has been granted. Supports both system-level and project-level robot accounts.

The command supports multiple ways to identify the robot account:
- By providing the robot ID directly as an argument
- Without any arguments, which will prompt for robot selection

The displayed information includes:
- Basic details (ID, name, description)
- Temporal information (creation date, expiration date, remaining time)
- Security details (disabled status)
- Detailed permissions breakdown by resource and action
- For system robots: permissions across multiple projects are shown separately

System-level robots can have permissions spanning multiple projects, while
project-level robots are scoped to a single project.

Examples:
  # View robot by ID
  harbor-cli robot view 123

  # Interactive selection (will prompt for robot)
  harbor-cli robot view

```sh
harbor robot view [robotID] [flags]
```

### Options

```sh
  -h, --help   help for view
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor robot](harbor-robot.md)	 - Manage robot accounts

