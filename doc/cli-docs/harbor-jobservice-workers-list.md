---
title: harbor jobservice workers list
weight: 85
---
## harbor jobservice workers list

### Description

##### List workers in a pool

### Synopsis

Display all workers in the specified Harbor jobservice worker pool. Requires system admin privileges.

```sh
harbor jobservice workers list [pool-id] [flags]
```

### Examples

```sh
harbor jobservice workers list <pool-id>
```

### Options

```sh
  -h, --help   help for list
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice workers](harbor-jobservice-workers.md)	 - Manage workers

