---
title: harbor jobservice schedules pause all
weight: 95
---
## harbor jobservice schedules pause-all

### Description

##### Pause all schedules

### Synopsis

Pause the global scheduler and all schedules.

```sh
harbor jobservice schedules pause-all [flags]
```

### Examples

```sh
harbor jobservice schedules pause-all
```

### Options

```sh
  -h, --help   help for pause-all
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice schedules](harbor-jobservice-schedules.md)	 - Manage schedules (list/status/pause-all/resume-all)

