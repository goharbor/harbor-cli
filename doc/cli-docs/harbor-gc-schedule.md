---
title: harbor gc schedule
weight: 40
---
## harbor gc schedule

### Description

##### Get the current GC schedule

### Synopsis

Retrieve the configuration and schedule parameters for automatic Garbage Collection.

```sh
harbor gc schedule [flags]
```

### Examples

```sh
  harbor gc schedule
```

### Options

```sh
  -h, --help   help for schedule
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection in Harbor

