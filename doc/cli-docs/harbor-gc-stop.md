---
title: harbor gc stop
weight: 75
---
## harbor gc stop

### Description

##### Stop a running GC job

### Synopsis

Stop a running Garbage Collection job in Harbor.

This command displays a list of currently running or pending GC jobs and 
allows you to select one to stop. Only jobs with status "running" or "pending" 
can be stopped.

Examples:
  # Stop a running GC job interactively
  harbor-cli gc stop

Notes:
  - Only jobs that are currently running or pending can be stopped
  - Jobs that have already completed cannot be stopped
  - Use 'harbor-cli gc list' to view all GC jobs and their statuses

```sh
harbor gc stop [flags]
```

### Options

```sh
  -h, --help   help for stop
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection

