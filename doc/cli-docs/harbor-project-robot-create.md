---
title: harbor project robot create
weight: 15
---
## harbor project robot create

### Description

##### create robot

```sh
harbor project robot create [flags]
```

### Options

```sh
  -a, --all-permission             Select all permissions for the robot account
      --description string         description of the robot account
      --duration int               set expiration of robot account in days
  -h, --help                       help for create
      --name string                name of the robot account
      --project string             set project name
  -r, --robot-config-file string   YAML file with robot configuration
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project robot](harbor-project-robot.md)	 - Manage robot accounts

