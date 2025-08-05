---
title: harbor replication log
weight: 10
---
## harbor replication log

### Description

##### get replication execution logs by execution and task id

### Synopsis

Get the logs of a specific replication execution and task by their IDs. If no IDs are provided, it will prompt the user to select them interactively.

```sh
harbor replication log [EXECUTION_ID] [TASK_ID] [flags]
```

### Examples

```sh
  harbor replication log -e 12345 -t 67890
		  harbor replication log --execution-id 12345 --task-id 67890
		  harbor replication log --execution-id 12345
  harbor replication log
```

### Options

```sh
  -e, --execution-id int   Replication execution ID
  -h, --help               help for log
  -t, --task-id int        Replication task ID
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor replication](harbor-replication.md)	 - Manage replications

