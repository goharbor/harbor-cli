---
title: harbor jobservice pools list
weight: 70
---
## harbor jobservice pools list

### Description

##### List all worker pools

### Synopsis

Display all Harbor jobservice worker pools with host, concurrency and start time. Requires system admin privileges.

```sh
harbor jobservice pools list [flags]
```

### Examples

```sh
harbor jobservice pools list
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

* [harbor jobservice pools](harbor-jobservice-pools.md)	 - Manage worker pools

