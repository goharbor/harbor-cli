---
title: harbor gc trigger
weight: 55
---
## harbor gc trigger

### Description

##### Trigger Garbage Collection immediately

### Synopsis

Start a manual Garbage Collection job immediately in Harbor registry.

```sh
harbor gc trigger [flags]
```

### Examples

```sh
  harbor gc trigger --delete-untagged --dry-run=false
```

### Options

```sh
      --delete-untagged   Delete untagged artifacts
      --dry-run           Simulate the GC process without deleting actual blobs
  -h, --help              help for trigger
  -i, --interactive       Trigger Garbage Collection interactively
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection in Harbor

