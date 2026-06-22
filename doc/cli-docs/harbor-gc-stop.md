---
title: harbor gc stop
weight: 50
---
## harbor gc stop

### Description

##### Stop a running GC execution

### Synopsis

Stop a currently running Garbage Collection job by its run ID.

```sh
harbor gc stop [gc-id] [flags]
```

### Examples

```sh
  harbor gc stop 12
```

### Options

```sh
  -h, --help   help for stop
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection in Harbor

