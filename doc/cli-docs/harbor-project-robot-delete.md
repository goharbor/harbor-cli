---
title: harbor project robot delete
weight: 5
---
## harbor project robot delete

### Description

##### delete robot by id

### Synopsis

Delete a robot account from a Harbor project.

This command permanently removes a robot account from Harbor. Once deleted,
the robot's credentials will no longer be valid, and any automated processes
using those credentials will fail.

The command supports multiple ways to identify the robot account to delete:
- By providing the robot ID directly as an argument
- By specifying a project with the --project flag and selecting the robot interactively
- Without any arguments, which will prompt for both project and robot selection

Important considerations:
- Deletion is permanent and cannot be undone
- All access tokens for the robot will be invalidated immediately
- Any systems using the robot's credentials will need to be updated

Examples:
  # Delete robot by ID
  harbor-cli project robot delete 123

  # Delete robot by selecting from a specific project
  harbor-cli project robot delete --project myproject

  # Interactive deletion (will prompt for project and robot selection)
  harbor-cli project robot delete

```sh
harbor project robot delete [robotID] [flags]
```

### Options

```sh
  -h, --help             help for delete
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

