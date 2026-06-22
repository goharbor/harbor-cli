---
title: harbor jobservice queues pause
weight: 80
---
## harbor jobservice queues pause

### Description

##### Pause queue(s) (--type or --interactive)

### Synopsis

Pause a job queue or all queues.

```sh
harbor jobservice queues pause [flags]
```

### Examples

```sh
harbor jobservice queues pause --type REPLICATION
harbor jobservice queues pause --type REPLICATION --type RETENTION
harbor jobservice queues pause --type all
```

### Options

```sh
  -h, --help           help for pause
  -i, --interactive    Interactive mode to choose queue type(s) instead of passing --type
      --type strings   Job type(s) to pause (repeat flag or comma-separate values; use 'all' for all queues)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice queues](harbor-jobservice-queues.md)	 - Manage job queues (list, stop, pause, resume)

