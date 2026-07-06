---
title: harbor gc
weight: 15
---
## harbor gc

### Description

##### Manage Garbage Collection in Harbor

### Synopsis

Use this command to manage registry-wide Garbage Collection (GC) in your Harbor instance.

Garbage Collection cleans up deleted or orphaned blobs/tags in the registry to free up storage space.
This command supports listing execution history.

### Examples

```sh
  # View Garbage Collection execution history
  harbor gc history
```

### Options

```sh
  -h, --help   help for gc
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor gc history](harbor-gc-history.md)	 - Get GC execution history

