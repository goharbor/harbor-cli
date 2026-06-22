---
title: harbor gc history
weight: 60
---
## harbor gc history

### Description

##### Get GC execution history

### Synopsis

Retrieve the execution history of registry-wide Garbage Collection jobs.

```sh
harbor gc history [flags]
```

### Examples

```sh
  harbor gc history --page 1 --page-size 10
```

### Options

```sh
  -h, --help            help for history
      --page int        Page number (default 1)
      --page-size int   Size of per page (default 10)
  -q, --query string    Query string to query resources
      --sort string     Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection in Harbor

