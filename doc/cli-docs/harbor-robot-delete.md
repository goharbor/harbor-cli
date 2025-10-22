---
title: harbor robot delete
weight: 10
---
## harbor robot delete

### Description

##### delete robot by name

### Synopsis

Delete a robot account from Harbor.

This command permanently removes a robot account from Harbor. Once deleted,
the robot's credentials will no longer be valid, and any automated processes
using those credentials will fail.

The command supports multiple ways to identify the robot account to delete:
- By providing the robot name directly as an argument
- Without any arguments, which will prompt for robot selection

Important considerations:
- Deletion is permanent and cannot be undone
- All access tokens for the robot will be invalidated immediately
- Any systems using the robot's credentials will need to be updated
- For system robots, access across all projects will be revoked

Examples:
  # Delete robot by name
  harbor-cli robot delete robot_robotname

  # Interactive deletion (will prompt for robot selection)
  harbor-cli robot delete

```sh
harbor robot delete [robotName] [flags]
```

### Options

```sh
  -h, --help   help for delete
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor robot](harbor-robot.md)	 - Manage robot accounts

