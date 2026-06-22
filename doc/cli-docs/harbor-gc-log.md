---
title: harbor gc log
weight: 65
---
## harbor gc log

### Description

##### Get GC execution log

### Synopsis

Retrieve the execution log of a specific Garbage Collection run by its ID.

```sh
harbor gc log [gc-id] [flags]
```

### Examples

```sh
  harbor gc log 12
```

### Options

```sh
  -h, --help   help for log
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection in Harbor

