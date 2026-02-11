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

  # List GC history with multiple sort fields
  harbor gc list --sort creation_time --sort -job_status

  # Filter GC history by status (exact match)
  harbor gc list --match job_status=Success

  # Filter GC history by fuzzy match
  harbor gc list --fuzzy job_name=gc

```sh
harbor gc list [flags]
```

### Options

```sh
      --fuzzy strings   Fuzzy match filter (key=value)
  -h, --help            help for list
      --match strings   Exact match filter (key=value)
  -p, --page int        Page number (default 1)
  -s, --page-size int   Size of per page (default 10)
      --range strings   Range filter (key=min~max)
      --sort strings    Sort the resource list (e.g. --sort creation_time --sort -update_time)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection

