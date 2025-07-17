---
title: harbor replication executions view
weight: 65
---
## harbor replication executions view

### Description

##### get replication execution by id

### Synopsis

Get a specific replication execution by its ID. If no ID is provided, it will prompt the user to select one interactively. If the --no-tasks flag is set, it will not list associated tasks.

```sh
harbor replication executions view [ID] [flags]
```

### Examples

```sh
  harbor replication executions view 12345
  harbor replication executions view --no-tasks
```

### Options

```sh
  -h, --help       help for view
      --no-tasks   Do not list associated tasks
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor replication executions](harbor-replication-executions.md)	 - Manage replication executions

