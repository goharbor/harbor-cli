---
title: harbor gc list
weight: 60
---
## harbor gc list

### Description

##### List GC history

### Synopsis

List GC (Garbage Collection) history in Harbor.

This command displays a list of GC executions with their status, creation time,
and other details. You can control the output using pagination flags and format options.

Examples:
  # List GC history with default pagination (page 1, 10 items per page)
  harbor gc list

  # List GC history with custom pagination
  harbor gc list --page 2 --page-size 20

  # List GC history with sorting by creation time (newest first)
  harbor gc list --sort -creation_time

  # Filter GC history by status
  harbor gc list -q status=Success

```sh
harbor gc list [flags]
```

### Options

```sh
  -h, --help            help for list
  -p, --page int        Page number (default 1)
  -s, --page-size int   Size of per page (default 10)
  -q, --query string    Query string to query resources
      --sort string     Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection

