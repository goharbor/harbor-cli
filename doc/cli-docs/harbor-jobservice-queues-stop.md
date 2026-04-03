---
title: harbor jobservice queues stop
weight: 60
---
## harbor jobservice queues stop

### Description

##### Stop queue(s) (--type or --interactive)

### Synopsis

Stop a job queue or all queues.

```sh
harbor jobservice queues stop [flags]
```

### Examples

```sh
harbor jobservice queues stop --type REPLICATION
harbor jobservice queues stop --type REPLICATION --type RETENTION
harbor jobservice queues stop --type all
```

### Options

```sh
  -h, --help           help for stop
  -i, --interactive    Interactive mode to choose queue type(s) instead of passing --type
      --type strings   Job type(s) to stop (repeat flag or comma-separate values; use 'all' for all queues)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice queues](harbor-jobservice-queues.md)	 - Manage job queues (list, stop, pause, resume)

