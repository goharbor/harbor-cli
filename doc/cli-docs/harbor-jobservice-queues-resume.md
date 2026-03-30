---
title: harbor jobservice queues resume
weight: 65
---
## harbor jobservice queues resume

### Description

##### Resume queue(s) (--type or --interactive)

### Synopsis

Resume a paused job queue or all queues.

```sh
harbor jobservice queues resume [flags]
```

### Examples

```sh
harbor jobservice queues resume --type REPLICATION
harbor jobservice queues resume --type REPLICATION --type RETENTION
harbor jobservice queues resume --type all
```

### Options

```sh
  -h, --help           help for resume
  -i, --interactive    Interactive mode to choose queue type(s) instead of passing --type
      --type strings   Job type(s) to resume (repeat flag or comma-separate values; use 'all' for all queues)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice queues](harbor-jobservice-queues.md)	 - Manage job queues (list, stop, pause, resume)

