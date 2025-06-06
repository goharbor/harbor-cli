---
title: harbor project robot update
weight: 55
---
## harbor project robot update

### Description

##### update robot by id

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
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project robot](harbor-project-robot.md)	 - Manage robot accounts

