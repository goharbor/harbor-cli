---
title: harbor jobservice queues list
weight: 60
---
## harbor jobservice queues list

### Description

##### List all job queues

### Synopsis

Display all job queues with their pending job counts and latency.

```sh
harbor jobservice queues list [flags]
```

### Examples

```sh
harbor jobservice queues list
```

### Options

```sh
  -h, --help   help for list
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice queues](harbor-jobservice-queues.md)	 - Manage job queues (list, stop, pause, resume)

