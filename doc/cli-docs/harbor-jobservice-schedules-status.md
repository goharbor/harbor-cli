---
title: harbor jobservice schedules status
weight: 70
---
## harbor jobservice schedules status

### Description

##### Show scheduler status

### Synopsis

Display whether the global scheduler is paused or running.

```sh
harbor jobservice schedules status [flags]
```

### Examples

```sh
harbor jobservice schedules status
```

### Options

```sh
  -h, --help   help for status
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice schedules](harbor-jobservice-schedules.md)	 - Manage schedules (list/status/pause-all/resume-all)

