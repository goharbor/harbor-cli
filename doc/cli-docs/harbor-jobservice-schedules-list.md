---
title: harbor jobservice schedules list
weight: 10
---
## harbor jobservice schedules list

### Description

##### List schedules (supports --page and --page-size)

### Synopsis

Display all job schedules with pagination support.

```sh
harbor jobservice schedules list [flags]
```

### Examples

```sh
harbor jobservice schedules list --page 1 --page-size 20
```

### Options

```sh
  -h, --help            help for list
      --page int        Page number (default 1)
      --page-size int   Number of items per page (default 20)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice schedules](harbor-jobservice-schedules.md)	 - Manage schedules (list/status/pause-all/resume-all)

